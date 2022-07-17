package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

type metadataLookupStage struct {
	baseStage
	leafCtx *context.LeafExecuteContext
}

func NewMetadataLookupStage(leafCtx *context.LeafExecuteContext) Stage {
	return &metadataLookupStage{
		baseStage: baseStage{
			stageType: MetadataLookup,
		},
		leafCtx: leafCtx,
	}
}
func (stage *metadataLookupStage) Plan() PlanNode {
	execPlan := NewRootPlanNode()
	execCtx := stage.leafCtx.StorageExecuteContext
	database := stage.leafCtx.Database
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
	storageExecuteCtx := stage.leafCtx.StorageExecuteContext
	shardIDs := storageExecuteCtx.ShardIDs
	storageExecuteCtx.ShardContexts = make([]*flow.ShardExecuteContext, len(shardIDs))
	for shardIdx := range shardIDs {
		shardExecuteCtx := flow.NewShardExecuteContext(storageExecuteCtx)
		storageExecuteCtx.ShardContexts[shardIdx] = shardExecuteCtx
		stages = append(stages, NewShardScanStage(stage.leafCtx, shardExecuteCtx, shardIDs[shardIdx]))
	}
	return
}
