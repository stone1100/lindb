package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
	storagequery "github.com/lindb/lindb/query/storage"
	"github.com/lindb/lindb/tsdb"
)

type groupingStage struct {
	executeCtx *flow.DataLoadContext
	shard      tsdb.Shard
}

func NewGroupingStage() Stage {
	return &groupingStage{}
}

func (stage *groupingStage) Plan() storagequery.PlanNode {
	execPlan := storagequery.NewRootPlanNode()
	// add find grouping node
	execPlan.AddChild(storagequery.NewPlanNode(operator.NewGroupingFind(stage.executeCtx)))
	return execPlan
}

func (stage *groupingStage) NextStages() (stages []Stage) {
	if !stage.executeCtx.HasGroupingData() {
		// time segments sorted by family
		timeSegments := stage.executeCtx.ShardExecuteCtx.TimeSegmentContext.GetTimeSegments()
		for segmentIdx := range timeSegments {
			// add data load stage based on time segment, one by one
			stages = append(stages, NewDataLoadStage())
		}
	}
	return
}
