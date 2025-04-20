// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Filter

import "github.com/antlr4-go/antlr/v4"

// FilterListener is a complete listener for a parse tree produced by FilterParser.
type FilterListener interface {
	antlr.ParseTreeListener

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)
}
