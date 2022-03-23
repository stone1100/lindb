// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package storagequery

import (
	"errors"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

// for testing
var (
	newTagSearchFunc          = newTagSearch
	newStorageExecutePlanFunc = newStorageExecutePlan
	newSeriesSearchFunc       = newSeriesSearch
	newBuildGroupTaskFunc     = newBuildGroupTask
	newDataLoadTaskFunc       = newDataLoadTask
)

var (
	errNoShardID         = errors.New("there is no shard id in search condition")
	errNoShardInDatabase = errors.New("there is no shard in database storage engine")
	errShardNotFound     = errors.New("shard not found in database storage engine")
	errShardNumNotMatch  = errors.New("got shard size not equals input shard size")
)

type timeSpanCtx struct {
	decoders [][]*encoding.TSDDecoder
	loaders  []flow.DataLoader // item maybe DataLoader is nil
}

// groupingResult represents the grouping context result
type groupingResult struct {
	groupingCtx series.GroupingContext
}

// groupedSeriesResult represents grouped series for group by query
type groupedSeriesResult struct {
	// tag values(ids) => low series ids.
	// if no grouping tag values is empty string value.
	groupedSeries map[string][]uint16
}

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Grouping if it needs
// 3) Scanning and Loading
// 4) Down sampling
// 5) Simple aggregation
type storageExecutor struct {
	database tsdb.Database
	ctx      *storageExecuteContext
	shards   []tsdb.Shard

	metricID           metric.ID
	fields             field.Metas
	storageExecutePlan *storageExecutePlan

	queryFlow flow.StorageQueryFlow

	queryTimeRange     timeutil.TimeRange
	queryInterval      timeutil.Interval
	queryIntervalRatio int

	// group by query need
	mutex              sync.Mutex
	groupByTagKeyIDs   []tag.Meta
	tagValueIDs        []*roaring.Bitmap // for group by query store tag value ids for each group tag key
	pendingForShard    atomic.Int32
	pendingForGrouping atomic.Int32
	collecting         atomic.Bool
}

// newStorageMetricQuery creates the execution which queries the data of storage engine
func newStorageMetricQuery(
	queryFlow flow.StorageQueryFlow,
	database tsdb.Database,
	storageExecuteCtx StorageExecuteContext,
) storageMetricQuery {
	ctx := storageExecuteCtx.(*storageExecuteContext)
	return &storageExecutor{
		database:  database,
		ctx:       ctx,
		queryFlow: queryFlow,
	}
}

// Execute executes search logic in storage level,
// 1) validation input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *storageExecutor) Execute() {
	// do query validation
	if err := e.validation(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	// get shard by given query shard id list
	for _, shardID := range e.ctx.shardIDs {
		shard, ok := e.database.GetShard(shardID)
		// if shard exist, add shard to query list
		if ok {
			e.shards = append(e.shards, shard)
		}
	}

	// check got shards if valid
	if err := e.checkShards(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	plan := newStorageExecutePlanFunc(e.ctx.query.Namespace, e.database.Metadata(), e.ctx.query)
	t := newStoragePlanTask(e.ctx, plan)

	if err := t.Run(); err != nil {
		e.queryFlow.Complete(err)
		return
	}
	condition := e.ctx.query.Condition
	if condition != nil {
		tagSearch := newTagSearchFunc(e.ctx.query.Namespace, e.ctx.query.MetricName,
			e.ctx.query.Condition, e.database.Metadata())
		t = newTagFilterTask(e.ctx, tagSearch)
		if err := t.Run(); err != nil {
			e.queryFlow.Complete(err)
			return
		}
	}

	e.metricID = plan.metricID
	e.fields = plan.getFields()
	e.storageExecutePlan = plan
	if e.ctx.query.HasGroupBy() {
		groupByTagKeyIDs := e.storageExecutePlan.groupByKeyIDs()
		groupByTagValueIDs := make([]*roaring.Bitmap, len(e.groupByTagKeyIDs))

		for idx, tagMeta := range groupByTagKeyIDs {
			tagValueIDs, err := e.database.Metadata().TagMetadata().GetTagValueIDsForTag(tagMeta.ID)
			if err != nil {
				e.queryFlow.Complete(err)
				return
			}
			groupByTagValueIDs[idx] = tagValueIDs
		}
		e.groupByTagKeyIDs = groupByTagKeyIDs
		e.tagValueIDs = groupByTagValueIDs
	}

	option := e.database.GetOption()
	// TODO need get storage interval by query time if has rollup config
	interval := option.Intervals[0].Interval
	e.queryTimeRange, e.queryIntervalRatio, e.queryInterval = downSamplingTimeRange(
		e.ctx.query.Interval, interval, e.ctx.query.TimeRange)

	// prepare storage query flow
	e.queryFlow.Prepare(e.queryInterval, e.queryIntervalRatio, e.queryTimeRange, plan.getAggregatorSpecs())

	// execute query flow
	e.executeQuery()
}

// executeQuery executes query flow for each shard
func (e *storageExecutor) executeQuery() {
	e.pendingForShard.Store(int32(len(e.shards)))
	for idx := range e.shards {
		shard := e.shards[idx]
		e.queryFlow.Filtering(func() {
			defer func() {
				// finish shard query
				e.pendingForShard.Dec()
				// try start collect tag values
				e.collectGroupByTagValues()
			}()
			// 1. get series ids by query condition
			executeCtx := flow.NewStorageExecuteContext(e.ctx.query, e.metricID, e.fields)
			executeCtx.TagFilterResult = e.ctx.tagFilterResult
			t := newSeriesIDsSearchTask(executeCtx, shard)
			err := t.Run()
			if err != nil && !errors.Is(err, constants.ErrNotFound) {
				// maybe series ids not found in shard, so ignore not found err
				e.queryFlow.Complete(err)
			}
			// if series ids not found
			if executeCtx.SeriesIDsAfterFiltering.IsEmpty() {
				return
			}

			rs := newTimeSpanResultSet(len(e.fields))
			// 2. filter data each data family in shard
			t = newFamilyFilterTask(executeCtx, shard, rs)
			err = t.Run()
			if err != nil && !errors.Is(err, constants.ErrNotFound) {
				// maybe data not exist in shard, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			if rs.isEmpty() {
				// data not found
				return
			}

			executeCtx.GroupByTagKeys = e.groupByTagKeyIDs
			executeCtx.TagValueIDsForGroupByTagKey = e.tagValueIDs

			// 3. execute group by
			e.pendingForGrouping.Inc()
			e.queryFlow.Grouping(func() {
				defer func() {
					e.pendingForGrouping.Dec()
					// try start collect tag values
					e.collectGroupByTagValues()
				}()
				e.executeGroupBy(shard, rs, rs.getSeriesIDs())
			})
		})
	}
}

// executeGroupBy executes the query flow, step as below:
// 1. grouping
// 2. loading
func (e *storageExecutor) executeGroupBy(shard tsdb.Shard, rs *timeSpanResultSet, seriesIDs *roaring.Bitmap) {
	groupingResult := &groupingResult{}
	var groupingCtx series.GroupingContext
	// time spans sorted by family
	ts := rs.getTimeSpans()
	if e.ctx.query.HasGroupBy() {
		// 1. grouping, if it has grouping, do group by tag keys, else just split series ids as batch first,
		// get grouping context if it needs
		tagKeys := make([]tag.KeyID, len(e.groupByTagKeyIDs))
		for idx, tagKeyID := range e.groupByTagKeyIDs {
			tagKeys[idx] = tagKeyID.ID
		}
		t := newGroupingContextFindTask(e.ctx, shard, tagKeys, seriesIDs, groupingResult)
		err := t.Run()
		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			// maybe group by not found, so ignore not found err
			e.queryFlow.Complete(err)
			return
		}
		if groupingResult.groupingCtx == nil {
			return
		}
		groupingCtx = groupingResult.groupingCtx
	}
	seriesIDsHighKeys := seriesIDs.GetHighKeys()
	e.pendingForGrouping.Add(int32(len(seriesIDsHighKeys)))
	var groupWait atomic.Int32
	groupWait.Add(int32(len(seriesIDsHighKeys)))

	for seriesIDHighKeyIdx, seriesIDHighKey := range seriesIDsHighKeys {
		seriesIDHighKey := seriesIDHighKey
		// be carefully, need use new variable for variable scope problem
		containerOfSeries := seriesIDs.GetContainerAtIndex(seriesIDHighKeyIdx)

		// grouping based on group by tag keys for each container
		e.queryFlow.Grouping(func() {
			defer func() {
				groupWait.Dec()
				if groupingCtx != nil && groupWait.Load() == 0 {
					// current group by query completed, need merge group by tag value ids
					e.mergeGroupByTagValueIDs(groupingCtx.GetGroupByTagValueIDs())
				}
				e.pendingForGrouping.Dec()
				// try start collect tag values for group by query
				e.collectGroupByTagValues()
			}()
			groupedResult := &groupedSeriesResult{}
			t := newBuildGroupTaskFunc(e.ctx, shard, groupingCtx, seriesIDHighKey, containerOfSeries, groupedResult)
			if err := t.Run(); err != nil {
				e.queryFlow.Complete(err)
				return
			}

			e.queryFlow.Load(func() {
				timeSpanCtxs := make([]*timeSpanCtx, len(ts))
				for idx, span := range ts {
					tsc := &timeSpanCtx{}
					timeSpanCtxs[idx] = tsc
					// 3.load data by grouped lowSeriesIDs
					t := newDataLoadTaskFunc(e.ctx, shard, e.queryFlow, span, tsc,
						seriesIDHighKey, containerOfSeries)
					if err := t.Run(); err != nil {
						e.queryFlow.Complete(err)
						return
					}
				}
				grouped := groupedResult.groupedSeries
				fieldAggList := make(aggregation.FieldAggregates, len(e.fields))
				aggSpecs := e.storageExecutePlan.getAggregatorSpecs()

				for idx := range e.fields {
					fieldAggList[idx] = aggregation.NewSeriesAggregator(
						e.ctx.query.Interval,
						e.queryIntervalRatio,
						e.ctx.query.TimeRange,
						aggSpecs[idx])
				}

				defer func() {
					if r := recover(); r != nil {
						storageQueryFlowLogger.Error("executeGroupBy",
							logger.Any("error", r),
							logger.Stack())
					}
				}()
				// tag values => low series ids
				for tags, lowSeriesIDs := range grouped {
					// scan metric data from storage(memory/file)
					for _, lowSeriesID := range lowSeriesIDs {
						for timeIdx, span := range ts { // family => result set
							// loads the metric data by given series id from load result.
							tsc := timeSpanCtxs[timeIdx]
							for resultSetIdx, loader := range tsc.loaders {
								if loader == nil {
									continue
								}
								// load field series data by series ids
								slotRange2, fieldSpanBinary := loader.Load(lowSeriesID)
								for fieldIndex := range fieldSpanBinary {
									spanBinary := fieldSpanBinary[fieldIndex]
									fieldsTSDDecoders := tsc.decoders[fieldIndex]
									if spanBinary != nil {
										if fieldsTSDDecoders[resultSetIdx] == nil {
											d := encoding.GetTSDDecoder()
											fieldsTSDDecoders[resultSetIdx] = d
										}
										fieldsTSDDecoders[resultSetIdx].ResetWithTimeRange(spanBinary, slotRange2.Start, slotRange2.End)
									}
								}
							}

							for idx, fieldSeries := range tsc.decoders {
								var agg aggregation.FieldAggregator
								var ok bool
								agg, ok = fieldAggList[idx].GetAggregator(span.familyTime)
								if !ok {
									continue
								}
								start, end := agg.SlotRange()
								target := timeutil.SlotRange{
									Start: uint16(start),
									End:   uint16(end),
								}
								aggregation.DownSamplingMultiSeriesInto(
									target, uint16(e.queryIntervalRatio), 0, // same family, base slot = 0
									e.fields[idx].Type,
									fieldSeries,
									agg.AggregateBySlot,
								)
							}
						}
					}
					e.queryFlow.Reduce(tags, fieldAggList.ResultSet(tags))
					// reset aggregate context
					fieldAggList.Reset()
				}
			})
		})
	}
}

// mergeGroupByTagValueIDs merges group by tag value ids for each shard
func (e *storageExecutor) mergeGroupByTagValueIDs(tagValueIDs []*roaring.Bitmap) {
	if tagValueIDs == nil {
		return
	}
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for idx, tagVIDs := range e.tagValueIDs {
		if tagVIDs == nil {
			e.tagValueIDs[idx] = tagValueIDs[idx]
		} else {
			tagVIDs.Or(tagValueIDs[idx])
		}
	}
}

// collectGroupByTagValues collects group tag values
func (e *storageExecutor) collectGroupByTagValues() {
	// all shard pending query tasks and grouping task completed, start collect tag values
	if e.pendingForShard.Load() == 0 && e.pendingForGrouping.Load() == 0 {
		if e.collecting.CAS(false, true) {
			for idx, tagKeyID := range e.groupByTagKeyIDs {
				tagKey := tagKeyID
				tagValueIDs := e.tagValueIDs[idx]
				tagIndex := idx
				if tagValueIDs == nil || tagValueIDs.IsEmpty() {
					e.queryFlow.ReduceTagValues(tagIndex, nil)
					continue
				}
				e.queryFlow.Load(func() {
					tagValues := make(map[uint32]string)
					t := newCollectTagValuesTask(e.ctx, e.database.Metadata(), tagKey, tagValueIDs, tagValues)
					if err := t.Run(); err != nil {
						e.queryFlow.Complete(err)
						return
					}
					e.queryFlow.ReduceTagValues(tagIndex, tagValues)
				})
			}
		}
	}
}

// validation validates query input params are valid
func (e *storageExecutor) validation() error {
	// check input shardIDs if empty
	if len(e.ctx.shardIDs) == 0 {
		return errNoShardID
	}
	numOfShards := e.database.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return errNoShardInDatabase
	}

	return nil
}

// checkShards checks got shards if valid
func (e *storageExecutor) checkShards() error {
	numOfShards := len(e.shards)
	if numOfShards == 0 {
		return errShardNotFound
	}
	numOfShardIDs := len(e.ctx.shardIDs)
	if numOfShards != numOfShardIDs {
		return errShardNumNotMatch
	}
	return nil
}
