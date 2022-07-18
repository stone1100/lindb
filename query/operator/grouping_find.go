package operator

import "github.com/lindb/lindb/flow"

type groupingFind struct {
	executeCtx *flow.DataLoadContext
}

func NewGroupingFind(executeCtx *flow.DataLoadContext) Operator {
	return &groupingFind{
		executeCtx: executeCtx,
	}
}

func (op *groupingFind) Execute() error {
	op.executeCtx.Grouping()
	if op.executeCtx.ShardExecuteCtx.GroupingContext != nil {
		// build group by data, grouped series: tags => series IDs(based on low series ids)
		op.executeCtx.ShardExecuteCtx.GroupingContext.BuildGroup(op.executeCtx)
	} else {
		op.executeCtx.PrepareAggregatorWithoutGrouping()
	}
	return nil
}
