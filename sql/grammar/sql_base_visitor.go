// Code generated from ./sql/grammar/SQL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr4-go/antlr/v4"

type BaseSQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitUseStmt(ctx *UseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSetLimitStmt(ctx *SetLimitStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStmt(ctx *ShowStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMasterStmt(ctx *ShowMasterStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowBrokersStmt(ctx *ShowBrokersStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowRequestsStmt(ctx *ShowRequestsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitRecoverStorageStmt(ctx *RecoverStorageStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowLimitStmt(ctx *ShowLimitStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMetadatasStmt(ctx *ShowMetadatasStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowAliveStmt(ctx *ShowAliveStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowReplicationsStmt(ctx *ShowReplicationsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStateStmt(ctx *ShowStateStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreateBrokerStmt(ctx *CreateBrokerStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowNamespacesStmt(ctx *ShowNamespacesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMetricsStmt(ctx *ShowMetricsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowFieldsStmt(ctx *ShowFieldsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQueryStmt(ctx *QueryStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSourceAndSelect(ctx *SourceAndSelectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSelectExpr(ctx *SelectExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFields(ctx *FieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitField(ctx *FieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitAlias(ctx *AliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFromClause(ctx *FromClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitConditionExpr(ctx *ConditionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitValueList(ctx *ValueListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTimeExpr(ctx *TimeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNowExpr(ctx *NowExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNowFunc(ctx *NowFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroupByClause(ctx *GroupByClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroupByKeys(ctx *GroupByKeysContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroupByKey(ctx *GroupByKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFillOption(ctx *FillOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitOrderByClause(ctx *OrderByClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSortField(ctx *SortFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSortFields(ctx *SortFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitHavingClause(ctx *HavingClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExpr(ctx *BoolExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprAtom(ctx *BoolExprAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBinaryExpr(ctx *BinaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBinaryOperator(ctx *BinaryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFieldExpr(ctx *FieldExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitStar(ctx *StarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDurationLit(ctx *DurationLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalItem(ctx *IntervalItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprFunc(ctx *ExprFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFuncName(ctx *FuncNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprFuncParams(ctx *ExprFuncParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFuncParam(ctx *FuncParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprAtom(ctx *ExprAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitProperty(ctx *PropertyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitValue(ctx *ValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntNumber(ctx *IntNumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDecNumber(ctx *DecNumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetricName(ctx *MetricNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagKey(ctx *TagKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagValue(ctx *TagValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPrefix(ctx *PrefixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithTagKey(ctx *WithTagKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitName(ctx *NameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitRequestID(ctx *RequestIDContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitToml(ctx *TomlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNonReservedWords(ctx *NonReservedWordsContext) interface{} {
	return v.VisitChildren(ctx)
}
