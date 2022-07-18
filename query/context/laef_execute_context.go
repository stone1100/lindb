package context

import (
	"fmt"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

type LeafExecuteContext struct {
	TaskCtx  *flow.TaskContext
	LeafNode *models.Leaf

	StorageExecuteContext *flow.StorageExecuteContext
	Database              tsdb.Database

	ServerFactory rpc.TaskServerFactory
	Req           *protoCommonV1.TaskRequest
}

func (ctx *LeafExecuteContext) SendResponse(resultData [][]byte, err error) {
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
