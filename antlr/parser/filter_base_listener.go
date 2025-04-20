// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Filter

import "github.com/antlr4-go/antlr/v4"

// BaseFilterListener is a complete listener for a parse tree produced by FilterParser.
type BaseFilterListener struct{}

var _ FilterListener = &BaseFilterListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseFilterListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseFilterListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseFilterListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseFilterListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseFilterListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseFilterListener) ExitFilter(ctx *FilterContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseFilterListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseFilterListener) ExitExpr(ctx *ExprContext) {}
