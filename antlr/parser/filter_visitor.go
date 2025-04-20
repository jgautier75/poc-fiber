// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Filter

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by FilterParser.
type FilterVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by FilterParser#filter.
	VisitFilter(ctx *FilterContext) interface{}

	// Visit a parse tree produced by FilterParser#expr.
	VisitExpr(ctx *ExprContext) interface{}
}
