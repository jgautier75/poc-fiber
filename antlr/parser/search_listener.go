package parser

import (
	"github.com/antlr4-go/antlr/v4"
)

const (
	Property           string = "PROPERTY"
	Comparison         string = "COMPARISON"
	Operator           string = "OPERATOR"
	Negation           string = "NEGATION"
	OpeningParenthesis string = "OPENING_PARENTHESIS"
	ClosingParenthesis string = "CLOSING_PARENTHESIS"
	Value              string = "VALUE"
)

type SearchExpression struct {
	Type      string
	TextValue string
}

// SearchListener implements the ParseTreeListener interface
type SearchListener struct {
	Expressions []SearchExpression
	ErrorNodes  []antlr.ErrorNode
}

// VisitTerminal is called when a terminal node (leaf node) is visited in the parse tree.
func (l *SearchListener) VisitTerminal(node antlr.TerminalNode) {
	tokenType := node.GetSymbol().GetTokenType()
	switch tokenType {
	case FilterLexerAND:
	case FilterLexerOR:
		l.Expressions = append(l.Expressions, SearchExpression{Type: Operator, TextValue: node.GetText()})
	case FilterLexerCOMPARISON:
		l.Expressions = append(l.Expressions, SearchExpression{Type: Comparison, TextValue: node.GetText()})
	case FilterLexerOPAR:
		l.Expressions = append(l.Expressions, SearchExpression{Type: OpeningParenthesis, TextValue: node.GetText()})
	case FilterLexerCPAR:
		l.Expressions = append(l.Expressions, SearchExpression{Type: ClosingParenthesis, TextValue: node.GetText()})
	case FilterLexerVALUE:
		l.Expressions = append(l.Expressions, SearchExpression{Type: Value, TextValue: node.GetText()})
	case FilterLexerPROPERTY:
		l.Expressions = append(l.Expressions, SearchExpression{Type: Property, TextValue: node.GetText()})
	case FilterLexerNOT:
		l.Expressions = append(l.Expressions, SearchExpression{Type: Negation, TextValue: node.GetText()})
	default:
	}
}

// VisitErrorNode is called when an error node is visited.
func (l *SearchListener) VisitErrorNode(node antlr.ErrorNode) {
	l.ErrorNodes = append(l.ErrorNodes, node)
}

// EnterEveryRule is called when entering any rule in the parse tree.
func (l *SearchListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// No need to implement
}

// ExitEveryRule is called when exiting any rule in the parse tree.
func (l *SearchListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	// No need to implement
}
