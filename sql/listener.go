// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package sql

import (
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

type listener struct {
	*grammar.BaseSQLListener

	queryStmt          *queryStmtParser
	metadataStmt       *metadataStmtParser
	stateStmt          *stateStmtParser
	metricMetadataStmt *metricMetadataStmtParser
	useStmt            *useStmtParser
	schemasStmt        *schemasStmtParser
	storageStmt        *storageStmtParser
	requestStmt        *requestStmtParser
	brokerStmt         *brokerStmtParser
	limitStmt          *limitStmtParser
}

// EnterQueryStmt is called when production queryStmt is entered.
func (l *listener) EnterQueryStmt(ctx *grammar.QueryStmtContext) {
	l.queryStmt = newQueryStmtParse(ctx.T_EXPLAIN() != nil)
}

// EnterShowMetadataTypesStmt is called when production showMetadataTypesStmt is entered.
func (l *listener) EnterShowMetadataTypesStmt(_ *grammar.ShowMetadataTypesStmtContext) {
	l.metadataStmt = newMetadataStmtParser(stmt.MetadataTypes)
}

// EnterShowRequestsStmt is called when production showRequestssStmt is entered.
func (l *listener) EnterShowRequestsStmt(_ *grammar.ShowRequestsStmtContext) {
	l.requestStmt = newRequestStmtParse()
}

// EnterRequestID is called when production requestID is entered.
func (l *listener) EnterRequestID(ctx *grammar.RequestIDContext) {
	l.requestStmt.visitRequestID(ctx)
}

// EnterShowMasterStmt is called when production showMasterStmt is entered.
func (l *listener) EnterShowMasterStmt(_ *grammar.ShowMasterStmtContext) {
	l.stateStmt = newStateStmtParse(stmt.Master)
}

// EnterShowAliveStmt is called when production showAliveStmt is entered.
func (l *listener) EnterShowAliveStmt(ctx *grammar.ShowAliveStmtContext) {
	switch {
	case ctx.T_ROOT() != nil:
		l.stateStmt = newStateStmtParse(stmt.RootAlive)
	case ctx.T_BROKER() != nil:
		l.stateStmt = newStateStmtParse(stmt.BrokerAlive)
	case ctx.T_STORAGE() != nil:
		l.stateStmt = newStateStmtParse(stmt.StorageAlive)
	}
}

// EnterRecoverStorageStmt is called when entering the recoverStorageStmt production.
func (l *listener) EnterRecoverStorageStmt(_ *grammar.RecoverStorageStmtContext) {
	l.storageStmt = newStorageStmtParse(stmt.StorageOpRecover)
}

// EnterDatabaseFilter is called when production databaseFilter is entered.
func (l *listener) EnterDatabaseFilter(ctx *grammar.DatabaseFilterContext) {
	l.stateStmt.visitDatabaseFilter(ctx)
}

// EnterShowBrokersStmt is called when production showBrokersStmt is entered.
func (l *listener) EnterShowBrokersStmt(_ *grammar.ShowBrokersStmtContext) {
	l.brokerStmt = newBrokerStmtParse(stmt.BrokerOpShow)
}

// EnterOptionClause is called when production optionClause is entered.
func (l *listener) EnterOptionClause(ctx *grammar.OptionClauseContext) {
	if l.schemasStmt != nil {
		l.schemasStmt.visitWithCfg(ctx)
	}
}

// EnterCreateBrokerStmt is called when production createBrokerStmt is entered.
func (l *listener) EnterCreateBrokerStmt(c *grammar.CreateBrokerStmtContext) {
	l.brokerStmt = newBrokerStmtParse(stmt.BrokerOpCreate)
}

// EnterCreateDatabaseStmt is called when entering the createDatabaseStmt production.
func (l *listener) EnterCreateDatabaseStmt(_ *grammar.CreateDatabaseStmtContext) {
	l.schemasStmt = newSchemasStmtParse(stmt.CreateDatabaseSchemaType)
}

// EnterDropDatabaseStmt is called when production dropDatabaseStmt is entered.
func (l *listener) EnterDropDatabaseStmt(_ *grammar.DropDatabaseStmtContext) {
	l.schemasStmt = newSchemasStmtParse(stmt.DropDatabaseSchemaType)
}

// EnterDatabaseName is called when production databaseName is entered.
func (l *listener) EnterDatabaseName(ctx *grammar.DatabaseNameContext) {
	l.schemasStmt.visitDatabaseName(ctx)
}

// EnterSetLimtStmt is called when production setLimitStmt is entered.
func (l *listener) EnterSetLimitStmt(ctx *grammar.SetLimitStmtContext) {
	l.limitStmt = newLimitStmtParse(stmt.SetLimit)
	l.limitStmt.visitToml(ctx.Toml())
}

// EnterShowLimtStmt is called when production showLimitStmt is entered.
func (l *listener) EnterShowLimitStmt(ctx *grammar.ShowLimitStmtContext) {
	l.limitStmt = newLimitStmtParse(stmt.ShowLimit)
}

// EnterUseStmt is called when production useStmt is entered.
func (l *listener) EnterUseStmt(ctx *grammar.UseStmtContext) {
	l.useStmt = newUseStmtParse()
	l.useStmt.visitName(ctx.Ident())
}

// EnterShowDatabaseStmt is called when production showDatabaseStmt is entered.
func (l *listener) EnterShowDatabasesStmt(_ *grammar.ShowDatabasesStmtContext) {
	l.schemasStmt = newSchemasStmtParse(stmt.DatabaseNameSchemaType)
}

// EnterShowNameSpacesStmt is called when production showNameSpacesStmt is entered.
func (l *listener) EnterShowNamespacesStmt(_ *grammar.ShowNamespacesStmtContext) {
	l.metricMetadataStmt = newMetricMetadataStmtParser(stmt.Namespace)
}

// EnterShowMetricsStmt is called when production showMetricsStmt is entered.
func (l *listener) EnterShowMetricsStmt(_ *grammar.ShowMetricsStmtContext) {
	l.metricMetadataStmt = newMetricMetadataStmtParser(stmt.Metric)
}

// EnterShowFieldsStmt is called when production showFieldsStmt is entered.
func (l *listener) EnterShowFieldsStmt(_ *grammar.ShowFieldsStmtContext) {
	l.metricMetadataStmt = newMetricMetadataStmtParser(stmt.Field)
}

// EnterShowTagKeysStmt is called when production showTagKeysStmt is entered.
func (l *listener) EnterShowTagKeysStmt(_ *grammar.ShowTagKeysStmtContext) {
	l.metricMetadataStmt = newMetricMetadataStmtParser(stmt.TagKey)
}

// EnterShowTagValuesStmt is called when production showTagValuesStmt is entered.
func (l *listener) EnterShowTagValuesStmt(_ *grammar.ShowTagValuesStmtContext) {
	l.metricMetadataStmt = newMetricMetadataStmtParser(stmt.TagValue)
}

// EnterNamespace is called when production namespace is entered.
func (l *listener) EnterNamespace(ctx *grammar.NamespaceContext) {
	switch {
	case l.queryStmt != nil:
		l.queryStmt.visitNamespace(ctx)
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.visitNamespace(ctx)
	}
}

// EnterWithTagKey is called when production withTagKey is entered.
func (l *listener) EnterWithTagKey(ctx *grammar.WithTagKeyContext) {
	if l.metricMetadataStmt != nil {
		l.metricMetadataStmt.visitWithTagKey(ctx)
	}
}

// EnterPrefix is called when production prefix is entered.
func (l *listener) EnterPrefix(ctx *grammar.PrefixContext) {
	if l.metricMetadataStmt != nil {
		l.metricMetadataStmt.visitPrefix(ctx)
	}
}

// EnterMetricName is called when production metricName is entered.
func (l *listener) EnterMetricName(ctx *grammar.MetricNameContext) {
	switch {
	case l.queryStmt != nil:
		l.queryStmt.visitMetricName(ctx)
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.visitMetricName(ctx)
	}
}

// EnterSelectExpr is called when production selectExpr is entered.
func (l *listener) EnterSelectExpr(_ *grammar.SelectExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.resetExprStack()
	}
}

