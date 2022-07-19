package context

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
	"sync"

	"go.uber.org/atomic"
)

type LeafGroupingContext struct {
	storageExecuteCtx *flow.StorageExecuteContext
	database          tsdb.Database

	collectGroupingTagsCompleted chan struct{}       // collect completed signal
	groupingRelatedTasks         atomic.Int32        // track how many tasks are pending
	collectRelatedTasks          atomic.Int32        // track if collect grouping tag value tasks completed
	tagsMap                      map[string]string   // tag value ids => tag values
	tagValuesMap                 []map[uint32]string // tag value id=> tag value for each group by tag key
	tagValues                    []string

	mutex sync.Mutex
}

func newLeafGroupingContext(storageExecuteCtx *flow.StorageExecuteContext, database tsdb.Database) *LeafGroupingContext {
	groupByKenLen := len(storageExecuteCtx.Query.GroupBy)
	return &LeafGroupingContext{
		storageExecuteCtx:            storageExecuteCtx,
		database:                     database,
		tagValuesMap:                 make([]map[uint32]string, groupByKenLen),
		tagsMap:                      make(map[string]string),
		tagValues:                    make([]string, groupByKenLen),
		collectGroupingTagsCompleted: make(chan struct{}),
		collectRelatedTasks:          *atomic.NewInt32(int32(groupByKenLen)),
	}
}

func (ctx *LeafGroupingContext) ForkGroupingTask() {
	ctx.groupingRelatedTasks.Inc()
}

func (ctx *LeafGroupingContext) CompleteGroupingTask() {
	ctx.groupingRelatedTasks.Dec()

	ctx.collectGroupByTagValues()
}

// collectGroupByTagValues collects group tag values
func (ctx *LeafGroupingContext) collectGroupByTagValues() {
	if ctx.groupingRelatedTasks.Load() != 0 {
		return
	}
	// all shard pending query tasks and grouping task completed, start collect tag values
	storageExecuteCtx := ctx.storageExecuteCtx
	metadata := ctx.database.Metadata()
	for idx, tagKeyID := range storageExecuteCtx.GroupByTags {
		tagKey := tagKeyID
		tagValueIDs := storageExecuteCtx.GroupingTagValueIDs[idx]
		tagIndex := idx
		if tagValueIDs == nil || tagValueIDs.IsEmpty() {
			ctx.reduceTagValues(tagIndex, nil)
			continue
		}

		t.queryFlow.Submit(flow.ScannerStage, func() {
			tagValues := make(map[uint32]string) // tag value id => tag value
			task := newCollectTagValuesTaskFunc(t.ctx, metadata, tagKey, tagValueIDs, tagValues)
			if err := task.Run(); err != nil {
				t.queryFlow.Complete(err)
				return
			}
			ctx.reduceTagValues(tagIndex, tagValues)
		})
	}
}

// reduceTagValues reduces the group by tag values
func (ctx *LeafGroupingContext) reduceTagValues(tagKeyIndex int, tagValues map[uint32]string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.tagValuesMap[tagKeyIndex] = tagValues
	if ctx.collectRelatedTasks.Dec() == 0 {
		// notify all collect tag value tasks completed
		close(ctx.collectGroupingTagsCompleted)
	}
}
