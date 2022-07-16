package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
	storagequery "github.com/lindb/lindb/query/storage"
	"github.com/lindb/lindb/tsdb"
)

type metadataStage struct {
	executeCtx *flow.StorageExecuteContext
	database   tsdb.Database
}

func NewMetadataStage() Stage {
	return &metadataStage{}
}
func (stage *metadataStage) Plan() storagequery.PlanNode {
	execPlan := storagequery.NewRootPlanNode()
	executeCtx := stage.executeCtx
	// add metadata lookup(name/tag/field etc.) node
	execPlan.AddChild(storagequery.NewPlanNode(operator.NewMetadataLookup(executeCtx, stage.database)))
	hasWhereCondition := executeCtx.Query.Condition != nil
	if hasWhereCondition {
		// add tag values lookup node
		execPlan.AddChild(storagequery.NewPlanNode(operator.NewTagValuesLookup(executeCtx, stage.database)))
	}
	return execPlan
}

func (stage *metadataStage) NextStages() []Stage {
	//TODO implement me
	panic("implement me")
}
