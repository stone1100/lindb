// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package memdb

import (
	"io"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("TSDB", "MemDB")

var (
	memDBScope               = linmetric.StorageRegistry.NewScope("lindb.tsdb.memdb")
	pageAllocatedCounterVec  = memDBScope.NewCounterVec("allocated_pages", "db")
	pageAllocatedFailuresVec = memDBScope.NewCounterVec("allocated_page_failures", "db")
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// AcquireWrite acquires writing data points
	AcquireWrite()
	// WithLock retrieves the lock of memdb, and returns the release function
	WithLock() (release func())
	// WriteRow must be called after WithLock
	// Used for batch write
	WriteRow(row *metric.StorageRow) error
	// CompleteWrite completes writing data points
	CompleteWrite()
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int64
	// DataFilter filters the data based on condition
	flow.DataFilter
	// Closer closes the memory database resource
	io.Closer
	// FamilyTime returns the family time of this memdb
	FamilyTime() int64
	// Uptime returns duration since created
	Uptime() time.Duration
	// Size returns the number of metric names.
	Size() int
}

type memoryDBMetrics struct {
	allocatedPages        *linmetric.BoundCounter
	allocatedPageFailures *linmetric.BoundCounter
}

func newMemoryDBMetrics(name string) *memoryDBMetrics {
	return &memoryDBMetrics{
		allocatedPages:        pageAllocatedCounterVec.WithTagValues(name),
		allocatedPageFailures: pageAllocatedFailuresVec.WithTagValues(name),
	}
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	FamilyTime int64
	Name       string
	BufferMgr  BufferManager
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID uint32

	timeutil.SlotRange // start/end time slot, metric level flush context
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	allocSize atomic.Int64 // allocated size

	familyTime int64
	name       string

	mStores *MetricBucketStore // metric id => mStoreINTF
	buf     DataPointBuffer

	writeCondition sync.WaitGroup
	rwMutex        sync.RWMutex // lock of create metric store

	metrics     memoryDBMetrics
	createdTime int64
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg MemoryDatabaseCfg) (MemoryDatabase, error) {
	buf, err := cfg.BufferMgr.AllocBuffer()
	if err != nil {
		return nil, err
	}
	return &memoryDatabase{
		familyTime:  cfg.FamilyTime,
		name:        cfg.Name,
		buf:         buf,
		mStores:     NewMetricBucketStore(),
		allocSize:   *atomic.NewInt64(0),
		metrics:     *newMemoryDBMetrics(cfg.Name),
		createdTime: fasttime.UnixNano(),
	}, err
}

func (md *memoryDatabase) FamilyTime() int64 { return md.familyTime }

func (md *memoryDatabase) metricBucketSize() int {
	var size int
	size += cap(md.mStores.values)*24 + 24
	for idx := range md.mStores.values {
		size += cap(md.mStores.values[idx])*8 + 24
	}
	return size
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricID metric.ID) (mStore mStoreINTF) {
	metricKey := uint32(metricID)
	mStore, ok := md.mStores.Get(metricKey)
	if !ok {
		// not found need create new metric store
		beforeMetricBucketSize := md.metricBucketSize()
		mStore = newMetricStore()
		// add metric-store size
		md.allocSize.Add(int64(mStore.Capacity()))
		// add metric-bucket increased
		md.mStores.Put(metricKey, mStore)
		md.allocSize.Add(int64(md.metricBucketSize() - beforeMetricBucketSize))
	}
	// found metric store in current memory database
	return
}

// AcquireWrite acquires writing data points
func (md *memoryDatabase) AcquireWrite() {
	md.writeCondition.Add(1)
}

// CompleteWrite completes writing data points
func (md *memoryDatabase) CompleteWrite() {
	md.writeCondition.Done()
}

func (md *memoryDatabase) WithLock() (release func()) {
	md.rwMutex.Lock()
	return md.rwMutex.Unlock
}

func (md *memoryDatabase) WriteRow(row *metric.StorageRow) error {
	mStore := md.getOrCreateMStore(row.MetricID)
	var size int
	beforeMStoreCapacity := mStore.Capacity()
	tStore, created := mStore.GetOrCreateTStore(row.SeriesID)
	if created {
		size += tStore.Capacity()
		size += mStore.Capacity() - beforeMStoreCapacity
	}
	written := false
	var fieldIDIdx = 0
	afterWrite := func(writtenLinFieldSize int) {
		fieldIDIdx++
		size += writtenLinFieldSize
		written = true
	}

	simpleFieldItr := row.NewSimpleFieldIterator()
	for simpleFieldItr.HasNext() {
		writtenLinFieldSize, err := md.writeLinField(
			row.SlotIndex,
			row.FieldIDs[fieldIDIdx],
			simpleFieldItr.NextType(),
			simpleFieldItr.NextValue(),
			mStore, tStore,
		)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	compoundFieldItr, ok := row.NewCompoundFieldIterator()

	var (
		err                 error
		writtenLinFieldSize int
	)
	if !ok {
		goto End
	}

	// write histogram_min
	if compoundFieldItr.Min() > 0 {
		writtenLinFieldSize, err = md.writeLinField(
			row.SlotIndex, row.FieldIDs[fieldIDIdx],
			field.MinField, compoundFieldItr.Min(),
			mStore, tStore)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	// write histogram_max
	if compoundFieldItr.Max() > 0 {
		writtenLinFieldSize, err = md.writeLinField(
			row.SlotIndex, row.FieldIDs[fieldIDIdx],
			field.MaxField, compoundFieldItr.Max(),
			mStore, tStore)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	// write histogram_sum
	writtenLinFieldSize, err = md.writeLinField(
		row.SlotIndex, row.FieldIDs[fieldIDIdx],
		field.SumField, compoundFieldItr.Sum(),
		mStore, tStore)
	if err != nil {
		return err
	}
	afterWrite(writtenLinFieldSize)

	// write histogram_count
	writtenLinFieldSize, err = md.writeLinField(
		row.SlotIndex, row.FieldIDs[fieldIDIdx],
		field.SumField, compoundFieldItr.Count(),
		mStore, tStore)
	if err != nil {
		return err
	}
	afterWrite(writtenLinFieldSize)

	// write __bucket_${boundary}
	// assume that length of ExplicitBounds equals to Values
	// data must be valid before write
	for compoundFieldItr.HasNextBucket() {
		writtenLinFieldSize, err = md.writeLinField(
			row.SlotIndex, row.FieldIDs[fieldIDIdx],
			field.HistogramField, compoundFieldItr.NextValue(),
			mStore, tStore)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}

End:
	if written {
		mStore.SetSlot(row.SlotIndex)
	}
	md.allocSize.Add(int64(size))
	return nil
}

func (md *memoryDatabase) writeLinField(
	slotIndex uint16,
	fieldID field.ID, fieldType field.Type, fieldValue float64,
	mStore mStoreINTF, tStore tStoreINTF,
) (writtenSize int, err error) {
	fStore, ok := tStore.GetFStore(fieldID)
	if !ok {
		buf, err := md.buf.AllocPage()
		if err != nil {
			md.metrics.allocatedPageFailures.Incr()
			return 0, err
		}
		md.metrics.allocatedPages.Incr()
		fStore = newFieldStore(buf, fieldID)
		writtenSize += fStore.Capacity()
		beforeTStoreSize := tStore.Capacity()
		tStore.InsertFStore(fStore)
		writtenSize += tStore.Capacity() - beforeTStoreSize
		// if write data success, add field into metric level for cache
		mStore.AddField(fieldID, fieldType)
	}
	beforeFStoreCapacity := fStore.Capacity()
	fStore.Write(fieldType, slotIndex, fieldValue)
	return writtenSize + fStore.Capacity() - beforeFStoreCapacity, nil
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder.
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher) error {
	// waiting current writing complete
	md.writeCondition.Wait()

	var flushCtx flushContext
	if err := md.mStores.WalkEntry(func(key uint32, value mStoreINTF) error {
		flushCtx.metricID = key
		if err := value.FlushMetricsDataTo(flusher, &flushContext{
			metricID: key,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return flusher.Close()
}

// Filter filters the data based on metric/seriesIDs,
// if it finds data then returns the flow.FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(executeCtx *flow.StorageExecuteContext) ([]flow.FilterResultSet, error) {
	md.rwMutex.RLock()
	defer md.rwMutex.RUnlock()

	mStore, ok := md.mStores.Get(uint32(executeCtx.MetricID))
	if !ok {
		return nil, nil
	}

	// TODO filter slot range
	return mStore.Filter(md.familyTime, executeCtx.SeriesIDsAfterFiltering, executeCtx.Fields)
}

// MemSize returns the time series database memory size
func (md *memoryDatabase) MemSize() int64 {
	return md.allocSize.Load()
}

// Close releases resources for current memory database.
func (md *memoryDatabase) Close() error {
	md.buf.Release()
	return nil
}

func (md *memoryDatabase) Uptime() time.Duration {
	return time.Duration(fasttime.UnixNano() - md.createdTime)
}

// Size returns the number of metric names.
func (md *memoryDatabase) Size() int {
	md.rwMutex.RLock()
	defer md.rwMutex.RUnlock()
	return md.mStores.Size()
}
