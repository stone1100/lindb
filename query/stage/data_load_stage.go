package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
)

type dataLoadStage struct {
	baseStage
	executeCtx *flow.DataLoadContext
	segmentRS  *flow.TimeSegmentResultSet
}

func NewDataLoadStage(executeCtx *flow.DataLoadContext, segmentRS *flow.TimeSegmentResultSet) Stage {
	return &dataLoadStage{
		baseStage: baseStage{
			stageType: DataLoad,
		},
		executeCtx: executeCtx,
		segmentRS:  segmentRS,
	}
}

func (stage *dataLoadStage) Plan() PlanNode {
	execPlan := NewRootPlanNode()
	shardExecuteCtx := stage.executeCtx.ShardExecuteCtx
	// add segment data load node.
	stage.segmentRS.IntervalRatio = uint16(shardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio)
	// calc base slot based on query interval and family time of storage
	queryInterval := shardExecuteCtx.StorageExecuteCtx.Query.Interval
	calc := queryInterval.Calculator()
	familyTimeForQuery := calc.CalcFamilyTime(stage.segmentRS.FamilyTime)
	stage.segmentRS.BaseTime = uint16(calc.CalcSlot(stage.segmentRS.FamilyTime, familyTimeForQuery, queryInterval.Int64()))
	stage.segmentRS.TargetRange = shardExecuteCtx.StorageExecuteCtx.CalcTargetSlotRange(familyTimeForQuery)

	for idx := range stage.segmentRS.FilterRS {
		execPlan.AddChild(NewPlanNode(
			operator.NewDataLoader(stage.executeCtx, stage.segmentRS, stage.segmentRS.FilterRS[idx])))
		execPlan.AddChild(NewPlanNode(operator.NewLeafReducer(stage.executeCtx)))
	}
	return execPlan
}

func (stage *dataLoadStage) NextStages() []Stage {
	// terminal stage, return empty stages.
	return nil
}
