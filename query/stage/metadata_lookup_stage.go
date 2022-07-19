package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

type metadataLookupStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext
}

func NewMetadataLookupStage(leafExecuteCtx *context.LeafExecuteContext) Stage {
	return &metadataLookupStage{
		baseStage: baseStage{
			stageType: MetadataLookup,
		},
		leafExecuteCtx: leafExecuteCtx,
	}
}
func (stage *metadataLookupStage) Plan() PlanNode {
	execPlan := NewRootPlanNode()
	execCtx := stage.leafExecuteCtx.StorageExecuteCtx
	database := stage.leafExecuteCtx.Database
	// add metadata lookup(name/tag/field etc.) node
	execPlan.AddChild(NewPlanNode(operator.NewMetadataLookup(execCtx, database)))
	hasWhereCondition := execCtx.Query.Condition != nil
	if hasWhereCondition {
		// add tag values lookup node if query has where condition
		execPlan.AddChild(NewPlanNode(operator.NewTagValuesLookup(execCtx, database)))
	}
	return execPlan
}

func (stage *metadataLookupStage) NextStages() (stages []Stage) {
	storageExecuteCtx := stage.leafExecuteCtx.StorageExecuteCtx
	shardIDs := storageExecuteCtx.ShardIDs
	storageExecuteCtx.ShardContexts = make([]*flow.ShardExecuteContext, len(shardIDs))
	for shardIdx := range shardIDs {
		shardExecuteCtx := flow.NewShardExecuteContext(storageExecuteCtx)
		storageExecuteCtx.ShardContexts[shardIdx] = shardExecuteCtx
		stages = append(stages, NewShardScanStage(stage.leafExecuteCtx, shardExecuteCtx, shardIDs[shardIdx]))
	}
	return
}
