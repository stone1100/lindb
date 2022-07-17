package stage

type Type int

const (
	Unknown Type = iota
	MetadataLookup
	ShardScan
	Grouping
	DataLoad
)

func (t Type) String() string {
	switch t {
	case MetadataLookup:
		return "MetadataLookup"
	case ShardScan:
		return "ShardScan"
	case Grouping:
		return "Grouping"
	case DataLoad:
		return "DataLoad"
	default:
		return "Unknown"
	}
}

type Stage interface {
	Plan() PlanNode
	NextStages() []Stage
	Submit(task func())
	Type() Type
}
