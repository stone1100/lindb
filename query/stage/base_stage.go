package stage

type execTask func(task func())

type baseStage struct {
	stageType Type

	exec execTask
}

func (stage *baseStage) Submit(task func()) {
	stage.exec(task)
}

func (stage *baseStage) Type() Type {
	return stage.stageType
}
