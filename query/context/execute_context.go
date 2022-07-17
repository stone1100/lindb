package context

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb"
)

type LeafExecuteContext struct {
	TaskCtx  *flow.TaskContext
	LeafNode *models.Leaf

	StorageExecuteContext *flow.StorageExecuteContext
	Database              tsdb.Database
}
