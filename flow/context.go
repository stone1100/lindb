package flow

import (
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/roaring"
)

// TagFilterResult represents the tag filter result, include tag key id and tag value ids.
type TagFilterResult struct {
	TagKeyID    tag.KeyID
	TagValueIDs *roaring.Bitmap
}

type StorageExecuteContext struct {
	// user input params
	Query    *stmt.Query
	MetricID metric.ID
	Fields   field.Metas

	SeriesIDsAfterFiltering *roaring.Bitmap
	// result which after tag condition metadata filter
	TagFilterResult             map[string]*TagFilterResult
	GroupByTagKeys              []tag.Meta
	TagValueIDsForGroupByTagKey []*roaring.Bitmap
}

func NewStorageExecuteContext(query *stmt.Query, metricID metric.ID, fields field.Metas) *StorageExecuteContext {
	return &StorageExecuteContext{
		Query:                   query,
		MetricID:                metricID,
		Fields:                  fields,
		SeriesIDsAfterFiltering: roaring.New(),
	}
}
