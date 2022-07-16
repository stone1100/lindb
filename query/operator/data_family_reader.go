package operator

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
)

type dataFamilyReader struct {
	executeCtx *flow.ShardExecuteContext
	family     tsdb.DataFamily
}

func NewDataFamilyReader(executeCtx *flow.ShardExecuteContext, family tsdb.DataFamily) Operator {
	return &dataFamilyReader{
		executeCtx: executeCtx,
		family:     family,
	}
}

func (op *dataFamilyReader) Execute() error {
	family := op.family
	resultSet, err := family.Filter(op.executeCtx)
	if err != nil {
		return err
	}
	for _, rs := range resultSet {
		op.executeCtx.TimeSegmentContext.AddFilterResultSet(family.Interval(), rs)
	}
	return nil
}
