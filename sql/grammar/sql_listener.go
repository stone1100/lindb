// Code generated from ./sql/grammar/SQL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr4-go/antlr/v4"

// SQLListener is a complete listener for a parse tree produced by SQLParser.
type SQLListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterUseStmt is called when entering the useStmt production.
	EnterUseStmt(c *UseStmtContext)

	// EnterSetLimitStmt is called when entering the setLimitStmt production.
	EnterSetLimitStmt(c *SetLimitStmtContext)

	// EnterShowStmt is called when entering the showStmt production.
	EnterShowStmt(c *ShowStmtContext)

	// EnterShowMasterStmt is called when entering the showMasterStmt production.
	EnterShowMasterStmt(c *ShowMasterStmtContext)

	// EnterShowBrokersStmt is called when entering the showBrokersStmt production.
	EnterShowBrokersStmt(c *ShowBrokersStmtContext)

	// EnterShowRequestsStmt is called when entering the showRequestsStmt production.
	EnterShowRequestsStmt(c *ShowRequestsStmtContext)

	// EnterRecoverStorageStmt is called when entering the recoverStorageStmt production.
	EnterRecoverStorageStmt(c *RecoverStorageStmtContext)

	// EnterShowLimitStmt is called when entering the showLimitStmt production.
	EnterShowLimitStmt(c *ShowLimitStmtContext)

	// EnterShowMetadataTypesStmt is called when entering the showMetadataTypesStmt production.
	EnterShowMetadataTypesStmt(c *ShowMetadataTypesStmtContext)

	// EnterShowMetadatasStmt is called when entering the showMetadatasStmt production.
	EnterShowMetadatasStmt(c *ShowMetadatasStmtContext)

	// EnterShowAliveStmt is called when entering the showAliveStmt production.
	EnterShowAliveStmt(c *ShowAliveStmtContext)

	// EnterShowReplicationsStmt is called when entering the showReplicationsStmt production.
	EnterShowReplicationsStmt(c *ShowReplicationsStmtContext)

	// EnterShowStateStmt is called when entering the showStateStmt production.
	EnterShowStateStmt(c *ShowStateStmtContext)

	// EnterCreateDatabaseStmt is called when entering the createDatabaseStmt production.
	EnterCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// EnterDropDatabaseStmt is called when entering the dropDatabaseStmt production.
	EnterDropDatabaseStmt(c *DropDatabaseStmtContext)

	// EnterShowDatabasesStmt is called when entering the showDatabasesStmt production.
	EnterShowDatabasesStmt(c *ShowDatabasesStmtContext)

	// EnterCreateBrokerStmt is called when entering the createBrokerStmt production.
	EnterCreateBrokerStmt(c *CreateBrokerStmtContext)

	// EnterShowNamespacesStmt is called when entering the showNamespacesStmt production.
	EnterShowNamespacesStmt(c *ShowNamespacesStmtContext)

	// EnterShowMetricsStmt is called when entering the showMetricsStmt production.
	EnterShowMetricsStmt(c *ShowMetricsStmtContext)

	// EnterShowFieldsStmt is called when entering the showFieldsStmt production.
	EnterShowFieldsStmt(c *ShowFieldsStmtContext)

	// EnterShowTagKeysStmt is called when entering the showTagKeysStmt production.
	EnterShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// EnterShowTagValuesStmt is called when entering the showTagValuesStmt production.
	EnterShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// EnterQueryStmt is called when entering the queryStmt production.
	EnterQueryStmt(c *QueryStmtContext)

	// EnterSourceAndSelect is called when entering the sourceAndSelect production.
	EnterSourceAndSelect(c *SourceAndSelectContext)

	// EnterSelectExpr is called when entering the selectExpr production.
	EnterSelectExpr(c *SelectExprContext)

	// EnterFields is called when entering the fields production.
	EnterFields(c *FieldsContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterAlias is called when entering the alias production.
	EnterAlias(c *AliasContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterConditionExpr is called when entering the conditionExpr production.
	EnterConditionExpr(c *ConditionExprContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterValueList is called when entering the valueList production.
	EnterValueList(c *ValueListContext)

	// EnterTimeExpr is called when entering the timeExpr production.
	EnterTimeExpr(c *TimeExprContext)

	// EnterNowExpr is called when entering the nowExpr production.
	EnterNowExpr(c *NowExprContext)

	// EnterNowFunc is called when entering the nowFunc production.
	EnterNowFunc(c *NowFuncContext)

	// EnterGroupByClause is called when entering the groupByClause production.
	EnterGroupByClause(c *GroupByClauseContext)

	// EnterGroupByKeys is called when entering the groupByKeys production.
	EnterGroupByKeys(c *GroupByKeysContext)

	// EnterGroupByKey is called when entering the groupByKey production.
	EnterGroupByKey(c *GroupByKeyContext)

	// EnterFillOption is called when entering the fillOption production.
	EnterFillOption(c *FillOptionContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterSortField is called when entering the sortField production.
	EnterSortField(c *SortFieldContext)

	// EnterSortFields is called when entering the sortFields production.
	EnterSortFields(c *SortFieldsContext)

	// EnterHavingClause is called when entering the havingClause production.
	EnterHavingClause(c *HavingClauseContext)

	// EnterBoolExpr is called when entering the boolExpr production.
	EnterBoolExpr(c *BoolExprContext)

	// EnterBoolExprLogicalOp is called when entering the boolExprLogicalOp production.
	EnterBoolExprLogicalOp(c *BoolExprLogicalOpContext)

	// EnterBoolExprAtom is called when entering the boolExprAtom production.
	EnterBoolExprAtom(c *BoolExprAtomContext)

	// EnterBinaryExpr is called when entering the binaryExpr production.
	EnterBinaryExpr(c *BinaryExprContext)

	// EnterBinaryOperator is called when entering the binaryOperator production.
	EnterBinaryOperator(c *BinaryOperatorContext)

	// EnterFieldExpr is called when entering the fieldExpr production.
	EnterFieldExpr(c *FieldExprContext)

	// EnterStar is called when entering the star production.
	EnterStar(c *StarContext)

	// EnterDurationLit is called when entering the durationLit production.
	EnterDurationLit(c *DurationLitContext)

	// EnterIntervalItem is called when entering the intervalItem production.
	EnterIntervalItem(c *IntervalItemContext)

	// EnterExprFunc is called when entering the exprFunc production.
	EnterExprFunc(c *ExprFuncContext)

	// EnterFuncName is called when entering the funcName production.
	EnterFuncName(c *FuncNameContext)

	// EnterExprFuncParams is called when entering the exprFuncParams production.
	EnterExprFuncParams(c *ExprFuncParamsContext)

	// EnterFuncParam is called when entering the funcParam production.
	EnterFuncParam(c *FuncParamContext)

	// EnterExprAtom is called when entering the exprAtom production.
	EnterExprAtom(c *ExprAtomContext)

	// EnterProperties is called when entering the properties production.
	EnterProperties(c *PropertiesContext)

	// EnterPropertyAssignments is called when entering the propertyAssignments production.
	EnterPropertyAssignments(c *PropertyAssignmentsContext)

	// EnterProperty is called when entering the property production.
	EnterProperty(c *PropertyContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterIntNumber is called when entering the intNumber production.
	EnterIntNumber(c *IntNumberContext)

	// EnterDecNumber is called when entering the decNumber production.
	EnterDecNumber(c *DecNumberContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterMetricName is called when entering the metricName production.
	EnterMetricName(c *MetricNameContext)

	// EnterTagKey is called when entering the tagKey production.
	EnterTagKey(c *TagKeyContext)

	// EnterTagValue is called when entering the tagValue production.
	EnterTagValue(c *TagValueContext)

	// EnterPrefix is called when entering the prefix production.
	EnterPrefix(c *PrefixContext)

	// EnterWithTagKey is called when entering the withTagKey production.
	EnterWithTagKey(c *WithTagKeyContext)

	// EnterNamespace is called when entering the namespace production.
	EnterNamespace(c *NamespaceContext)

	// EnterName is called when entering the name production.
	EnterName(c *NameContext)

	// EnterRequestID is called when entering the requestID production.
	EnterRequestID(c *RequestIDContext)

	// EnterToml is called when entering the toml production.
	EnterToml(c *TomlContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// EnterNonReservedWords is called when entering the nonReservedWords production.
	EnterNonReservedWords(c *NonReservedWordsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitUseStmt is called when exiting the useStmt production.
	ExitUseStmt(c *UseStmtContext)

	// ExitSetLimitStmt is called when exiting the setLimitStmt production.
	ExitSetLimitStmt(c *SetLimitStmtContext)

	// ExitShowStmt is called when exiting the showStmt production.
	ExitShowStmt(c *ShowStmtContext)

	// ExitShowMasterStmt is called when exiting the showMasterStmt production.
	ExitShowMasterStmt(c *ShowMasterStmtContext)

	// ExitShowBrokersStmt is called when exiting the showBrokersStmt production.
	ExitShowBrokersStmt(c *ShowBrokersStmtContext)

	// ExitShowRequestsStmt is called when exiting the showRequestsStmt production.
	ExitShowRequestsStmt(c *ShowRequestsStmtContext)

	// ExitRecoverStorageStmt is called when exiting the recoverStorageStmt production.
	ExitRecoverStorageStmt(c *RecoverStorageStmtContext)

	// ExitShowLimitStmt is called when exiting the showLimitStmt production.
	ExitShowLimitStmt(c *ShowLimitStmtContext)

	// ExitShowMetadataTypesStmt is called when exiting the showMetadataTypesStmt production.
	ExitShowMetadataTypesStmt(c *ShowMetadataTypesStmtContext)

	// ExitShowMetadatasStmt is called when exiting the showMetadatasStmt production.
	ExitShowMetadatasStmt(c *ShowMetadatasStmtContext)

	// ExitShowAliveStmt is called when exiting the showAliveStmt production.
	ExitShowAliveStmt(c *ShowAliveStmtContext)

	// ExitShowReplicationsStmt is called when exiting the showReplicationsStmt production.
	ExitShowReplicationsStmt(c *ShowReplicationsStmtContext)

	// ExitShowStateStmt is called when exiting the showStateStmt production.
	ExitShowStateStmt(c *ShowStateStmtContext)

	// ExitCreateDatabaseStmt is called when exiting the createDatabaseStmt production.
	ExitCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// ExitDropDatabaseStmt is called when exiting the dropDatabaseStmt production.
	ExitDropDatabaseStmt(c *DropDatabaseStmtContext)

	// ExitShowDatabasesStmt is called when exiting the showDatabasesStmt production.
	ExitShowDatabasesStmt(c *ShowDatabasesStmtContext)

	// ExitCreateBrokerStmt is called when exiting the createBrokerStmt production.
	ExitCreateBrokerStmt(c *CreateBrokerStmtContext)

	// ExitShowNamespacesStmt is called when exiting the showNamespacesStmt production.
	ExitShowNamespacesStmt(c *ShowNamespacesStmtContext)

	// ExitShowMetricsStmt is called when exiting the showMetricsStmt production.
	ExitShowMetricsStmt(c *ShowMetricsStmtContext)

	// ExitShowFieldsStmt is called when exiting the showFieldsStmt production.
	ExitShowFieldsStmt(c *ShowFieldsStmtContext)

	// ExitShowTagKeysStmt is called when exiting the showTagKeysStmt production.
	ExitShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// ExitShowTagValuesStmt is called when exiting the showTagValuesStmt production.
	ExitShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// ExitQueryStmt is called when exiting the queryStmt production.
	ExitQueryStmt(c *QueryStmtContext)

	// ExitSourceAndSelect is called when exiting the sourceAndSelect production.
	ExitSourceAndSelect(c *SourceAndSelectContext)

	// ExitSelectExpr is called when exiting the selectExpr production.
	ExitSelectExpr(c *SelectExprContext)

	// ExitFields is called when exiting the fields production.
	ExitFields(c *FieldsContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitAlias is called when exiting the alias production.
	ExitAlias(c *AliasContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitConditionExpr is called when exiting the conditionExpr production.
	ExitConditionExpr(c *ConditionExprContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitValueList is called when exiting the valueList production.
	ExitValueList(c *ValueListContext)

	// ExitTimeExpr is called when exiting the timeExpr production.
	ExitTimeExpr(c *TimeExprContext)

	// ExitNowExpr is called when exiting the nowExpr production.
	ExitNowExpr(c *NowExprContext)

	// ExitNowFunc is called when exiting the nowFunc production.
	ExitNowFunc(c *NowFuncContext)

	// ExitGroupByClause is called when exiting the groupByClause production.
	ExitGroupByClause(c *GroupByClauseContext)

	// ExitGroupByKeys is called when exiting the groupByKeys production.
	ExitGroupByKeys(c *GroupByKeysContext)

	// ExitGroupByKey is called when exiting the groupByKey production.
	ExitGroupByKey(c *GroupByKeyContext)

	// ExitFillOption is called when exiting the fillOption production.
	ExitFillOption(c *FillOptionContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitSortField is called when exiting the sortField production.
	ExitSortField(c *SortFieldContext)

	// ExitSortFields is called when exiting the sortFields production.
	ExitSortFields(c *SortFieldsContext)

	// ExitHavingClause is called when exiting the havingClause production.
	ExitHavingClause(c *HavingClauseContext)

	// ExitBoolExpr is called when exiting the boolExpr production.
	ExitBoolExpr(c *BoolExprContext)

	// ExitBoolExprLogicalOp is called when exiting the boolExprLogicalOp production.
	ExitBoolExprLogicalOp(c *BoolExprLogicalOpContext)

	// ExitBoolExprAtom is called when exiting the boolExprAtom production.
	ExitBoolExprAtom(c *BoolExprAtomContext)

	// ExitBinaryExpr is called when exiting the binaryExpr production.
	ExitBinaryExpr(c *BinaryExprContext)

	// ExitBinaryOperator is called when exiting the binaryOperator production.
	ExitBinaryOperator(c *BinaryOperatorContext)

	// ExitFieldExpr is called when exiting the fieldExpr production.
	ExitFieldExpr(c *FieldExprContext)

	// ExitStar is called when exiting the star production.
	ExitStar(c *StarContext)

	// ExitDurationLit is called when exiting the durationLit production.
	ExitDurationLit(c *DurationLitContext)

	// ExitIntervalItem is called when exiting the intervalItem production.
	ExitIntervalItem(c *IntervalItemContext)

	// ExitExprFunc is called when exiting the exprFunc production.
	ExitExprFunc(c *ExprFuncContext)

	// ExitFuncName is called when exiting the funcName production.
	ExitFuncName(c *FuncNameContext)

	// ExitExprFuncParams is called when exiting the exprFuncParams production.
	ExitExprFuncParams(c *ExprFuncParamsContext)

	// ExitFuncParam is called when exiting the funcParam production.
	ExitFuncParam(c *FuncParamContext)

	// ExitExprAtom is called when exiting the exprAtom production.
	ExitExprAtom(c *ExprAtomContext)

	// ExitProperties is called when exiting the properties production.
	ExitProperties(c *PropertiesContext)

	// ExitPropertyAssignments is called when exiting the propertyAssignments production.
	ExitPropertyAssignments(c *PropertyAssignmentsContext)

	// ExitProperty is called when exiting the property production.
	ExitProperty(c *PropertyContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitIntNumber is called when exiting the intNumber production.
	ExitIntNumber(c *IntNumberContext)

	// ExitDecNumber is called when exiting the decNumber production.
	ExitDecNumber(c *DecNumberContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitMetricName is called when exiting the metricName production.
	ExitMetricName(c *MetricNameContext)

	// ExitTagKey is called when exiting the tagKey production.
	ExitTagKey(c *TagKeyContext)

	// ExitTagValue is called when exiting the tagValue production.
	ExitTagValue(c *TagValueContext)

	// ExitPrefix is called when exiting the prefix production.
	ExitPrefix(c *PrefixContext)

	// ExitWithTagKey is called when exiting the withTagKey production.
	ExitWithTagKey(c *WithTagKeyContext)

	// ExitNamespace is called when exiting the namespace production.
	ExitNamespace(c *NamespaceContext)

	// ExitName is called when exiting the name production.
	ExitName(c *NameContext)

	// ExitRequestID is called when exiting the requestID production.
	ExitRequestID(c *RequestIDContext)

	// ExitToml is called when exiting the toml production.
	ExitToml(c *TomlContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)

	// ExitNonReservedWords is called when exiting the nonReservedWords production.
	ExitNonReservedWords(c *NonReservedWordsContext)
}
