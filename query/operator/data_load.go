package operator

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

type dataLoader struct {
	executeCtx *flow.DataLoadContext
	segmentRS  *flow.TimeSegmentResultSet
	rs         flow.FilterResultSet
}

func NewDataLoader(executeCtx *flow.DataLoadContext,
	segmentRS *flow.TimeSegmentResultSet, rs flow.FilterResultSet,
) Operator {
	return &dataLoader{
		executeCtx: executeCtx,
		segmentRS:  segmentRS,
		rs:         rs,
	}
}

func (op *dataLoader) Execute() error {
	seriesIDs := op.executeCtx.ShardExecuteCtx.SeriesIDsAfterFiltering // after group result
	// double filtering, maybe some series ids be filtered out when do grouping.
	// filter logic: forward_reader.go -> GetGroupingScanner
	if roaring.FastAnd(seriesIDs, op.rs.SeriesIDs()).IsEmpty() {
		return nil
	}
	loader := op.rs.Load(op.executeCtx)
	if loader == nil {
		// maybe return nil loader
		return nil
	}

	familyTime := op.segmentRS.FamilyTime
	targetSlotRange := op.segmentRS.TargetRange
	queryIntervalRatio := op.segmentRS.IntervalRatio
	baseSlot := op.segmentRS.BaseTime

	// load field series data by series ids
	op.executeCtx.Decoder = encoding.GetTSDDecoder()
	op.executeCtx.DownSampling = func(slotRange timeutil.SlotRange, lowSeriesIdx uint16, fieldIdx int, getter encoding.TSDValueGetter) {
		var agg aggregation.FieldAggregator
		seriesAggregator := op.executeCtx.GetSeriesAggregator(lowSeriesIdx, fieldIdx)

		var ok bool
		agg, ok = seriesAggregator.GetAggregator(familyTime)
		if !ok {
			return
		}
		aggregation.DownSampling(
			slotRange, targetSlotRange, queryIntervalRatio, baseSlot,
			getter,
			agg.AggregateBySlot,
		)
	}

	// loads the metric data by given series id from load result.
	// if found data need to do down sampling aggregate.
	loader.Load(op.executeCtx)
	// release tsd decoder back to pool for re-use.
	encoding.ReleaseTSDDecoder(op.executeCtx.Decoder)
	// after load, need to reduce the aggregator's result to query flow.
	op.executeCtx.Reduce(t.queryFlow.Reduce)
	return nil
}
