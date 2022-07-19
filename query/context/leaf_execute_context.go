package context

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

type LeafExecuteContext struct {
	TaskCtx  *flow.TaskContext
	LeafNode *models.Leaf

	StorageExecuteCtx *flow.StorageExecuteContext
	Database          tsdb.Database

	ServerFactory rpc.TaskServerFactory
	Req           *protoCommonV1.TaskRequest

	GroupingCtx *LeafGroupingContext
}

func NewLeafExecuteContext(taskCtx *flow.TaskContext,
	queryStmt *stmt.Query,
	req *protoCommonV1.TaskRequest,
	serverFactory rpc.TaskServerFactory,
	leafNode *models.Leaf,
	database tsdb.Database,
) *LeafExecuteContext {
	ctx := &LeafExecuteContext{
		TaskCtx:  taskCtx,
		LeafNode: leafNode,
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query:    queryStmt,
			ShardIDs: leafNode.ShardIDs,
		},
		Database:      database,
		ServerFactory: serverFactory,
		Req:           req,
	}
	// for group by query
	if queryStmt.HasGroupBy() {
		ctx.GroupingCtx = newLeafGroupingContext(ctx.StorageExecuteCtx, database)
	}
	return ctx
}

func (ctx *LeafExecuteContext) waitCollectGroupingTagsCompleted() {
	if ctx.StorageExecuteCtx.HasGroupingTagValueIDs() {
		// if it has grouping tag value ids, need wait collect group by tag values completed
		select {
		case <-ctx.TaskCtx.Ctx.Done():
			ctx.sendResponse(nil, ctx.TaskCtx.Ctx.Err())
			return
		case <-ctx.collectGroupingTagsCompleted:
		}
	}
}

func (ctx *LeafExecuteContext) SendResponse() {
	defer ctx.StorageExecuteCtx.Release()

	ctx.waitCollectGroupingTagsCompleted()

	numOfReceivers := len(ctx.LeafNode.Receivers)
	resultSet := make([][]byte, numOfReceivers)
	if ctx.ReduceAgg != nil {

		timeSeriesList := qf.makeTimeSeriesList()
		// root -> leaf task, return the raw total series
		if len(ctx.LeafNode.Receivers) == 1 {
			leaf2RootSeries := protoCommonV1.TimeSeriesList{
				TimeSeriesList: timeSeriesList,
				FieldAggSpecs:  qf.aggregatorSpecs,
			}
			leaf2RootSeriesPayload, _ := leaf2RootSeries.Marshal()
			resultSet[0] = leaf2RootSeriesPayload
		} else {
			// during intermediate task, time series will be grouped by hash
			// and send to multi intermediate receiver
			// hash mod -> series list
			var timeSeriesHashGroups = make([][]*protoCommonV1.TimeSeries, numOfReceivers)
			for _, ts := range timeSeriesList {
				h := xxhash.Sum64String(ts.Tags)
				index := int(h % uint64(numOfReceivers))
				timeSeriesHashGroups[index] = append(timeSeriesHashGroups[index], ts)
			}
			for idx, timeSeriesHashGroup := range timeSeriesHashGroups {
				leaf2IntermediateSeries := protoCommonV1.TimeSeriesList{
					TimeSeriesList: timeSeriesHashGroup,
					FieldAggSpecs:  qf.aggregatorSpecs,
				}
				leaf2IntermediatePayload, _ := leaf2IntermediateSeries.Marshal()
				resultSet[idx] = leaf2IntermediatePayload
			}
		}
	}
	ctx.sendResponse(resultSet, nil)
}

func (ctx *LeafExecuteContext) sendResponse(resultData [][]byte, err error) {
	var stats []byte
	// TODO need get stats when err
	errMsg := err.Error()
	// send result to upstream receivers
	for idx, receiver := range ctx.LeafNode.Receivers {
		stream := ctx.ServerFactory.GetStream(receiver.Indicator())
		if stream == nil {
			// TODO
			//storageQueryFlowLogger.Error("unable to get stream for write response",
			//	logger.String("target", receiver.Indicator()))
			//ctx.Complete(querypkg.ErrNoSendStream)
			fmt.Println("stream not found")
			break
		}
		var payload []byte
		if resultData != nil {
			payload = resultData[idx]
		}
		resp := &protoCommonV1.TaskResponse{
			TaskID:    ctx.Req.ParentTaskID,
			Type:      protoCommonV1.TaskType_Leaf,
			Completed: true,
			SendTime:  timeutil.NowNano(),
			Payload:   payload,
			Stats:     stats,
			ErrMsg:    errMsg,
		}
		if err0 := stream.Send(resp); err0 != nil {
			//storageQueryFlowLogger.Error("send storage query result", logger.Error(err0))
			fmt.Println(err)
		}
	}
}
