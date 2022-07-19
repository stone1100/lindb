package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

type shardScanStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext

	shardExecuteCtx *flow.ShardExecuteContext
	shardID         models.ShardID

	shard tsdb.Shard
}

func NewShardScanStage(leafExecuteCtx *context.LeafExecuteContext,
	shardExecuteCtx *flow.ShardExecuteContext, shardID models.ShardID,
) Stage {
	leafExecuteCtx.GroupingCtx.ForkGroupingTask()
	return &shardScanStage{
		baseStage: baseStage{
			stageType: ShardScan,
		},
		leafExecuteCtx:  leafExecuteCtx,
		shardExecuteCtx: shardExecuteCtx,
		shardID:         shardID,
	}
}

func (stage *shardScanStage) Plan() PlanNode {
	execPlan := NewRootPlanNode()
	shardExecuteCtx := stage.shardExecuteCtx
	queryStmt := shardExecuteCtx.StorageExecuteCtx.Query
	// if shard exist, add shard to query list
	if shard, ok := stage.leafExecuteCtx.Database.GetShard(stage.shardID); ok {
		stage.shard = shard
		families := shard.GetDataFamilies(queryStmt.StorageInterval.Type(), queryStmt.TimeRange)
		if len(families) == 0 {
			// no data family found
			return execPlan
		}
		if queryStmt.Condition != nil {
			// add shard level series filtering node
			execPlan.AddChild(NewPlanNode(operator.NewSeriesFiltering(shardExecuteCtx, shard)))
		} else {
			// add shard level all series lookup node
			execPlan.AddChild(NewPlanNode(operator.NewMetricAllSeries(shardExecuteCtx, shard)))
		}

		for idx := range families {
			family := families[idx]
			// add data family reader node, found series ids which match condition.
			execPlan.AddChild(NewPlanNode(operator.NewDataFamilyReader(shardExecuteCtx, family)))
		}

		if shardExecuteCtx.StorageExecuteCtx.Query.HasGroupBy() {
			// if it has grouping, do group by tag keys, else just split series ids as batch first,
			// get grouping context if it needs
			// group context find task maybe change shardExecuteContext.SeriesIDsAfterFiltering value.
			execPlan.AddChild(NewPlanNode(operator.NewGroupingContextBuild(shardExecuteCtx, shard)))
		}
	}
	return execPlan
}

func (stage *shardScanStage) NextStages() (stages []Stage) {
	if stage.shard == nil {
		// shard not found
		return
	}
	// if not grouping found, series id is empty.
	shardExecuteContext := stage.shardExecuteCtx
	seriesIDs := shardExecuteContext.SeriesIDsAfterFiltering
	seriesIDsHighKeys := seriesIDs.GetHighKeys()

	for seriesIDHighKeyIdx := range seriesIDsHighKeys {
		// be carefully, need use new variable for variable scope problem(closures)
		// ref: https://go.dev/doc/faq#closures_and_goroutines
		highSeriesIDIdx := seriesIDHighKeyIdx
		// grouping based on group by tag keys for each low series container
		lowSeriesIDs := seriesIDs.GetContainerAtIndex(highSeriesIDIdx)
		dataLoadCtx := &flow.DataLoadContext{
			ShardExecuteCtx:       shardExecuteContext,
			LowSeriesIDsContainer: lowSeriesIDs,
			SeriesIDHighKey:       seriesIDsHighKeys[highSeriesIDIdx],
			IsMultiField:          len(shardExecuteContext.StorageExecuteCtx.Fields) > 1,
			IsGrouping:            shardExecuteContext.StorageExecuteCtx.Query.HasGroupBy(),
		}

		stages = append(stages, NewGroupingStage(stage.leafExecuteCtx, dataLoadCtx, stage.shard))
	}
	return stages
}

func (stage *shardScanStage) Complete() {
	stage.leafExecuteCtx.GroupingCtx.CompleteGroupingTask()
}
