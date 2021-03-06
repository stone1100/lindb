package tblstore

import (
	"fmt"

	"github.com/lindb/lindb/series/field"

	"github.com/lindb/lindb/kv"
)

type metricsMetaMerger struct {
	flusher      *metricsMetaFlusher
	reader       *metricsMetaReader
	nopKVFlusher *kv.NopFlusher
	fieldMetas   []field.Meta
}

// NewMetricsMetaMerger returns a new merger for compacting MetricsMetaTable
func NewMetricsMetaMerger() kv.Merger {
	m := &metricsMetaMerger{
		reader:       NewMetricsMetaReader(nil).(*metricsMetaReader),
		nopKVFlusher: kv.NewNopFlusher(),
	}
	m.flusher = NewMetricsMetaFlusher(m.nopKVFlusher).(*metricsMetaFlusher)
	return m
}

func (m *metricsMetaMerger) Merge(
	key uint32,
	value [][]byte,
) (
	[]byte,
	error,
) {
	var hasData bool
	defer func() {
		m.flusher.Reset()
		m.fieldMetas = m.fieldMetas[:0]
	}()
	// flush tag-key
	for _, block := range value {
		tagMetaBlock, fieldMetaBlock := m.reader.readMetasBlock(block)
		tagItr := newTagKeyIDIterator(tagMetaBlock)
		for tagItr.HasNext() {
			hasData = true
			m.flusher.FlushTagKeyID(tagItr.Next())
		}
		fieldItr := newFieldMetaIterator(fieldMetaBlock)
		for fieldItr.HasNext() {
			hasData = true
			m.fieldMetas = append(m.fieldMetas, fieldItr.Next())
		}
	}
	// flush field-meta
	for _, fm := range m.fieldMetas {
		m.flusher.FlushFieldMeta(fm)
	}
	if !hasData {
		return nil, fmt.Errorf("no available blocks for compacting")
	}
	_ = m.flusher.FlushMetricMeta(key)
	return m.nopKVFlusher.Bytes(), nil
}
