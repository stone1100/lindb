// Code generated from ./sql/grammar/SQL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by SQLParser.
type SQLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SQLParser#useStmt.
	VisitUseStmt(ctx *UseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#setLimitStmt.
	VisitSetLimitStmt(ctx *SetLimitStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStmt.
	VisitShowStmt(ctx *ShowStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMasterStmt.
	VisitShowMasterStmt(ctx *ShowMasterStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showBrokersStmt.
	VisitShowBrokersStmt(ctx *ShowBrokersStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showRequestsStmt.
	VisitShowRequestsStmt(ctx *ShowRequestsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#recoverStorageStmt.
	VisitRecoverStorageStmt(ctx *RecoverStorageStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showLimitStmt.
	VisitShowLimitStmt(ctx *ShowLimitStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetadataTypesStmt.
	VisitShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetadatasStmt.
	VisitShowMetadatasStmt(ctx *ShowMetadatasStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showAliveStmt.
	VisitShowAliveStmt(ctx *ShowAliveStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showReplicationsStmt.
	VisitShowReplicationsStmt(ctx *ShowReplicationsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStateStmt.
	VisitShowStateStmt(ctx *ShowStateStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#createDatabaseStmt.
	VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#dropDatabaseStmt.
	VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showDatabasesStmt.
	VisitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#createBrokerStmt.
	VisitCreateBrokerStmt(ctx *CreateBrokerStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showNamespacesStmt.
	VisitShowNamespacesStmt(ctx *ShowNamespacesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetricsStmt.
	VisitShowMetricsStmt(ctx *ShowMetricsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showFieldsStmt.
	VisitShowFieldsStmt(ctx *ShowFieldsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagKeysStmt.
	VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagValuesStmt.
	VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#queryStmt.
	VisitQueryStmt(ctx *QueryStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#sourceAndSelect.
	VisitSourceAndSelect(ctx *SourceAndSelectContext) interface{}

	// Visit a parse tree produced by SQLParser#selectExpr.
	VisitSelectExpr(ctx *SelectExprContext) interface{}

	// Visit a parse tree produced by SQLParser#fields.
	VisitFields(ctx *FieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by SQLParser#alias.
	VisitAlias(ctx *AliasContext) interface{}

	// Visit a parse tree produced by SQLParser#fromClause.
	VisitFromClause(ctx *FromClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#whereClause.
	VisitWhereClause(ctx *WhereClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#conditionExpr.
	VisitConditionExpr(ctx *ConditionExprContext) interface{}

	// Visit a parse tree produced by SQLParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by SQLParser#valueList.
	VisitValueList(ctx *ValueListContext) interface{}

	// Visit a parse tree produced by SQLParser#timeExpr.
	VisitTimeExpr(ctx *TimeExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowExpr.
	VisitNowExpr(ctx *NowExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowFunc.
	VisitNowFunc(ctx *NowFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByClause.
	VisitGroupByClause(ctx *GroupByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByKeys.
	VisitGroupByKeys(ctx *GroupByKeysContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByKey.
	VisitGroupByKey(ctx *GroupByKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#fillOption.
	VisitFillOption(ctx *FillOptionContext) interface{}

	// Visit a parse tree produced by SQLParser#orderByClause.
	VisitOrderByClause(ctx *OrderByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#sortField.
	VisitSortField(ctx *SortFieldContext) interface{}

	// Visit a parse tree produced by SQLParser#sortFields.
	VisitSortFields(ctx *SortFieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#havingClause.
	VisitHavingClause(ctx *HavingClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExpr.
	VisitBoolExpr(ctx *BoolExprContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprLogicalOp.
	VisitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprAtom.
	VisitBoolExprAtom(ctx *BoolExprAtomContext) interface{}

	// Visit a parse tree produced by SQLParser#binaryExpr.
	VisitBinaryExpr(ctx *BinaryExprContext) interface{}

	// Visit a parse tree produced by SQLParser#binaryOperator.
	VisitBinaryOperator(ctx *BinaryOperatorContext) interface{}

	// Visit a parse tree produced by SQLParser#fieldExpr.
	VisitFieldExpr(ctx *FieldExprContext) interface{}

	// Visit a parse tree produced by SQLParser#star.
	VisitStar(ctx *StarContext) interface{}

	// Visit a parse tree produced by SQLParser#durationLit.
	VisitDurationLit(ctx *DurationLitContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalItem.
	VisitIntervalItem(ctx *IntervalItemContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFunc.
	VisitExprFunc(ctx *ExprFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#funcName.
	VisitFuncName(ctx *FuncNameContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFuncParams.
	VisitExprFuncParams(ctx *ExprFuncParamsContext) interface{}

	// Visit a parse tree produced by SQLParser#funcParam.
	VisitFuncParam(ctx *FuncParamContext) interface{}

	// Visit a parse tree produced by SQLParser#exprAtom.
	VisitExprAtom(ctx *ExprAtomContext) interface{}

	// Visit a parse tree produced by SQLParser#properties.
	VisitProperties(ctx *PropertiesContext) interface{}

	// Visit a parse tree produced by SQLParser#propertyAssignments.
	VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{}

	// Visit a parse tree produced by SQLParser#property.
	VisitProperty(ctx *PropertyContext) interface{}

	// Visit a parse tree produced by SQLParser#value.
	VisitValue(ctx *ValueContext) interface{}

	// Visit a parse tree produced by SQLParser#intNumber.
	VisitIntNumber(ctx *IntNumberContext) interface{}

	// Visit a parse tree produced by SQLParser#decNumber.
	VisitDecNumber(ctx *DecNumberContext) interface{}

	// Visit a parse tree produced by SQLParser#limitClause.
	VisitLimitClause(ctx *LimitClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#metricName.
	VisitMetricName(ctx *MetricNameContext) interface{}

	// Visit a parse tree produced by SQLParser#tagKey.
	VisitTagKey(ctx *TagKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#tagValue.
	VisitTagValue(ctx *TagValueContext) interface{}

	// Visit a parse tree produced by SQLParser#prefix.
	VisitPrefix(ctx *PrefixContext) interface{}

	// Visit a parse tree produced by SQLParser#withTagKey.
	VisitWithTagKey(ctx *WithTagKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#namespace.
	VisitNamespace(ctx *NamespaceContext) interface{}

	// Visit a parse tree produced by SQLParser#name.
	VisitName(ctx *NameContext) interface{}

	// Visit a parse tree produced by SQLParser#requestID.
	VisitRequestID(ctx *RequestIDContext) interface{}

	// Visit a parse tree produced by SQLParser#toml.
	VisitToml(ctx *TomlContext) interface{}

	// Visit a parse tree produced by SQLParser#ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by SQLParser#nonReservedWords.
	VisitNonReservedWords(ctx *NonReservedWordsContext) interface{}
}
