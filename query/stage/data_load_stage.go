package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
	storagequery "github.com/lindb/lindb/query/storage"
)

type dataLoadStage struct {
	executeCtx *flow.DataLoadContext
	segmentRS  *flow.TimeSegmentResultSet
}

func NewDataLoadStage() Stage {
	return &dataLoadStage{}
}

func (stage *dataLoadStage) Plan() storagequery.PlanNode {
	execPlan := storagequery.NewRootPlanNode()
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
		execPlan.AddChild(storagequery.NewPlanNode(
			operator.NewDataLoader(stage.executeCtx, stage.segmentRS, stage.segmentRS.FilterRS[idx])))
	}
	return execPlan
}

func (stage *dataLoadStage) NextStages() []Stage {
	return nil
}
