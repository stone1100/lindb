package tsdb

import (
	"sync"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

type Schema struct {
	fields  map[field.Name]field.Meta
	field   sync.Map
	tagKeys tag.Metas

	c    uint8
	lock sync.RWMutex
}

func (s *Schema) GetFieldID(fieldName field.Name) field.ID {
	f, ok := s.field.Load(fieldName)
	if !ok {
		s.lock.Lock()
		defer s.lock.Unlock()
		v := field.ID(s.c)
		s.c++
		s.field.Store(fieldName, field.Meta{
			ID:   v,
			Name: fieldName,
		})
		return v
	}
	return f.(field.Meta).ID
}

type memoryDatabase struct {
	// metric map[string]*Schema
	metric sync.Map

	lock sync.RWMutex
}

func newMemoryDatabse() *memoryDatabase {
	return &memoryDatabase{
		// metric: make(map[string]*Schema),
	}
}

func (s *shard) GetSeriesID(tags uint64) int32 {
	s.lock.Lock()
	defer s.lock.Unlock()

	id, ok := s.ids[tags]
	if ok {
		return id
	}
	id = s.seq.Inc()
	s.ids[tags] = id
	return id
}

func (s *memoryDatabase) GetFieldID(metric string, fieldName field.Name) field.ID {
	schema, ok := s.metric.Load(metric)
	if !ok {
		s.lock.Lock()
		schema = &Schema{}
		s.metric.Store(metric, schema)

		s.lock.Unlock()
	}

	return (schema.(*Schema)).GetFieldID(fieldName)
}

func (s *shard) GetFieldID(metric string, fieldName field.Name) field.ID {
	s.lock.Lock()
	defer s.lock.Unlock()

	schema, ok := s.metric[metric]
	if !ok {
		schema = &Schema{
			fields: make(map[field.Name]field.Meta),
		}
		s.metric[metric] = schema
		l := field.ID(len(schema.fields))
		schema.fields[fieldName] = field.Meta{
			ID:   l,
			Name: fieldName,
		}
		return l
	}
	m, ok := schema.fields[fieldName]
	if ok {
		return m.ID
	}
	l := field.ID(len(schema.fields))
	schema.fields[fieldName] = field.Meta{
		ID:   l,
		Name: fieldName,
	}
	return l
}
