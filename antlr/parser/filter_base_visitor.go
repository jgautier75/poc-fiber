// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Filter

import "github.com/antlr4-go/antlr/v4"

type BaseFilterVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseFilterVisitor) VisitFilter(ctx *FilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFilterVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}
