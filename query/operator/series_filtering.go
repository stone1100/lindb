package operator

import (
	"fmt"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/roaring"
)

type seriesFiltering struct {
	executeCtx *flow.ShardExecuteContext
	indexDB    indexdb.IndexDatabase

	err error
}

func NewSeriesFiltering(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &seriesFiltering{
		executeCtx: executeCtx,
		indexDB:    shard.IndexDatabase(),
	}
}

func (op *seriesFiltering) Execute() error {
	queryStmt := op.executeCtx.StorageExecuteCtx.Query
	// if it gets tag filter result do series ids searching
	_, seriesIDs := op.findSeriesIDsByExpr(queryStmt.Condition)
	if op.err != nil {
		return op.err
	}
	op.executeCtx.SeriesIDsAfterFiltering.Or(seriesIDs)
	return nil
}

// findSeriesIDsByExpr finds series ids by expr, recursion filter for expr
func (op *seriesFiltering) findSeriesIDsByExpr(condition stmt.Expr) (tag.KeyID, *roaring.Bitmap) {
	if condition == nil {
		return 0, roaring.New() // create an empty series ids for parent expr
	}
	if op.err != nil {
		return 0, roaring.New() // create an empty series ids for parent expr
	}
	switch expr := condition.(type) {
	case stmt.TagFilter:
		tagKey, seriesIDs, err := op.getSeriesIDsByExpr(expr)
		if err != nil {
			op.err = err
			return tagKey, roaring.New() // create an empty series ids for parent expr
		}
		return tagKey, seriesIDs
	case *stmt.ParenExpr:
		return op.findSeriesIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// get filter series ids
		tagKey, matchResult := op.findSeriesIDsByExpr(expr.Expr)
		// get all series ids for tag key
		all, err := op.indexDB.GetSeriesIDsForTag(tagKey)
		if err != nil {
			op.err = err
			return tagKey, roaring.New() // create an empty series ids for parent expr
		}
		// do and not got series ids not in 'a' list
		all.AndNot(matchResult)
		return 0, all
	case *stmt.BinaryExpr:
		_, left := op.findSeriesIDsByExpr(expr.Left)
		_, right := op.findSeriesIDsByExpr(expr.Right)
		if expr.Operator == stmt.AND {
			left.And(right)
		} else {
			left.Or(right)
		}
		return 0, left
	}
	return 0, roaring.New() // create an empty series ids for parent expr
}

// getTagKeyID returns the tag key id by tag key
func (op *seriesFiltering) getSeriesIDsByExpr(expr stmt.Expr) (tag.KeyID, *roaring.Bitmap, error) {
	tagValues, ok := op.executeCtx.StorageExecuteCtx.TagFilterResult[expr.Rewrite()]
	if !ok {
		return 0, nil, fmt.Errorf("%w, expr: %s", constants.ErrTagValueFilterResultNotFound, expr.Rewrite())
	}
	seriesIDs, err := op.indexDB.GetSeriesIDsByTagValueIDs(tagValues.TagKeyID, tagValues.TagValueIDs)
	if err != nil {
		return 0, nil, err
	}
	return tagValues.TagKeyID, seriesIDs, nil
}
