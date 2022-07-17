package storagequery

//
//
//type storagePhysicalPlan struct {
//	database   tsdb.Database
//	executeCtx *flow.StorageExecuteContext
//	planTree   stage.PlanNode // root
//}
//
//func NewStoragePhysicalPlan(database tsdb.Database, shardIDs []models.ShardID, query *stmt.Query) *storagePhysicalPlan {
//	return &storagePhysicalPlan{
//		executeCtx: &flow.StorageExecuteContext{
//			Query:    query,
//			ShardIDs: shardIDs,
//			TagKeys:  make(map[string]tag.KeyID),
//		},
//	}
//}
//
//func (p *storagePhysicalPlan) Plan() error {
//	// do query validation
//	if err := p.validation(); err != nil {
//		return err
//	}
//	executeCtx := p.executeCtx
//
//	// add metadata lookup(name/tag/field etc.) node
//	p.planTree.AddChild(NewPlanNode(operator.NewMetadataLookup(executeCtx, p.database)))
//	hasWhereCondition := executeCtx.Query.Condition != nil
//	if hasWhereCondition {
//		// add tag values lookup node
//		p.planTree.AddChild(NewPlanNode(operator.NewTagValuesLookup(executeCtx, p.database)))
//	}
//
//	numOfShardIDs := len(executeCtx.ShardIDs)
//	foundShard := 0
//	// get shard by given query shard id list
//	for shardIdx, shardID := range executeCtx.ShardIDs {
//		// build shard execute context
//		shardExecuteCtx := flow.NewShardExecuteContext(executeCtx)
//		executeCtx.ShardContexts[shardIdx] = shardExecuteCtx
//
//		// if shard exist, add shard to query list
//		if shard, ok := p.database.GetShard(shardID); ok {
//			// add shard reader node
//			shardPlanNode := NewAsyncPlanNode(operator.NewShardReader())
//
//			if hasWhereCondition {
//				// add shard level series filtering node
//				shardPlanNode.AddChild(NewPlanNode(operator.NewSeriesFiltering(shardExecuteCtx, shard)))
//			} else {
//				// add shard level all series lookup node
//				shardPlanNode.AddChild(NewPlanNode(operator.NewMetricAllSeries(shardExecuteCtx, shard)))
//			}
//
//			families := shard.GetDataFamilies(executeCtx.Query.StorageInterval.Type(), executeCtx.Query.TimeRange)
//			for idx := range families {
//				family := families[idx]
//				// add data family reader node
//				shardPlanNode.AddChild(NewPlanNode(operator.NewDataFamilyReader(shardExecuteCtx, family)))
//			}
//			p.planTree.AddChild(shardPlanNode)
//			foundShard++
//		}
//	}
//	// check got shards if valid
//	if foundShard == 0 {
//		return errShardNotFound
//	}
//	if foundShard != numOfShardIDs {
//		return errShardNumNotMatch
//	}
//	return nil
//}
//
//// validation validates query input params are valid.
//func (p *storagePhysicalPlan) validation() error {
//	// check input shardIDs if empty
//	if len(p.executeCtx.ShardIDs) == 0 {
//		return errNoShardID
//	}
//	numOfShards := p.database.NumOfShards()
//	// check engine has shard
//	if numOfShards == 0 {
//		return errNoShardInDatabase
//	}
//	return nil
//}
//
//type executor struct {
//	queryFlow flow.StorageQueryFlow
//}
//
//func (e *executor) Execute(tree PlanNode) {
//	e.execute(tree)
//}
//
//func (e *executor) execute(node PlanNode) {
//	if node == nil {
//		return
//	}
//
//	exec := func(n PlanNode) {
//		n.Execute()
//		children := n.Children()
//		for idx := range children {
//			child := children[idx]
//			e.execute(child)
//		}
//	}
//
//	if node.Async() {
//		e.queryFlow.Submit(flow.GroupingStage, func() {
//			exec(node)
//		})
//	} else {
//		exec(node)
//	}
//}
