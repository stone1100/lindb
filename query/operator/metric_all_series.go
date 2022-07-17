package operator

import (
	"fmt"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

type metricAllSeries struct {
	executeCtx *flow.ShardExecuteContext
	indexDB    indexdb.IndexDatabase

	err error
}

func NewMetricAllSeries(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &metricAllSeries{
		executeCtx: executeCtx,
		indexDB:    shard.IndexDatabase(),
	}
}

func (op *metricAllSeries) Execute() error {
	queryStmt := op.executeCtx.StorageExecuteCtx.Query
	// get series ids for metric level
	seriesIDs, err := op.indexDB.GetSeriesIDsForMetric(queryStmt.Namespace, queryStmt.MetricName)
	if err != nil {
		return err
	}
	if !queryStmt.HasGroupBy() {
		// add series id without tags, maybe metric has too many series, but one series without tags
		seriesIDs.Add(series.IDWithoutTags)
	}
	fmt.Println(seriesIDs.ToArray())
	op.executeCtx.SeriesIDsAfterFiltering.Or(seriesIDs)
	return nil
}
