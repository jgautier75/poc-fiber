package parser

import (
	"github.com/antlr4-go/antlr/v4"
)

func FromInputString(inStr string) ([]SearchExpression, []antlr.ErrorNode) {
	is := antlr.NewInputStream(inStr)
	lexer := NewFilterLexer(is)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	fParser := NewFilterParser(tokenStream)
	listener := SearchListener{}
	antlr.NewParseTreeWalker().Walk(&listener, fParser.Filter())
	return listener.Expressions, listener.ErrorNodes
}
