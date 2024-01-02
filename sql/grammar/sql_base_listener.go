// Code generated from ./sql/grammar/SQL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr4-go/antlr/v4"

// BaseSQLListener is a complete listener for a parse tree produced by SQLParser.
type BaseSQLListener struct{}

var _ SQLListener = &BaseSQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseSQLListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseSQLListener) ExitStatement(ctx *StatementContext) {}

// EnterUseStmt is called when production useStmt is entered.
func (s *BaseSQLListener) EnterUseStmt(ctx *UseStmtContext) {}

// ExitUseStmt is called when production useStmt is exited.
func (s *BaseSQLListener) ExitUseStmt(ctx *UseStmtContext) {}

// EnterSetLimitStmt is called when production setLimitStmt is entered.
func (s *BaseSQLListener) EnterSetLimitStmt(ctx *SetLimitStmtContext) {}

// ExitSetLimitStmt is called when production setLimitStmt is exited.
func (s *BaseSQLListener) ExitSetLimitStmt(ctx *SetLimitStmtContext) {}

// EnterShowStmt is called when production showStmt is entered.
func (s *BaseSQLListener) EnterShowStmt(ctx *ShowStmtContext) {}

// ExitShowStmt is called when production showStmt is exited.
func (s *BaseSQLListener) ExitShowStmt(ctx *ShowStmtContext) {}

// EnterShowMasterStmt is called when production showMasterStmt is entered.
func (s *BaseSQLListener) EnterShowMasterStmt(ctx *ShowMasterStmtContext) {}

// ExitShowMasterStmt is called when production showMasterStmt is exited.
func (s *BaseSQLListener) ExitShowMasterStmt(ctx *ShowMasterStmtContext) {}

// EnterShowBrokersStmt is called when production showBrokersStmt is entered.
func (s *BaseSQLListener) EnterShowBrokersStmt(ctx *ShowBrokersStmtContext) {}

// ExitShowBrokersStmt is called when production showBrokersStmt is exited.
func (s *BaseSQLListener) ExitShowBrokersStmt(ctx *ShowBrokersStmtContext) {}

// EnterShowRequestsStmt is called when production showRequestsStmt is entered.
func (s *BaseSQLListener) EnterShowRequestsStmt(ctx *ShowRequestsStmtContext) {}

// ExitShowRequestsStmt is called when production showRequestsStmt is exited.
func (s *BaseSQLListener) ExitShowRequestsStmt(ctx *ShowRequestsStmtContext) {}

// EnterRecoverStorageStmt is called when production recoverStorageStmt is entered.
func (s *BaseSQLListener) EnterRecoverStorageStmt(ctx *RecoverStorageStmtContext) {}

// ExitRecoverStorageStmt is called when production recoverStorageStmt is exited.
func (s *BaseSQLListener) ExitRecoverStorageStmt(ctx *RecoverStorageStmtContext) {}

// EnterShowLimitStmt is called when production showLimitStmt is entered.
func (s *BaseSQLListener) EnterShowLimitStmt(ctx *ShowLimitStmtContext) {}

// ExitShowLimitStmt is called when production showLimitStmt is exited.
func (s *BaseSQLListener) ExitShowLimitStmt(ctx *ShowLimitStmtContext) {}

// EnterShowMetadataTypesStmt is called when production showMetadataTypesStmt is entered.
func (s *BaseSQLListener) EnterShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) {}

// ExitShowMetadataTypesStmt is called when production showMetadataTypesStmt is exited.
func (s *BaseSQLListener) ExitShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) {}

// EnterShowMetadatasStmt is called when production showMetadatasStmt is entered.
func (s *BaseSQLListener) EnterShowMetadatasStmt(ctx *ShowMetadatasStmtContext) {}

// ExitShowMetadatasStmt is called when production showMetadatasStmt is exited.
func (s *BaseSQLListener) ExitShowMetadatasStmt(ctx *ShowMetadatasStmtContext) {}

// EnterShowAliveStmt is called when production showAliveStmt is entered.
func (s *BaseSQLListener) EnterShowAliveStmt(ctx *ShowAliveStmtContext) {}

// ExitShowAliveStmt is called when production showAliveStmt is exited.
func (s *BaseSQLListener) ExitShowAliveStmt(ctx *ShowAliveStmtContext) {}

// EnterShowReplicationsStmt is called when production showReplicationsStmt is entered.
func (s *BaseSQLListener) EnterShowReplicationsStmt(ctx *ShowReplicationsStmtContext) {}