// EnterWhereClause is called when production whereClause is entered.
func (l *listener) EnterWhereClause(_ *grammar.WhereClauseContext) {
	if l.queryStmt != nil {
		l.queryStmt.resetExprStack()
	}
}

// EnterFieldExpr is called when production fieldExpr is entered.
func (l *listener) EnterFieldExpr(ctx *grammar.FieldExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitFieldExpr(ctx)
	}
}

// ExitFieldExpr is called when production fieldExpr is exited.
func (l *listener) ExitFieldExpr(ctx *grammar.FieldExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.completeFieldExpr(ctx)
	}
}

// EnterFuncName is called when production exprFunc is entered.
func (l *listener) EnterFuncName(ctx *grammar.FuncNameContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitFuncName(ctx)
	}
}

// ExitExprFunc is called when production exprFunc is exited.
func (l *listener) ExitExprFunc(_ *grammar.ExprFuncContext) {
	if l.queryStmt != nil {
		l.queryStmt.completeFuncExpr()
	}
}

// EnterExprAtom is called when production exprAtom is entered.
func (l *listener) EnterExprAtom(ctx *grammar.ExprAtomContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitExprAtom(ctx)
	}
}

// EnterAlias is called when production alias is entered.
func (l *listener) EnterAlias(ctx *grammar.AliasContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitAlias(ctx)
	}
}

