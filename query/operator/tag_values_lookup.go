package operator

import (
	"fmt"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

type tagValuesLookup struct {
	executeCtx *flow.StorageExecuteContext
	metadata   metadb.Metadata

	err error
}

func NewTagValuesLookup(executeCtx *flow.StorageExecuteContext, database tsdb.Database) Operator {
	return &tagValuesLookup{
		executeCtx: executeCtx,
		metadata:   database.Metadata(),
	}
}

func (op *tagValuesLookup) Execute() error {
	op.executeCtx.TagFilterResult = make(map[string]*flow.TagFilterResult)
	op.findTagValueIDsByExpr(op.executeCtx.Query.Condition)
	return op.err
}

// findTagValueIDsByExpr finds tag value ids by expr, recursion filter for expr
func (op *tagValuesLookup) findTagValueIDsByExpr(expr stmt.Expr) {
	if expr == nil {
		return
	}
	if op.err != nil {
		return
	}
	switch expr := expr.(type) {
	case stmt.TagFilter:
		tagKeyID, err := op.getTagKeyID(expr.TagKey())
		if err != nil {
			op.err = err
			return
		}
		tagValueIDs, err := op.metadata.TagMetadata().FindTagValueDsByExpr(tagKeyID, expr)
		if err != nil {
			op.err = err
			return
		}
		if tagValueIDs != nil && !tagValueIDs.IsEmpty() {
			// save atomic tag filter result
			op.executeCtx.TagFilterResult[expr.Rewrite()] = &flow.TagFilterResult{
				TagKeyID:    tagKeyID,
				TagValueIDs: tagValueIDs,
			}
		}
	case *stmt.ParenExpr:
		op.findTagValueIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// find tag value id by expr => (not tag filter) => tag filter
		op.findTagValueIDsByExpr(expr.Expr)
	case *stmt.BinaryExpr:
		if expr.Operator != stmt.AND && expr.Operator != stmt.OR {
			op.err = fmt.Errorf("wrong binary operator in tag filter: %s", stmt.BinaryOPString(expr.Operator))
			return
		}
		op.findTagValueIDsByExpr(expr.Left)
		op.findTagValueIDsByExpr(expr.Right)
	}
}

// getTagKeyID returns the tag key id by tag key
func (op *tagValuesLookup) getTagKeyID(tagKey string) (tag.KeyID, error) {
	// try to get tag key from context
	if tagKeyID, ok := op.executeCtx.TagKeys[tagKey]; ok {
		return tagKeyID, nil
	}
	queryStmt := op.executeCtx.Query
	tagKeyID, err := op.metadata.MetadataDatabase().GetTagKeyID(queryStmt.Namespace, queryStmt.MetricName, tagKey)
	if err != nil {
		return 0, err
	}
	op.executeCtx.TagKeys[tagKey] = tagKeyID
	return tagKeyID, nil
}
