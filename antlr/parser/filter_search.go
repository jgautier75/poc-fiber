package parser

import (
	"poc-fiber/commons"

	"github.com/antlr4-go/antlr/v4"
)

func FromInputString(inStr string) (expressions []SearchExpression, errorNodes []antlr.ErrorNode, listenerErrors CustomErrorListener) {

	if inStr == "" {
		var nilListener CustomErrorListener
		searchExpressions := make([]SearchExpression, 0)
		return searchExpressions, nil, nilListener
	}

	is := antlr.NewInputStream(inStr)
	lexer := NewFilterLexer(is)
	lexer.RemoveErrorListeners()
	errListener := NewCustomErrorListener()
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	fParser := NewFilterParser(tokenStream)
	fParser.AddErrorListener(errListener)
	listener := SearchListener{}
	antlr.NewParseTreeWalker().Walk(&listener, fParser.Filter())
	return listener.Expressions, listener.ErrorNodes, *errListener
}

func ConvertErrorNodes(status int, errorNodes []antlr.ErrorNode) commons.ApiError {
	apiErrorDetails := make([]commons.ApiErrorDetails, 1)
	for _, err := range errorNodes {
		detail := commons.ApiErrorDetails{
			Field:  err.GetSymbol().GetText(),
			Detail: err.GetText(),
		}
		apiErrorDetails = append(apiErrorDetails, detail)
	}
	apiError := commons.ApiError{
		Code:    status,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: commons.SearchFilter,
		Details: apiErrorDetails,
	}
	return apiError
}
