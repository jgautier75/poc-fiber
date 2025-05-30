// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type FilterLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var FilterLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func filterlexerLexerInit() {
	staticData := &FilterLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'('", "')'", "'or'", "'and'", "'not'", "", "'gt'", "'ge'", "'lt'",
		"'le'", "'eq'", "'ne'", "'lk'", "", "'true'", "'false'",
	}
	staticData.SymbolicNames = []string{
		"", "OPAR", "CPAR", "OR", "AND", "NOT", "COMPARISON", "GT", "GE", "LT",
		"LE", "EQ", "NE", "LK", "VALUE", "TRUE", "FALSE", "PROPERTY", "STRING",
		"INT", "FLOAT", "SPACE",
	}
	staticData.RuleNames = []string{
		"OPAR", "CPAR", "OR", "AND", "NOT", "COMPARISON", "GT", "GE", "LT",
		"LE", "EQ", "NE", "LK", "VALUE", "TRUE", "FALSE", "PROPERTY", "ALLOWED_CHARACTERS",
		"STRING", "INT", "FLOAT", "DIGIT", "SPACE",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 21, 170, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 5, 1, 5, 1, 5, 3, 5, 70, 8, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7,
		1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 98,
		8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 1, 15, 1, 16, 4, 16, 112, 8, 16, 11, 16, 12, 16, 113, 1, 16, 1, 16,
		4, 16, 118, 8, 16, 11, 16, 12, 16, 119, 5, 16, 122, 8, 16, 10, 16, 12,
		16, 125, 9, 16, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 18, 5, 18, 133, 8,
		18, 10, 18, 12, 18, 136, 9, 18, 1, 18, 1, 18, 1, 19, 4, 19, 141, 8, 19,
		11, 19, 12, 19, 142, 1, 20, 4, 20, 146, 8, 20, 11, 20, 12, 20, 147, 1,
		20, 1, 20, 5, 20, 152, 8, 20, 10, 20, 12, 20, 155, 9, 20, 1, 20, 1, 20,
		4, 20, 159, 8, 20, 11, 20, 12, 20, 160, 3, 20, 163, 8, 20, 1, 21, 1, 21,
		1, 22, 1, 22, 1, 22, 1, 22, 0, 0, 23, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11,
		6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15,
		31, 16, 33, 17, 35, 0, 37, 18, 39, 19, 41, 20, 43, 0, 45, 21, 1, 0, 4,
		5, 0, 45, 45, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 10, 10, 13, 13, 39,
		39, 1, 0, 48, 57, 3, 0, 9, 10, 13, 13, 32, 32, 187, 0, 1, 1, 0, 0, 0, 0,
		3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0,
		11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0,
		0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0,
		0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0,
		0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 45, 1,
		0, 0, 0, 1, 47, 1, 0, 0, 0, 3, 49, 1, 0, 0, 0, 5, 51, 1, 0, 0, 0, 7, 54,
		1, 0, 0, 0, 9, 58, 1, 0, 0, 0, 11, 69, 1, 0, 0, 0, 13, 71, 1, 0, 0, 0,
		15, 74, 1, 0, 0, 0, 17, 77, 1, 0, 0, 0, 19, 80, 1, 0, 0, 0, 21, 83, 1,
		0, 0, 0, 23, 86, 1, 0, 0, 0, 25, 89, 1, 0, 0, 0, 27, 97, 1, 0, 0, 0, 29,
		99, 1, 0, 0, 0, 31, 104, 1, 0, 0, 0, 33, 111, 1, 0, 0, 0, 35, 126, 1, 0,
		0, 0, 37, 128, 1, 0, 0, 0, 39, 140, 1, 0, 0, 0, 41, 162, 1, 0, 0, 0, 43,
		164, 1, 0, 0, 0, 45, 166, 1, 0, 0, 0, 47, 48, 5, 40, 0, 0, 48, 2, 1, 0,
		0, 0, 49, 50, 5, 41, 0, 0, 50, 4, 1, 0, 0, 0, 51, 52, 5, 111, 0, 0, 52,
		53, 5, 114, 0, 0, 53, 6, 1, 0, 0, 0, 54, 55, 5, 97, 0, 0, 55, 56, 5, 110,
		0, 0, 56, 57, 5, 100, 0, 0, 57, 8, 1, 0, 0, 0, 58, 59, 5, 110, 0, 0, 59,
		60, 5, 111, 0, 0, 60, 61, 5, 116, 0, 0, 61, 10, 1, 0, 0, 0, 62, 70, 3,
		13, 6, 0, 63, 70, 3, 15, 7, 0, 64, 70, 3, 17, 8, 0, 65, 70, 3, 19, 9, 0,
		66, 70, 3, 21, 10, 0, 67, 70, 3, 23, 11, 0, 68, 70, 3, 25, 12, 0, 69, 62,
		1, 0, 0, 0, 69, 63, 1, 0, 0, 0, 69, 64, 1, 0, 0, 0, 69, 65, 1, 0, 0, 0,
		69, 66, 1, 0, 0, 0, 69, 67, 1, 0, 0, 0, 69, 68, 1, 0, 0, 0, 70, 12, 1,
		0, 0, 0, 71, 72, 5, 103, 0, 0, 72, 73, 5, 116, 0, 0, 73, 14, 1, 0, 0, 0,
		74, 75, 5, 103, 0, 0, 75, 76, 5, 101, 0, 0, 76, 16, 1, 0, 0, 0, 77, 78,
		5, 108, 0, 0, 78, 79, 5, 116, 0, 0, 79, 18, 1, 0, 0, 0, 80, 81, 5, 108,
		0, 0, 81, 82, 5, 101, 0, 0, 82, 20, 1, 0, 0, 0, 83, 84, 5, 101, 0, 0, 84,
		85, 5, 113, 0, 0, 85, 22, 1, 0, 0, 0, 86, 87, 5, 110, 0, 0, 87, 88, 5,
		101, 0, 0, 88, 24, 1, 0, 0, 0, 89, 90, 5, 108, 0, 0, 90, 91, 5, 107, 0,
		0, 91, 26, 1, 0, 0, 0, 92, 98, 3, 29, 14, 0, 93, 98, 3, 31, 15, 0, 94,
		98, 3, 39, 19, 0, 95, 98, 3, 41, 20, 0, 96, 98, 3, 37, 18, 0, 97, 92, 1,
		0, 0, 0, 97, 93, 1, 0, 0, 0, 97, 94, 1, 0, 0, 0, 97, 95, 1, 0, 0, 0, 97,
		96, 1, 0, 0, 0, 98, 28, 1, 0, 0, 0, 99, 100, 5, 116, 0, 0, 100, 101, 5,
		114, 0, 0, 101, 102, 5, 117, 0, 0, 102, 103, 5, 101, 0, 0, 103, 30, 1,
		0, 0, 0, 104, 105, 5, 102, 0, 0, 105, 106, 5, 97, 0, 0, 106, 107, 5, 108,
		0, 0, 107, 108, 5, 115, 0, 0, 108, 109, 5, 101, 0, 0, 109, 32, 1, 0, 0,
		0, 110, 112, 3, 35, 17, 0, 111, 110, 1, 0, 0, 0, 112, 113, 1, 0, 0, 0,
		113, 111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114, 123, 1, 0, 0, 0, 115,
		117, 5, 46, 0, 0, 116, 118, 3, 35, 17, 0, 117, 116, 1, 0, 0, 0, 118, 119,
		1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 122, 1, 0,
		0, 0, 121, 115, 1, 0, 0, 0, 122, 125, 1, 0, 0, 0, 123, 121, 1, 0, 0, 0,
		123, 124, 1, 0, 0, 0, 124, 34, 1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 126, 127,
		7, 0, 0, 0, 127, 36, 1, 0, 0, 0, 128, 134, 5, 39, 0, 0, 129, 133, 8, 1,
		0, 0, 130, 131, 5, 39, 0, 0, 131, 133, 5, 39, 0, 0, 132, 129, 1, 0, 0,
		0, 132, 130, 1, 0, 0, 0, 133, 136, 1, 0, 0, 0, 134, 132, 1, 0, 0, 0, 134,
		135, 1, 0, 0, 0, 135, 137, 1, 0, 0, 0, 136, 134, 1, 0, 0, 0, 137, 138,
		5, 39, 0, 0, 138, 38, 1, 0, 0, 0, 139, 141, 3, 43, 21, 0, 140, 139, 1,
		0, 0, 0, 141, 142, 1, 0, 0, 0, 142, 140, 1, 0, 0, 0, 142, 143, 1, 0, 0,
		0, 143, 40, 1, 0, 0, 0, 144, 146, 3, 43, 21, 0, 145, 144, 1, 0, 0, 0, 146,
		147, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148, 149,
		1, 0, 0, 0, 149, 153, 5, 46, 0, 0, 150, 152, 3, 43, 21, 0, 151, 150, 1,
		0, 0, 0, 152, 155, 1, 0, 0, 0, 153, 151, 1, 0, 0, 0, 153, 154, 1, 0, 0,
		0, 154, 163, 1, 0, 0, 0, 155, 153, 1, 0, 0, 0, 156, 158, 5, 46, 0, 0, 157,
		159, 3, 43, 21, 0, 158, 157, 1, 0, 0, 0, 159, 160, 1, 0, 0, 0, 160, 158,
		1, 0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 163, 1, 0, 0, 0, 162, 145, 1, 0,
		0, 0, 162, 156, 1, 0, 0, 0, 163, 42, 1, 0, 0, 0, 164, 165, 7, 2, 0, 0,
		165, 44, 1, 0, 0, 0, 166, 167, 7, 3, 0, 0, 167, 168, 1, 0, 0, 0, 168, 169,
		6, 22, 0, 0, 169, 46, 1, 0, 0, 0, 13, 0, 69, 97, 113, 119, 123, 132, 134,
		142, 147, 153, 160, 162, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// FilterLexerInit initializes any static state used to implement FilterLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewFilterLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func FilterLexerInit() {
	staticData := &FilterLexerLexerStaticData
	staticData.once.Do(filterlexerLexerInit)
}

// NewFilterLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewFilterLexer(input antlr.CharStream) *FilterLexer {
	FilterLexerInit()
	l := new(FilterLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &FilterLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Filter.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// FilterLexer tokens.
const (
	FilterLexerOPAR       = 1
	FilterLexerCPAR       = 2
	FilterLexerOR         = 3
	FilterLexerAND        = 4
	FilterLexerNOT        = 5
	FilterLexerCOMPARISON = 6
	FilterLexerGT         = 7
	FilterLexerGE         = 8
	FilterLexerLT         = 9
	FilterLexerLE         = 10
	FilterLexerEQ         = 11
	FilterLexerNE         = 12
	FilterLexerLK         = 13
	FilterLexerVALUE      = 14
	FilterLexerTRUE       = 15
	FilterLexerFALSE      = 16
	FilterLexerPROPERTY   = 17
	FilterLexerSTRING     = 18
	FilterLexerINT        = 19
	FilterLexerFLOAT      = 20
	FilterLexerSPACE      = 21
)
