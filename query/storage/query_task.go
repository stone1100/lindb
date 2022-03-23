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
	"fmt"
	"strings"
	"time"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

// baseQueryTask represents base query task stats track for task execute cost
type baseQueryTask struct {
	start time.Time
	cost  time.Duration
}

// BeforeRun invokes before task run function
func (t *baseQueryTask) BeforeRun() {
	t.start = time.Now()
}

// Run executes task logic
func (t *baseQueryTask) Run() error {
	return nil
}

// AfterRun invokes after task run function
func (t *baseQueryTask) AfterRun() {
	t.cost = time.Since(t.start)
}

// queryStatTask represents the query stat task
type queryStatTask struct {
	task flow.QueryTask
}

// BeforeRun invokes before task run function
func (t *queryStatTask) BeforeRun() {
}

// Run executes query cost stat
func (t *queryStatTask) Run() error {
	t.task.BeforeRun()
	defer func() {
		t.task.AfterRun()
	}()
	return t.task.Run()
}

// AfterRun invokes after task run function
func (t *queryStatTask) AfterRun() {
}

// storagePlanTask represents storage execute plan task
type storagePlanTask struct {
	baseQueryTask
	ctx  *storageExecuteContext
	plan *storageExecutePlan
}

// newStoragePlanTask creates storage execute plan task
func newStoragePlanTask(ctx *storageExecuteContext, plan *storageExecutePlan) flow.QueryTask {
	task := &storagePlanTask{
		ctx:  ctx,
		plan: plan,
	}
	if ctx.query.Explain {
		// if need explain query, use queryStatTask
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes storage execute plan
func (t *storagePlanTask) Run() error {
	return t.plan.Plan()
}

// AfterRun invokes after execute plan, collects plan stats
func (t *storagePlanTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.stats.SetPlanCost(t.cost)
}

// tagFilterTask represents tag filtering task based on where condition
type tagFilterTask struct {
	baseQueryTask
	ctx       *storageExecuteContext
	tagSearch TagSearch
}

// newTagFilterTask creates tag filtering task
func newTagFilterTask(ctx *storageExecuteContext, tagSearch TagSearch) flow.QueryTask {
	task := &tagFilterTask{
		ctx:       ctx,
		tagSearch: tagSearch,
	}
	if ctx.query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes tag filtering based on where condition
func (t *tagFilterTask) Run() error {
	tagFilterResult, err := t.tagSearch.Filter()
	if err != nil {
		return err
	}
	if len(tagFilterResult) == 0 {
		// filter not match, return not found
		return constants.ErrNotFound
	}
	// set tag filter result
	t.ctx.setTagFilterResult(tagFilterResult)
	return nil
}

// AfterRun invokes after tag filtering, collects tag filtering stats
func (t *tagFilterTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.stats.SetTagFilterCost(t.cost)
}

// seriesIDsSearchTask represents series ids search task based on tag filtering result set
type seriesIDsSearchTask struct {
	baseQueryTask

	executeCtx *flow.StorageExecuteContext
	shard      tsdb.Shard
}

// newSeriesIDsSearchTask creates series ids search task
func newSeriesIDsSearchTask(executeCtx *flow.StorageExecuteContext, shard tsdb.Shard) flow.QueryTask {
	task := &seriesIDsSearchTask{
		executeCtx: executeCtx,
		shard:      shard,
	}
	if executeCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes series ids search based on tag filtering result
func (t *seriesIDsSearchTask) Run() (err error) {
	condition := t.executeCtx.Query.Condition
	var seriesIDs *roaring.Bitmap
	if condition != nil {
		// if it gets tag filter result do series ids searching
		seriesSearch := newSeriesSearchFunc(t.shard.IndexDatabase(), t.executeCtx.TagFilterResult, t.executeCtx.Query.Condition)
		seriesIDs, err = seriesSearch.Search()
	} else {
		// get series ids for metric level
		seriesIDs, err = t.shard.IndexDatabase().GetSeriesIDsForMetric(t.executeCtx.Query.Namespace,
			t.executeCtx.Query.MetricName)
		if err == nil && !t.executeCtx.Query.HasGroupBy() {
			// add series id without tags, maybe metric has too many series, but one series without tags
			seriesIDs.Add(series.IDWithoutTags)
		}
	}
	if err == nil && seriesIDs != nil {
		t.executeCtx.SeriesIDsAfterFiltering.Or(seriesIDs)
	}
	return
}

// AfterRun invokes after series ids search, collects the series ids search stats
func (t *seriesIDsSearchTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	//TODO
	//t.ctx.stats.SetShardSeriesIDsSearchStats(t.shard.ShardID(), t.result.GetCardinality(), t.cost)
}

// familyFilterTask represents family data filtering task
type familyFilterTask struct {
	baseQueryTask

	executeCtx *flow.StorageExecuteContext
	shard      tsdb.Shard

	rs *timeSpanResultSet
}

// newFamilyFilterTask creates family data filtering task
func newFamilyFilterTask(executeCtx *flow.StorageExecuteContext, shard tsdb.Shard, rs *timeSpanResultSet) flow.QueryTask {
	task := &familyFilterTask{
		executeCtx: executeCtx,
		shard:      shard,
		rs:         rs,
	}
	if executeCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes file data filtering based on series ids and time range for each data family
func (t *familyFilterTask) Run() error {
	families := t.shard.GetDataFamilies(t.executeCtx.Query.Interval.Type(), t.executeCtx.Query.TimeRange)
	if len(families) == 0 {
		return nil
	}
	for idx := range families {
		family := families[idx]
		// execute data family search in background goroutine
		resultSet, err := family.Filter(t.executeCtx)
		if err != nil {
			return err
		}
		for _, rs := range resultSet {
			t.rs.addFilterResultSet(family.Interval(), rs)
		}
	}
	return nil
}

// AfterRun invokes after file data filtering, collects the file data filtering stats
func (t *familyFilterTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	//TODO
	//t.ctx.stats.SetShardKVDataFilterCost(t.shard.ShardID(), t.cost)
}

// groupingContextFindTask represents group by context find task
type groupingContextFindTask struct {
	baseQueryTask

	executeCtx *flow.StorageExecuteContext
	shard      tsdb.Shard
	seriesIDs  *roaring.Bitmap
	result     *groupingResult
}

// newGroupingContextFindTask creates the group by context find task
func newGroupingContextFindTask(executeCtx *flow.StorageExecuteContext, shard tsdb.Shard, result *groupingResult,
) flow.QueryTask {
	task := &groupingContextFindTask{
		executeCtx: executeCtx,
		shard:            shard,
		result:           result,
	}
	if executeCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes group by context finding based on group by tag key ids
func (t *groupingContextFindTask) Run() error {
	gCtx, err := t.shard.IndexDatabase().GetGroupingContext(t.groupByTagKeyIDs, t.seriesIDs)
	if err != nil {
		return err
	}

	t.result.groupingCtx = gCtx
	return nil
}

// AfterRun invokes after group by context, collects the find group by context stats
func (t *groupingContextFindTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.stats.SetShardGroupingCost(t.shard.ShardID(), t.cost)
}

// buildGroupTask represents build grouped tag value ids => series ids mapping
type buildGroupTask struct {
	baseQueryTask
	ctx         *storageExecuteContext
	shard       tsdb.Shard
	groupingCtx series.GroupingContext
	highKey     uint16
	container   roaring.Container
	result      *groupedSeriesResult
}

// newBuildGroupTask creates build group task
func newBuildGroupTask(ctx *storageExecuteContext, shard tsdb.Shard,
	groupingCtx series.GroupingContext, highKey uint16, container roaring.Container,
	result *groupedSeriesResult,
) flow.QueryTask {
	task := &buildGroupTask{
		ctx:         ctx,
		shard:       shard,
		groupingCtx: groupingCtx,
		highKey:     highKey,
		container:   container,
		result:      result,
	}
	if ctx.query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes grouped series ids(tag value ids=>series ids mapping)
func (t *buildGroupTask) Run() error {
	if t.groupingCtx != nil {
		// build group by data, grouped series: tags => series IDs
		t.result.groupedSeries = t.groupingCtx.BuildGroup(t.highKey, t.container)
	} else {
		t.result.groupedSeries = map[string][]uint16{"": t.container.ToArray()}
	}
	return nil
}

// AfterRun invokes after build grouped series, collects build stats
func (t *buildGroupTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.stats.SetShardGroupBuildStats(t.shard.ShardID(), t.cost)
}

// dataLoadTask represents data load task based on filtering result set
type dataLoadTask struct {
	baseQueryTask
	ctx       *storageExecuteContext
	spanCtx   *timeSpanCtx
	shard     tsdb.Shard
	queryFlow flow.StorageQueryFlow
	timeSpan  *timeSpan
	highKey   uint16
	seriesIDs roaring.Container
}

// newDataLoadTask creates the data load task
func newDataLoadTask(
	ctx *storageExecuteContext,
	shard tsdb.Shard,
	queryFlow flow.StorageQueryFlow,
	timeSpan *timeSpan,
	spanCtx *timeSpanCtx,
	highKey uint16, seriesIDs roaring.Container,
) flow.QueryTask {
	task := &dataLoadTask{
		ctx:       ctx,
		shard:     shard,
		queryFlow: queryFlow,
		timeSpan:  timeSpan,
		spanCtx:   spanCtx,
		highKey:   highKey,
		seriesIDs: seriesIDs,
	}
	if ctx.query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes data load based on filtering result set
func (t *dataLoadTask) Run() error {
	t.spanCtx.loaders = make([]flow.DataLoader, len(t.timeSpan.resultSets))
	t.spanCtx.decoders = make([][]*encoding.TSDDecoder, t.timeSpan.fieldCount)
	for i := range t.spanCtx.decoders {
		t.spanCtx.decoders[i] = make([]*encoding.TSDDecoder, len(t.timeSpan.resultSets))
	}
	for idx, rs := range t.timeSpan.resultSets {
		// maybe return nil loader
		t.spanCtx.loaders[idx] = rs.Load(t.highKey, t.seriesIDs)
	}
	return nil
}

// AfterRun invokes after data load, collects the data load stats
func (t *dataLoadTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	// TODO need modify
	identifiers := strings.Split(t.timeSpan.identifier, fmt.Sprintf("shard/%d/segment", t.shard.ShardID()))
	var identifier string
	if len(identifiers) > 1 {
		identifier = identifiers[1]
	} else {
		identifier = identifiers[0]
	}
	t.ctx.stats.SetShardScanStats(t.shard.ShardID(), identifier, t.cost)
}

// collectTagValuesTask represents collect tag values by tag value ids
type collectTagValuesTask struct {
	baseQueryTask
	ctx         *storageExecuteContext
	metadata    metadb.Metadata
	tagKey      tag.Meta
	tagValueIDs *roaring.Bitmap
	tagValues   map[uint32]string
}

// newCollectTagValuesTask creates the collect tag values task
func newCollectTagValuesTask(ctx *storageExecuteContext, metadata metadb.Metadata,
	tagKey tag.Meta, tagValueIDs *roaring.Bitmap, tagValues map[uint32]string,
) flow.QueryTask {
	task := &collectTagValuesTask{
		ctx:         ctx,
		metadata:    metadata,
		tagKey:      tagKey,
		tagValueIDs: tagValueIDs,
		tagValues:   tagValues,
	}
	if ctx.query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes collect tag values by ids
func (t *collectTagValuesTask) Run() error {
	return t.metadata.TagMetadata().CollectTagValues(tag.KeyID(t.tagKey.ID), t.tagValueIDs, t.tagValues)
}

// AfterRun invokes after tag value collect, collects execution stats
func (t *collectTagValuesTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.stats.SetCollectTagValuesStats(t.tagKey.Key, t.cost)
}
