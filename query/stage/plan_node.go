package stage

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

type PlanNode interface {
	Async() bool
	Execute() error
	Children() []PlanNode
	AddChild(node PlanNode)
}

type planNode struct {
	op    operator.Operator
	async bool

	children []PlanNode
}

func NewRootPlanNode() PlanNode {
	return &planNode{}
}

func NewPlanNode(op operator.Operator) PlanNode {
	return &planNode{
		op: op,
	}
}

func NewAsyncPlanNode(op operator.Operator) PlanNode {
	return &planNode{
		op:    op,
		async: true,
	}
}

func (p *planNode) Async() bool {
	return p.async
}

func (p *planNode) Execute() error {
	if p.op != nil {
		return p.op.Execute()
	}
	return nil
}

func (p *planNode) Children() []PlanNode {
	return p.children
}

func (p *planNode) AddChild(child PlanNode) {
	p.children = append(p.children, child)
}

type storagePhysicalPlan struct {
	database   tsdb.Database
	executeCtx *flow.StorageExecuteContext
	planTree   PlanNode // root
}
