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
	"strconv"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// baseStmtParser represents metadata statement parser
type baseStmtParser struct {
	namespace  string
	metricName string

	exprStack *collections.Stack
	condition stmt.Expr

	limit int

	err error
}

// visitLimit visits when production limit expression is entered
func (b *baseStmtParser) visitLimit(ctx *grammar.LimitClauseContext) {
	if ctx.L_INT() == nil {
		return
	}
	limit, err := strconv.ParseInt(ctx.L_INT().GetText(), 10, 32)
	if err != nil {
		b.err = err
		return
	}
	b.limit = int(limit)
}

// visitMetricName visits when production metricName expression is entered
func (b *baseStmtParser) visitMetricName(ctx *grammar.MetricNameContext) {
	b.metricName = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitPrefix visits when production namespace expression is entered
func (b *baseStmtParser) visitNamespace(ctx *grammar.NamespaceContext) {
	b.namespace = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitExpression visits when production expression is entered.
func (b *baseStmtParser) visitExpression(ctx *grammar.ExpressionContext) {
	var expr stmt.Expr
	switch {
	case ctx.Name() != nil:
		expr = b.createExpr(ctx.Name(), ctx)
	case ctx.T_OPEN_P() != nil:
		expr = &stmt.ParenExpr{}
	case ctx.T_AND() != nil:
		expr = &stmt.BinaryExpr{Operator: stmt.AND}
	case ctx.T_OR() != nil:
		expr = &stmt.BinaryExpr{Operator: stmt.OR}
	}

	b.exprStack.Push(expr)
}

// visitTagValue visits when production tag value expression is entered
func (b *baseStmtParser) visitTagValue(ctx *grammar.ValueContext) {
	if b.exprStack.Empty() {
		return
	}
	filterExpr := b.exprStack.Peek()
	value := strutil.GetStringValue(ctx.STRING().GetText())
	switch expr := filterExpr.(type) {
	case *stmt.NotExpr:
		b.setExprValue(expr.Expr, value)
	case stmt.Expr:
		b.setExprValue(expr, value)
	}
}

// setExprValue sets tag value for tag filter expression
func (b *baseStmtParser) setExprValue(expr stmt.Expr, value string) {
	switch e := expr.(type) {
	case *stmt.EqualsExpr:
		e.Value = value
	case *stmt.LikeExpr:
		e.Value = value
	case *stmt.RegexExpr:
		e.Regexp = value
	case *stmt.InExpr:
		e.Values = append(e.Values, value)
	}
}

// createExpr creates filer expr like equals, like, in and regex etc.
func (b *baseStmtParser) createExpr(name grammar.INameContext,
	ctx *grammar.ExpressionContext) stmt.Expr {
	var expr stmt.Expr
	if nameCtx, ok := name.(*grammar.NameContext); ok {
		nameStr := strutil.GetStringValue(nameCtx.Ident().GetText())
		switch {
		case ctx.T_EQUAL() != nil:
			expr = &stmt.EqualsExpr{Key: nameStr}
		case ctx.T_LIKE() != nil:
			if ctx.T_NOT() != nil {
				expr = &stmt.NotExpr{Expr: &stmt.LikeExpr{Key: nameStr}}
			} else {
				expr = &stmt.LikeExpr{Key: nameStr}
			}
		case ctx.T_REGEXP() != nil:
			expr = &stmt.RegexExpr{Key: nameStr}
		case ctx.T_NEQREGEXP() != nil:
			expr = &stmt.NotExpr{Expr: &stmt.RegexExpr{Key: nameStr}}
		case ctx.T_NOTEQUAL() != nil || ctx.T_NOTEQUAL2() != nil:
			expr = &stmt.NotExpr{Expr: &stmt.EqualsExpr{Key: nameStr}}
		case ctx.T_IN() != nil:
			if ctx.T_NOT() != nil {
				expr = &stmt.NotExpr{Expr: &stmt.InExpr{Key: nameStr}}
			} else {
				expr = &stmt.InExpr{Key: nameStr}
			}
		}
	}
	return expr
}

// completeExpression completes a filter expression for query condition
func (b *baseStmtParser) completeExpression() {
	expr := b.exprStack.Pop()
	e, ok := expr.(stmt.Expr)
	if !ok {
		return
	}
	if !b.exprStack.Empty() {
		parent := b.exprStack.Peek()
		switch parentExpr := parent.(type) {
		case *stmt.BinaryExpr:
			if parentExpr.Left == nil {
				parentExpr.Left = e
			} else if parentExpr.Right == nil {
				parentExpr.Right = e
			}
		case *stmt.ParenExpr:
			parentExpr.Expr = e
		}
	}
	b.condition = e
}

// setExprParam sets expr's param(call,paren,binary)
func (b *baseStmtParser) setExprParam(param stmt.Expr) {
	if b.exprStack.Empty() {
		return
	}

	switch expr := b.exprStack.Peek().(type) {
	case *stmt.CallExpr:
		expr.Params = append(expr.Params, param)
	case *stmt.ParenExpr:
		expr.Expr = param
	case *stmt.BinaryExpr:
		if expr.Left == nil {
			expr.Left = param
		} else if expr.Right == nil {
			expr.Right = param
		}
	default:
	}
}