// EnterLimitClause is called when production limitClause is entered.
func (l *listener) EnterLimitClause(ctx *grammar.LimitClauseContext) {
	switch {
	case l.queryStmt != nil:
		l.queryStmt.visitLimit(ctx)
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.visitLimit(ctx)
	}
}

// EnterTagFilterExpr is called when production tagFilterExpr is entered.
func (l *listener) EnterExpression(ctx *grammar.ExpressionContext) {
	switch {
	case l.stateStmt != nil:
		l.stateStmt.visitExpression(ctx)
	case l.metadataStmt != nil:
		l.metadataStmt.visitExpression(ctx)
	case l.queryStmt != nil:
		l.queryStmt.visitExpression(ctx)
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.visitExpression(ctx)
	}
}

// ExitTagFilterExpr is called when production tagValueList is exited.
func (l *listener) ExitExpression(_ *grammar.ExpressionContext) {
	switch {
	case l.stateStmt != nil:
		l.stateStmt.completeExpression()
	case l.metadataStmt != nil:
		l.metadataStmt.completeExpression()
	case l.queryStmt != nil:
		l.queryStmt.completeExpression()
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.completeExpression()
	}
}

// EnterTagValue is called when production tagValue is entered.
func (l *listener) EnterTagValue(ctx *grammar.TagValueContext) {
	switch {
	case l.queryStmt != nil:
		l.queryStmt.visitTagValue(ctx)
	case l.metricMetadataStmt != nil:
		l.metricMetadataStmt.visitTagValue(ctx)
	}
}

// EnterTimeRangeExpr is called when production timeRangeExpr is entered.
func (l *listener) EnterTimeRangeExpr(ctx *grammar.TimeRangeExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitTimeRangeExpr(ctx)
	}
}

// EnterGroupByKey is called when production groupByClause is entered.
func (l *listener) EnterGroupByKey(ctx *grammar.GroupByKeyContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitGroupByKey(ctx)
	}
}

// EnterSortField is called when production sortField is entered.
func (l *listener) EnterSortField(ctx *grammar.SortFieldContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitSortField(ctx)
	}
}

// ExitSortField is called when production sortField is exited.
func (l *listener) ExitSortField(ctx *grammar.SortFieldContext) {
	if l.queryStmt != nil {
		l.queryStmt.completeSortField(ctx)
	}
}

// EnterHavingClause is called when production havingClause is entered.
func (l *listener) EnterHavingClause(ctx *grammar.HavingClauseContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitHaving(ctx)
	}
}

// ExitHavingClause is called when production havingClause is exited.
func (l *listener) ExitHavingClause(ctx *grammar.HavingClauseContext) {
	if l.queryStmt != nil {
		l.queryStmt.completeHaving(ctx)
	}
}

// EnterBoolExprAtom is called when production boolExprAtom is entered.
func (l *listener) EnterBoolExprAtom(ctx *grammar.BoolExprAtomContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitBoolExprAtom(ctx)
	}
}

// EnterBoolExpr is called when production boolExpr is entered.
func (l *listener) EnterBoolExpr(ctx *grammar.BoolExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.visitBoolExpr(ctx)
	}
}

// ExitBoolExpr is called when production boolExpr is exited.
func (l *listener) ExitBoolExpr(ctx *grammar.BoolExprContext) {
	if l.queryStmt != nil {
		l.queryStmt.completeBoolExpr(ctx)
	}
}

// statement returns query statement, if failure return error
func (l *listener) statement() (stmt.Statement, error) {
	switch {
	case l.useStmt != nil:
		return l.useStmt.build()
	case l.metadataStmt != nil:
		return l.metadataStmt.build()
	case l.storageStmt != nil:
		return l.storageStmt.build()
	case l.schemasStmt != nil:
		return l.schemasStmt.build()
	case l.queryStmt != nil:
		return l.queryStmt.build()
	case l.metricMetadataStmt != nil:
		return l.metricMetadataStmt.build()
	case l.stateStmt != nil:
		return l.stateStmt.build()
	case l.requestStmt != nil:
		return l.requestStmt.build()
	case l.brokerStmt != nil:
		return l.brokerStmt.build()
	case l.limitStmt != nil:
		return l.limitStmt.build()
	default:
		return nil, nil
	}
}
