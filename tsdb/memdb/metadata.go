package memdb

import (
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

type Schema struct {
	fields  sync.Map // field name => field meta
	tagKeys tag.Metas

	fieldIdx uint8
	lock     sync.Mutex // lock for create tags/fields(like table columns)
}

func (s *Schema) GetFieldID(fieldName field.Name, fieldType field.Type) field.ID {
	fieldMeta, ok := s.fields.Load(fieldName)
	if ok {
		return fieldMeta.(*field.Meta).ID
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	// double check
	fieldMeta, ok = s.fields.Load(fieldName)
	if ok {
		return fieldMeta.(*field.Meta).ID
	}

	fieldID := field.ID(s.fieldIdx)
	s.fieldIdx++
	s.fields.Store(fieldName, &field.Meta{
		ID:   fieldID,
		Name: fieldName,
		Type: fieldType,
	})
	return fieldID
}

func (s *Schema) GetTagMeta(tagKey string) *TagMeta {
	return nil
}

// db level
type MetricMetadata struct {
	metric sync.Map // metric name => metric schema(Schema)

	dict sync.Map

	lock sync.Mutex // lock for create metric schema
}

func NewMetricMetadata() *MetricMetadata {
	return &MetricMetadata{}
}

func (mm *MetricMetadata) GetID(val string) uint32 {
	return 0
}

func (mm *MetricMetadata) GetSchema(metric string) *Schema {
	schema, ok := mm.metric.Load(metric)
	if ok {
		return schema.(*Schema)
	}

	mm.lock.Lock()
	defer mm.lock.Unlock()

	// double check
	schema, ok = mm.metric.Load(metric)
	if ok {
		return schema.(*Schema)
	}

	// create new metric schema
	metricSchema := &Schema{}
	mm.metric.Store(metric, metricSchema)
	return metricSchema
}

type TagMeta struct {
	id tag.KeyID

	tagValueIDs *roaring.Bitmap
}

func (tm *TagMeta) GetTagValueID(val string) uint32 {
	return 0
}

// shard level
type TimeSeries struct {
	forwardIndex sync.Map // time series hash => time series id(global id mapping)

	invertedIndex map[string]map[string]*roaring.Bitmap // tag key => tag value => bitmap(time series id)
	metricIndex   sync.Map                              // metric id => bitmap(time series id)

	filter *bloom.BloomFilter

	seq uint32

	lock sync.Mutex
}

func (ts *TimeSeries) GetTimeSeriesID(tsHash uint64, createCallback func(newTimeSeriesID uint32)) uint32 {
	// TODO: check from file
	// exist := ts.filter.Test([]byte{})
	// if !exist {
	//
	// }

	tsID, ok := ts.forwardIndex.Load(tsHash)
	if ok {
		return tsID.(uint32)
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	// double check
	tsID, ok = ts.forwardIndex.Load(tsHash)
	if ok {
		return tsID.(uint32)
	}

	// create new time series ID
	id := ts.seq
	ts.seq++
	ts.forwardIndex.Store(tsHash, id)

	// invoke create time series callback
	createCallback(id)

	return id
}

func (ts *TimeSeries) BuildIndex(schema *Schema, row *metric.StorageRow) {
	_ = ts.GetTimeSeriesID(row.TagsHash(), func(newTimeSeriesID uint32) {
		it := row.NewKeyValueIterator()
		for it.HasNext() {
			tagKey := it.NextKey()
			tagValue := it.NextValue()
			index, ok := ts.invertedIndex[string(tagKey)]

			if !ok {
				// create new series ids for new tag value
				seriesIDs := roaring.BitmapOf(newTimeSeriesID)
				ts.invertedIndex[string(tagKey)] = map[string]*roaring.Bitmap{
					string(tagValue): seriesIDs,
				}
			} else {
				seriesIDs, ok := index[string(tagValue)]
				if !ok {
					index[string(tagValue)] = roaring.BitmapOf(newTimeSeriesID)
				} else {
					seriesIDs.Add(newTimeSeriesID)
				}
			}
		}

		// FIXME: no tag
	})
}
