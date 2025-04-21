package parser

import (
	"fmt"
	"testing"
)

func TestFilterSearch(t *testing.T) {
	expressions, errorNodes, _ := FromInputString("lname eq 'hopper'")
	for _, errnode := range errorNodes {
		fmt.Printf("error node [%s]", errnode.GetText())
	}
	for _, expr := range expressions {
		fmt.Printf("type [%s] value [%s]", expr.Type, expr.TextValue)
	}
}
