package storagequery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lindb/lindb/query/stage"
	"sync"
)

type pipelineStateMachine struct {
	stages map[string]stage.Stage

	completedCallbackFn func()
	mutex               sync.Mutex
}

func (sm *pipelineStateMachine) startStage(stageID string, stage stage.Stage) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.stages[stageID] = stage
}

func (sm *pipelineStateMachine) completeStage(stageID string, _ stage.Stage) {
	isCompleted := false

	sm.mutex.Lock()
	delete(sm.stages, stageID)
	// all stage executed completed
	isCompleted = len(sm.stages) == 0
	sm.mutex.Unlock()

	if isCompleted && sm.completedCallbackFn != nil {
		// check if all stages execute completed
		sm.completedCallbackFn()
	}
}

func (sm *pipelineStateMachine) comp(completedCallbackFn func()) {
	sm.completedCallbackFn = completedCallbackFn
}

type pipeline struct {
	sm *pipelineStateMachine
}

func NewExecutePipeline() *pipeline {
	return &pipeline{
		sm: &pipelineStateMachine{
			stages: make(map[string]stage.Stage),
		},
	}
}

func (p *pipeline) Execute(stage stage.Stage) {
	p.executeStage(stage)
}

func (p *pipeline) Complete(fn func()) {
	fn()
}

func (p *pipeline) executeStage(stage stage.Stage) {
	if stage == nil {
		return
	}

	stageID := uuid.New().String()

	p.sm.startStage(stageID, stage)
	defer func() {
		p.sm.completeStage(stageID, stage)
		stage.Complete()
	}()

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
