package tsdb

import (
	"sync"

	"github.com/lindb/lindb/series/field"
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

// db level
type MetricMetadata struct {
	metric sync.Map   // metric name => metric schema(Schema)
	lock   sync.Mutex // lock for create metric schema
}

func newMemoryDatabse() *MetricMetadata {
	return &MetricMetadata{}
}

func (s *MetricMetadata) GetSchema(metric string) *Schema {
	schema, ok := s.metric.Load(metric)
	if ok {
		return schema.(*Schema)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// double check
	schema, ok = s.metric.Load(metric)
	if ok {
		return schema.(*Schema)
	}

	// create new metric schema
	metricSchema := &Schema{}
	s.metric.Store(metric, metricSchema)
	return metricSchema
}

// shard level
type TimeSeries struct {
	forwardIDs sync.Map // time series hash => time series id

	seq uint32

	lock sync.Mutex
}

func (ts *TimeSeries) GetTimeSeriesID(tsHash uint64, createCallback func(newTimeSeriesID uint32)) uint32 {
	tsID, ok := ts.forwardIDs.Load(tsHash)
	if ok {
		return tsID.(uint32)
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	// double check
	tsID, ok = ts.forwardIDs.Load(tsHash)
	if ok {
		return tsID.(uint32)
	}

	// create new time series ID
	id := ts.seq
	ts.seq++
	ts.forwardIDs.Store(tsHash, id)

	// invoke create time series callback
	createCallback(id)

	return id
}

// shard level
type TimeSeriesIndex struct {
	index sync.Map // tag value id => bitmap(series id)
}

type MetadataDict struct {
	dict sync.Map // string => id
}