// ExitShowReplicationsStmt is called when production showReplicationsStmt is exited.
func (s *BaseSQLListener) ExitShowReplicationsStmt(ctx *ShowReplicationsStmtContext) {}

// EnterShowStateStmt is called when production showStateStmt is entered.
func (s *BaseSQLListener) EnterShowStateStmt(ctx *ShowStateStmtContext) {}

// ExitShowStateStmt is called when production showStateStmt is exited.
func (s *BaseSQLListener) ExitShowStateStmt(ctx *ShowStateStmtContext) {}

// EnterCreateDatabaseStmt is called when production createDatabaseStmt is entered.
func (s *BaseSQLListener) EnterCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) {}

// ExitCreateDatabaseStmt is called when production createDatabaseStmt is exited.
func (s *BaseSQLListener) ExitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) {}

// EnterDropDatabaseStmt is called when production dropDatabaseStmt is entered.
func (s *BaseSQLListener) EnterDropDatabaseStmt(ctx *DropDatabaseStmtContext) {}

// ExitDropDatabaseStmt is called when production dropDatabaseStmt is exited.
func (s *BaseSQLListener) ExitDropDatabaseStmt(ctx *DropDatabaseStmtContext) {}

// EnterShowDatabasesStmt is called when production showDatabasesStmt is entered.
func (s *BaseSQLListener) EnterShowDatabasesStmt(ctx *ShowDatabasesStmtContext) {}

// ExitShowDatabasesStmt is called when production showDatabasesStmt is exited.
func (s *BaseSQLListener) ExitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) {}

// EnterCreateBrokerStmt is called when production createBrokerStmt is entered.
func (s *BaseSQLListener) EnterCreateBrokerStmt(ctx *CreateBrokerStmtContext) {}

// ExitCreateBrokerStmt is called when production createBrokerStmt is exited.
func (s *BaseSQLListener) ExitCreateBrokerStmt(ctx *CreateBrokerStmtContext) {}

// EnterShowNamespacesStmt is called when production showNamespacesStmt is entered.
func (s *BaseSQLListener) EnterShowNamespacesStmt(ctx *ShowNamespacesStmtContext) {}

// ExitShowNamespacesStmt is called when production showNamespacesStmt is exited.
func (s *BaseSQLListener) ExitShowNamespacesStmt(ctx *ShowNamespacesStmtContext) {}

// EnterShowMetricsStmt is called when production showMetricsStmt is entered.
func (s *BaseSQLListener) EnterShowMetricsStmt(ctx *ShowMetricsStmtContext) {}

// ExitShowMetricsStmt is called when production showMetricsStmt is exited.
func (s *BaseSQLListener) ExitShowMetricsStmt(ctx *ShowMetricsStmtContext) {}

// EnterShowFieldsStmt is called when production showFieldsStmt is entered.
func (s *BaseSQLListener) EnterShowFieldsStmt(ctx *ShowFieldsStmtContext) {}

// ExitShowFieldsStmt is called when production showFieldsStmt is exited.
func (s *BaseSQLListener) ExitShowFieldsStmt(ctx *ShowFieldsStmtContext) {}

// EnterShowTagKeysStmt is called when production showTagKeysStmt is entered.
func (s *BaseSQLListener) EnterShowTagKeysStmt(ctx *ShowTagKeysStmtContext) {}

// ExitShowTagKeysStmt is called when production showTagKeysStmt is exited.
func (s *BaseSQLListener) ExitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) {}

// EnterShowTagValuesStmt is called when production showTagValuesStmt is entered.
func (s *BaseSQLListener) EnterShowTagValuesStmt(ctx *ShowTagValuesStmtContext) {}

// ExitShowTagValuesStmt is called when production showTagValuesStmt is exited.
func (s *BaseSQLListener) ExitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) {}

// EnterQueryStmt is called when production queryStmt is entered.
func (s *BaseSQLListener) EnterQueryStmt(ctx *QueryStmtContext) {}

// ExitQueryStmt is called when production queryStmt is exited.
func (s *BaseSQLListener) ExitQueryStmt(ctx *QueryStmtContext) {}

// EnterSourceAndSelect is called when production sourceAndSelect is entered.
func (s *BaseSQLListener) EnterSourceAndSelect(ctx *SourceAndSelectContext) {}

