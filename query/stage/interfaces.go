package stage

import storagequery "github.com/lindb/lindb/query/storage"

type Stage interface {
	Plan() storagequery.PlanNode
	NextStages() []Stage
}
