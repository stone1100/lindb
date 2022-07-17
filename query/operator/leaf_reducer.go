package operator

import (
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
)

type leafReducer struct {
	executeCtx *flow.DataLoadContext
	reduceAgg  aggregation.GroupingAggregator
	lock       sync.Mutex
}

func NewLeafReducer(executeCtx *flow.DataLoadContext) Operator {
	storageExecuteCtx := executeCtx.ShardExecuteCtx.StorageExecuteCtx
	aggregatorSpecs := storageExecuteCtx.AggregatorSpecs
	reduceAgg := aggregation.NewGroupingAggregator(storageExecuteCtx.Query.Interval,
		storageExecuteCtx.Query.IntervalRatio, storageExecuteCtx.Query.TimeRange, aggregatorSpecs)
	return &leafReducer{
		reduceAgg:  reduceAgg,
		executeCtx: executeCtx,
	}
}

// Reduce reduces the down sampling aggregator's result.
func (op *leafReducer) Reduce(it series.GroupedIterator) {
	op.lock.Lock()
	defer op.lock.Unlock()

	op.reduceAgg.Aggregate(it)
}

func (op *leafReducer) Execute() error {
	// after load, need to reduce the aggregator's result to query flow.
	op.executeCtx.Reduce(op.Reduce)
	return nil
}
