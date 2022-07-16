package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/operator"
	storagequery "github.com/lindb/lindb/query/storage"
	"github.com/lindb/lindb/tsdb"
)

type shardScanStage struct {
	executeCtx *flow.ShardExecuteContext
	shardID    models.ShardID

	database tsdb.Database
}

func NewShardScanStage() Stage {
	return &shardScanStage{}
}

func (stage *shardScanStage) Plan() storagequery.PlanNode {
	execPlan := storagequery.NewRootPlanNode()
	executeCtx := stage.executeCtx
	queryStmt := executeCtx.StorageExecuteCtx.Query
	// if shard exist, add shard to query list
	if shard, ok := stage.database.GetShard(stage.shardID); ok {
		families := shard.GetDataFamilies(queryStmt.StorageInterval.Type(), queryStmt.TimeRange)
		if len(families) == 0 {
			// no data family found
			return execPlan
		}
		if queryStmt.Condition != nil {
			// add shard level series filtering node
			execPlan.AddChild(storagequery.NewPlanNode(operator.NewSeriesFiltering(executeCtx, shard)))
		} else {
			// add shard level all series lookup node
			execPlan.AddChild(storagequery.NewPlanNode(operator.NewMetricAllSeries(executeCtx, shard)))
		}

		for idx := range families {
			family := families[idx]
			// add data family reader node, found series ids which match condition.
			execPlan.AddChild(storagequery.NewPlanNode(operator.NewDataFamilyReader(executeCtx, family)))
		}

		if stage.executeCtx.StorageExecuteCtx.Query.HasGroupBy() {
			// 1. grouping, if it has grouping, do group by tag keys, else just split series ids as batch first,
			// get grouping context if it needs
			// group context find task maybe change shardExecuteContext.SeriesIDsAfterFiltering value.
			execPlan.AddChild(storagequery.NewPlanNode(operator.NewGroupingContextBuild(stage.executeCtx, shard)))
		}
	}
	return execPlan
}

func (stage *shardScanStage) NextStages() (stages []Stage) {
	// if not grouping found, series id is empty.
	seriesIDs := stage.executeCtx.SeriesIDsAfterFiltering
	seriesIDsHighKeys := seriesIDs.GetHighKeys()
	shardExecuteContext := stage.executeCtx

	for seriesIDHighKeyIdx := range seriesIDsHighKeys {
		// be carefully, need use new variable for variable scope problem(closures)
		// ref: https://go.dev/doc/faq#closures_and_goroutines
		highSeriesIDIdx := seriesIDHighKeyIdx
		// grouping based on group by tag keys for each low series container
		//TODO
		stages = append(stages, NewGroupingStage())

	}
	return stages
}
