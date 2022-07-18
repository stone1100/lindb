package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

type groupingStage struct {
	baseStage
	executeCtx *flow.DataLoadContext
	shard      tsdb.Shard
}

func NewGroupingStage(executeCtx *flow.DataLoadContext, shard tsdb.Shard) Stage {
	return &groupingStage{
		baseStage: baseStage{
			stageType: Grouping,
		},
		executeCtx: executeCtx,
		shard:      shard,
	}
}

func (stage *groupingStage) Plan() PlanNode {
	execPlan := NewRootPlanNode()
	// add find grouping node
	execPlan.AddChild(NewPlanNode(operator.NewGroupingFind(stage.executeCtx)))
	return execPlan
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
