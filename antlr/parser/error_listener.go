package parser

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

// CustomErrorListener is a custom implementation of antlr.ErrorListener.
type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []string
}

// NewCustomErrorListener creates a new instance of CustomErrorListener.
func NewCustomErrorListener() *CustomErrorListener {
	return &CustomErrorListener{
		Errors: []string{},
	}
}

// SyntaxError is triggered when the parser encounters a syntax error.
func (l *CustomErrorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol interface{},
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	errorMessage := fmt.Sprintf("Syntax error at line %d:%d - %s", line, column, msg)
	l.Errors = append(l.Errors, errorMessage)
}

// HasErrors checks if any errors were recorded.
func (l *CustomErrorListener) HasErrors() bool {
	return len(l.Errors) > 0
}

// GetErrors returns the list of errors.
func (l *CustomErrorListener) GetErrors() []string {
	return l.Errors
}