// ExitSourceAndSelect is called when production sourceAndSelect is exited.
func (s *BaseSQLListener) ExitSourceAndSelect(ctx *SourceAndSelectContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (s *BaseSQLListener) EnterSelectExpr(ctx *SelectExprContext) {}

// ExitSelectExpr is called when production selectExpr is exited.
func (s *BaseSQLListener) ExitSelectExpr(ctx *SelectExprContext) {}

// EnterFields is called when production fields is entered.
func (s *BaseSQLListener) EnterFields(ctx *FieldsContext) {}

// ExitFields is called when production fields is exited.
func (s *BaseSQLListener) ExitFields(ctx *FieldsContext) {}

// EnterField is called when production field is entered.
func (s *BaseSQLListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseSQLListener) ExitField(ctx *FieldContext) {}

// EnterAlias is called when production alias is entered.
func (s *BaseSQLListener) EnterAlias(ctx *AliasContext) {}

// ExitAlias is called when production alias is exited.
func (s *BaseSQLListener) ExitAlias(ctx *AliasContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BaseSQLListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BaseSQLListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseSQLListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseSQLListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterConditionExpr is called when production conditionExpr is entered.
func (s *BaseSQLListener) EnterConditionExpr(ctx *ConditionExprContext) {}

// ExitConditionExpr is called when production conditionExpr is exited.
func (s *BaseSQLListener) ExitConditionExpr(ctx *ConditionExprContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseSQLListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseSQLListener) ExitExpression(ctx *ExpressionContext) {}

// EnterValueList is called when production valueList is entered.
func (s *BaseSQLListener) EnterValueList(ctx *ValueListContext) {}

// ExitValueList is called when production valueList is exited.
func (s *BaseSQLListener) ExitValueList(ctx *ValueListContext) {}

// EnterTimeExpr is called when production timeExpr is entered.
func (s *BaseSQLListener) EnterTimeExpr(ctx *TimeExprContext) {}

// ExitTimeExpr is called when production timeExpr is exited.
func (s *BaseSQLListener) ExitTimeExpr(ctx *TimeExprContext) {}

// EnterNowExpr is called when production nowExpr is entered.
func (s *BaseSQLListener) EnterNowExpr(ctx *NowExprContext) {}

// ExitNowExpr is called when production nowExpr is exited.
func (s *BaseSQLListener) ExitNowExpr(ctx *NowExprContext) {}

// EnterNowFunc is called when production nowFunc is entered.
func (s *BaseSQLListener) EnterNowFunc(ctx *NowFuncContext) {}

// ExitNowFunc is called when production nowFunc is exited.
func (s *BaseSQLListener) ExitNowFunc(ctx *NowFuncContext) {}

// EnterGroupByClause is called when production groupByClause is entered.
func (s *BaseSQLListener) EnterGroupByClause(ctx *GroupByClauseContext) {}

// ExitGroupByClause is called when production groupByClause is exited.
func (s *BaseSQLListener) ExitGroupByClause(ctx *GroupByClauseContext) {}

// EnterGroupByKeys is called when production groupByKeys is entered.
func (s *BaseSQLListener) EnterGroupByKeys(ctx *GroupByKeysContext) {}

// ExitGroupByKeys is called when production groupByKeys is exited.
func (s *BaseSQLListener) ExitGroupByKeys(ctx *GroupByKeysContext) {}

// EnterGroupByKey is called when production groupByKey is entered.
func (s *BaseSQLListener) EnterGroupByKey(ctx *GroupByKeyContext) {}

// ExitGroupByKey is called when production groupByKey is exited.
func (s *BaseSQLListener) ExitGroupByKey(ctx *GroupByKeyContext) {}

// EnterFillOption is called when production fillOption is entered.
func (s *BaseSQLListener) EnterFillOption(ctx *FillOptionContext) {}

// ExitFillOption is called when production fillOption is exited.
func (s *BaseSQLListener) ExitFillOption(ctx *FillOptionContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseSQLListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseSQLListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterSortField is called when production sortField is entered.
func (s *BaseSQLListener) EnterSortField(ctx *SortFieldContext) {}

// ExitSortField is called when production sortField is exited.
func (s *BaseSQLListener) ExitSortField(ctx *SortFieldContext) {}

// EnterSortFields is called when production sortFields is entered.
func (s *BaseSQLListener) EnterSortFields(ctx *SortFieldsContext) {}

// ExitSortFields is called when production sortFields is exited.
func (s *BaseSQLListener) ExitSortFields(ctx *SortFieldsContext) {}

// EnterHavingClause is called when production havingClause is entered.
func (s *BaseSQLListener) EnterHavingClause(ctx *HavingClauseContext) {}

// ExitHavingClause is called when production havingClause is exited.
func (s *BaseSQLListener) ExitHavingClause(ctx *HavingClauseContext) {}

// EnterBoolExpr is called when production boolExpr is entered.
func (s *BaseSQLListener) EnterBoolExpr(ctx *BoolExprContext) {}

// ExitBoolExpr is called when production boolExpr is exited.
func (s *BaseSQLListener) ExitBoolExpr(ctx *BoolExprContext) {}

// EnterBoolExprLogicalOp is called when production boolExprLogicalOp is entered.
func (s *BaseSQLListener) EnterBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) {}

// ExitBoolExprLogicalOp is called when production boolExprLogicalOp is exited.
func (s *BaseSQLListener) ExitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) {}

// EnterBoolExprAtom is called when production boolExprAtom is entered.
func (s *BaseSQLListener) EnterBoolExprAtom(ctx *BoolExprAtomContext) {}

// ExitBoolExprAtom is called when production boolExprAtom is exited.
func (s *BaseSQLListener) ExitBoolExprAtom(ctx *BoolExprAtomContext) {}

// EnterBinaryExpr is called when production binaryExpr is entered.
func (s *BaseSQLListener) EnterBinaryExpr(ctx *BinaryExprContext) {}

// ExitBinaryExpr is called when production binaryExpr is exited.
func (s *BaseSQLListener) ExitBinaryExpr(ctx *BinaryExprContext) {}

// EnterBinaryOperator is called when production binaryOperator is entered.
func (s *BaseSQLListener) EnterBinaryOperator(ctx *BinaryOperatorContext) {}

// ExitBinaryOperator is called when production binaryOperator is exited.
func (s *BaseSQLListener) ExitBinaryOperator(ctx *BinaryOperatorContext) {}

// EnterFieldExpr is called when production fieldExpr is entered.
func (s *BaseSQLListener) EnterFieldExpr(ctx *FieldExprContext) {}

// ExitFieldExpr is called when production fieldExpr is exited.
func (s *BaseSQLListener) ExitFieldExpr(ctx *FieldExprContext) {}

// EnterStar is called when production star is entered.
func (s *BaseSQLListener) EnterStar(ctx *StarContext) {}

// ExitStar is called when production star is exited.
func (s *BaseSQLListener) ExitStar(ctx *StarContext) {}

// EnterDurationLit is called when production durationLit is entered.
func (s *BaseSQLListener) EnterDurationLit(ctx *DurationLitContext) {}

// ExitDurationLit is called when production durationLit is exited.
func (s *BaseSQLListener) ExitDurationLit(ctx *DurationLitContext) {}

// EnterIntervalItem is called when production intervalItem is entered.
func (s *BaseSQLListener) EnterIntervalItem(ctx *IntervalItemContext) {}

// ExitIntervalItem is called when production intervalItem is exited.
func (s *BaseSQLListener) ExitIntervalItem(ctx *IntervalItemContext) {}

// EnterExprFunc is called when production exprFunc is entered.
func (s *BaseSQLListener) EnterExprFunc(ctx *ExprFuncContext) {}

// ExitExprFunc is called when production exprFunc is exited.
func (s *BaseSQLListener) ExitExprFunc(ctx *ExprFuncContext) {}

// EnterFuncName is called when production funcName is entered.
func (s *BaseSQLListener) EnterFuncName(ctx *FuncNameContext) {}

// ExitFuncName is called when production funcName is exited.
func (s *BaseSQLListener) ExitFuncName(ctx *FuncNameContext) {}

// EnterExprFuncParams is called when production exprFuncParams is entered.
func (s *BaseSQLListener) EnterExprFuncParams(ctx *ExprFuncParamsContext) {}

// ExitExprFuncParams is called when production exprFuncParams is exited.
func (s *BaseSQLListener) ExitExprFuncParams(ctx *ExprFuncParamsContext) {}

// EnterFuncParam is called when production funcParam is entered.
func (s *BaseSQLListener) EnterFuncParam(ctx *FuncParamContext) {}

// ExitFuncParam is called when production funcParam is exited.
func (s *BaseSQLListener) ExitFuncParam(ctx *FuncParamContext) {}

// EnterExprAtom is called when production exprAtom is entered.
func (s *BaseSQLListener) EnterExprAtom(ctx *ExprAtomContext) {}

// ExitExprAtom is called when production exprAtom is exited.
func (s *BaseSQLListener) ExitExprAtom(ctx *ExprAtomContext) {}

// EnterProperties is called when production properties is entered.
func (s *BaseSQLListener) EnterProperties(ctx *PropertiesContext) {}

// ExitProperties is called when production properties is exited.
func (s *BaseSQLListener) ExitProperties(ctx *PropertiesContext) {}

// EnterPropertyAssignments is called when production propertyAssignments is entered.
func (s *BaseSQLListener) EnterPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// ExitPropertyAssignments is called when production propertyAssignments is exited.
func (s *BaseSQLListener) ExitPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// EnterProperty is called when production property is entered.
func (s *BaseSQLListener) EnterProperty(ctx *PropertyContext) {}

// ExitProperty is called when production property is exited.
func (s *BaseSQLListener) ExitProperty(ctx *PropertyContext) {}

// EnterValue is called when production value is entered.
func (s *BaseSQLListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BaseSQLListener) ExitValue(ctx *ValueContext) {}

// EnterIntNumber is called when production intNumber is entered.
func (s *BaseSQLListener) EnterIntNumber(ctx *IntNumberContext) {}

// ExitIntNumber is called when production intNumber is exited.
func (s *BaseSQLListener) ExitIntNumber(ctx *IntNumberContext) {}

// EnterDecNumber is called when production decNumber is entered.
func (s *BaseSQLListener) EnterDecNumber(ctx *DecNumberContext) {}

// ExitDecNumber is called when production decNumber is exited.
func (s *BaseSQLListener) ExitDecNumber(ctx *DecNumberContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseSQLListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseSQLListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterMetricName is called when production metricName is entered.
func (s *BaseSQLListener) EnterMetricName(ctx *MetricNameContext) {}

// ExitMetricName is called when production metricName is exited.
func (s *BaseSQLListener) ExitMetricName(ctx *MetricNameContext) {}

// EnterTagKey is called when production tagKey is entered.
func (s *BaseSQLListener) EnterTagKey(ctx *TagKeyContext) {}

// ExitTagKey is called when production tagKey is exited.
func (s *BaseSQLListener) ExitTagKey(ctx *TagKeyContext) {}

// EnterTagValue is called when production tagValue is entered.
func (s *BaseSQLListener) EnterTagValue(ctx *TagValueContext) {}

// ExitTagValue is called when production tagValue is exited.
func (s *BaseSQLListener) ExitTagValue(ctx *TagValueContext) {}

// EnterPrefix is called when production prefix is entered.
func (s *BaseSQLListener) EnterPrefix(ctx *PrefixContext) {}

// ExitPrefix is called when production prefix is exited.
func (s *BaseSQLListener) ExitPrefix(ctx *PrefixContext) {}

// EnterWithTagKey is called when production withTagKey is entered.
func (s *BaseSQLListener) EnterWithTagKey(ctx *WithTagKeyContext) {}

// ExitWithTagKey is called when production withTagKey is exited.
func (s *BaseSQLListener) ExitWithTagKey(ctx *WithTagKeyContext) {}

// EnterNamespace is called when production namespace is entered.
func (s *BaseSQLListener) EnterNamespace(ctx *NamespaceContext) {}

// ExitNamespace is called when production namespace is exited.
func (s *BaseSQLListener) ExitNamespace(ctx *NamespaceContext) {}

// EnterName is called when production name is entered.
func (s *BaseSQLListener) EnterName(ctx *NameContext) {}

// ExitName is called when production name is exited.
func (s *BaseSQLListener) ExitName(ctx *NameContext) {}

// EnterRequestID is called when production requestID is entered.
func (s *BaseSQLListener) EnterRequestID(ctx *RequestIDContext) {}

// ExitRequestID is called when production requestID is exited.
func (s *BaseSQLListener) ExitRequestID(ctx *RequestIDContext) {}

// EnterToml is called when production toml is entered.
func (s *BaseSQLListener) EnterToml(ctx *TomlContext) {}

// ExitToml is called when production toml is exited.
func (s *BaseSQLListener) ExitToml(ctx *TomlContext) {}

// EnterIdent is called when production ident is entered.
func (s *BaseSQLListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BaseSQLListener) ExitIdent(ctx *IdentContext) {}

// EnterNonReservedWords is called when production nonReservedWords is entered.
func (s *BaseSQLListener) EnterNonReservedWords(ctx *NonReservedWordsContext) {}

// ExitNonReservedWords is called when production nonReservedWords is exited.
func (s *BaseSQLListener) ExitNonReservedWords(ctx *NonReservedWordsContext) {}
