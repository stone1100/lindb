package storagequery

import (
	"fmt"
	"github.com/lindb/lindb/query/stage"
)

type pipeline struct {
}

func NewExecutePipeline() *pipeline {
	return &pipeline{}
}

func (p *pipeline) Execute(stage stage.Stage) {
	p.executeStage(stage)
}

func (p *pipeline) executeStage(stage stage.Stage) {
	if stage == nil {
		return
	}
	fmt.Println(stage.Type())
	p.execute(stage, stage.Plan())

	nextStages := stage.NextStages()
	for idx := range nextStages {
		p.executeStage(nextStages[idx])
	}
}

func (p *pipeline) execute(execStage stage.Stage, node stage.PlanNode) {
	if node == nil {
		return
	}

	// define plan node execute function
	exec := func(n stage.PlanNode) {
		// execute current plan node logic
		if err := n.Execute(); err != nil {
			// TODO check error
			fmt.Println(err)
			return
		}
		// if it has child node, need execute child node logic
		children := n.Children()
		for idx := range children {
			p.execute(execStage, children[idx])
		}
	}

	if node.Async() {
		execStage.Submit(func() {
			exec(node)
		})
	} else {
		exec(node)
	}
}
