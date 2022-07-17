package operator

//
//type segmentLoad struct {
//	executeCtx *flow.DataLoadContext
//	segmentCtx *flow.TimeSegmentResultSet
//}
//
//func NewSegmentLoad() Operator {
//	return &segmentLoad{}
//}
//
//func (op *segmentLoad) Execute() error {
//	queryIntervalRatio := op.executeCtx.ShardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio
//	seriesIDs := op.executeCtx.ShardExecuteCtx.SeriesIDsAfterFiltering // after group result
//	// calc base slot based on query interval and family time of storage
//	queryInterval := op.executeCtx.ShardExecuteCtx.StorageExecuteCtx.Query.Interval
//	calc := queryInterval.Calculator()
//	familyTimeForQuery := calc.CalcFamilyTime(op.segmentCtx.FamilyTime)
//	baseSlot := uint16(calc.CalcSlot(op.segmentCtx.FamilyTime, familyTimeForQuery, queryInterval.Int64()))
//	targetSlotRange := op.executeCtx.ShardExecuteCtx.StorageExecuteCtx.CalcTargetSlotRange(familyTimeForQuery)
//
//	for idx, rs := range op.segmentCtx.FilterRS {
//		// double filtering, maybe some series ids be filtered out when do grouping.
//		// filter logic: forward_reader.go -> GetGroupingScanner
//		if roaring.FastAnd(seriesIDs, rs.SeriesIDs()).IsEmpty() {
//			continue
//		}
//		// maybe return nil loader
//		loader := rs.Load(op.executeCtx)
//		if loader == nil {
//			continue
//		}
//
//		// load field series data by series ids
//		op.executeCtx.Decoder = encoding.GetTSDDecoder()
//		op.executeCtx.DownSampling = func(slotRange timeutil.SlotRange, lowSeriesIdx uint16, fieldIdx int, getter encoding.TSDValueGetter) {
//			var agg aggregation.FieldAggregator
//			seriesAggregator := op.executeCtx.GetSeriesAggregator(lowSeriesIdx, fieldIdx)
//
//			var ok bool
//			agg, ok = seriesAggregator.GetAggregator(op.segmentCtx.FamilyTime)
//			if !ok {
//				return
//			}
//			aggregation.DownSampling(
//				slotRange, targetSlotRange, uint16(queryIntervalRatio), baseSlot,
//				getter,
//				agg.AggregateBySlot,
//			)
//		}
//
//		// loads the metric data by given series id from load result.
//		// if found data need to do down sampling aggregate.
//		loader.Load(op.executeCtx)
//		// release tsd decoder back to pool for re-use.
//		encoding.ReleaseTSDDecoder(op.executeCtx.Decoder)
//		// after load, need to reduce the aggregator's result to query flow.
//		op.executeCtx.Reduce(t.queryFlow.Reduce)
//	}
//	return nil
//}
