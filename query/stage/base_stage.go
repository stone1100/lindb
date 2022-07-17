package stage

type baseStage struct {
	stageType Type
}

func (stage *baseStage) Submit(task func()) {
	task()
}

func (stage *baseStage) Type() Type {
	return stage.stageType
}
