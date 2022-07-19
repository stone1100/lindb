package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

type groupingStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext
	executeCtx     *flow.DataLoadContext
	shard          tsdb.Shard
}

func NewGroupingStage(leafExecuteCtx *context.LeafExecuteContext, executeCtx *flow.DataLoadContext, shard tsdb.Shard) Stage {
	leafExecuteCtx.GroupingCtx.ForkGroupingTask()
	return &groupingStage{
		baseStage: baseStage{
			stageType: Grouping,
		},
		leafExecuteCtx: leafExecuteCtx,
		executeCtx:     executeCtx,
		shard:          shard,
	}
}

func (stage *groupingStage) Plan() PlanNode {
	// add find grouping node
	return NewPlanNode(operator.NewGroupingFind(stage.executeCtx))
}

func (stage *groupingStage) NextStages() (stages []Stage) {
	if stage.executeCtx.IsGrouping && len(stage.executeCtx.GroupingSeriesAgg) == 0 {
		// if not found any grouping tags, terminal.
		return
	}
	// time segments sorted by family time
	timeSegments := stage.executeCtx.ShardExecuteCtx.TimeSegmentContext.GetTimeSegments()
	for segmentIdx := range timeSegments {
		// add data load stage based on time segment, one by one
		stages = append(stages, NewDataLoadStage(stage.executeCtx, timeSegments[segmentIdx]))
	}
	return
}

func (stage *groupingStage) Complete() {
	stage.leafExecuteCtx.GroupingCtx.CompleteGroupingTask()
}
