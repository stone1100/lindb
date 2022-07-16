package operator

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
)

type groupingContextBuild struct {
	executeCtx *flow.ShardExecuteContext
	shard      tsdb.Shard
}

func NewGroupingContextBuild(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &groupingContextBuild{
		executeCtx: executeCtx,
		shard:      shard,
	}
}

func (op *groupingContextBuild) Execute() error {
	return op.shard.IndexDatabase().GetGroupingContext(op.executeCtx)
}
