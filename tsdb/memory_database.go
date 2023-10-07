package tsdb

import (
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

type Schema struct {
	fields  field.Metas
	tagKeys tag.Metas
}

type metricMetadata struct {
}

type memoryDatabase struct {
}
