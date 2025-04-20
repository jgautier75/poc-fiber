// Code generated from Filter.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Filter

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type FilterParser struct {
	*antlr.BaseParser
}

var FilterParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func filterParserInit() {
	staticData := &FilterParserStaticData
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
		"filter", "expr",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 21, 35, 2, 0, 7, 0, 2, 1, 7, 1, 1, 0, 4, 0, 6, 8, 0, 11, 0, 12, 0,
		7, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 3, 1, 22, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 30, 8, 1,
		10, 1, 12, 1, 33, 9, 1, 1, 1, 0, 1, 2, 2, 0, 2, 0, 0, 37, 0, 5, 1, 0, 0,
		0, 2, 21, 1, 0, 0, 0, 4, 6, 3, 2, 1, 0, 5, 4, 1, 0, 0, 0, 6, 7, 1, 0, 0,
		0, 7, 5, 1, 0, 0, 0, 7, 8, 1, 0, 0, 0, 8, 9, 1, 0, 0, 0, 9, 10, 5, 0, 0,
		1, 10, 1, 1, 0, 0, 0, 11, 12, 6, 1, -1, 0, 12, 13, 5, 1, 0, 0, 13, 14,
		3, 2, 1, 0, 14, 15, 5, 2, 0, 0, 15, 22, 1, 0, 0, 0, 16, 17, 5, 5, 0, 0,
		17, 22, 3, 2, 1, 4, 18, 19, 5, 17, 0, 0, 19, 20, 5, 6, 0, 0, 20, 22, 5,
		14, 0, 0, 21, 11, 1, 0, 0, 0, 21, 16, 1, 0, 0, 0, 21, 18, 1, 0, 0, 0, 22,
		31, 1, 0, 0, 0, 23, 24, 10, 2, 0, 0, 24, 25, 5, 4, 0, 0, 25, 30, 3, 2,
		1, 3, 26, 27, 10, 1, 0, 0, 27, 28, 5, 3, 0, 0, 28, 30, 3, 2, 1, 2, 29,
		23, 1, 0, 0, 0, 29, 26, 1, 0, 0, 0, 30, 33, 1, 0, 0, 0, 31, 29, 1, 0, 0,
		0, 31, 32, 1, 0, 0, 0, 32, 3, 1, 0, 0, 0, 33, 31, 1, 0, 0, 0, 4, 7, 21,
		29, 31,
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

// FilterParserInit initializes any static state used to implement FilterParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewFilterParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func FilterParserInit() {
	staticData := &FilterParserStaticData
	staticData.once.Do(filterParserInit)
}

// NewFilterParser produces a new parser instance for the optional input antlr.TokenStream.
func NewFilterParser(input antlr.TokenStream) *FilterParser {
	FilterParserInit()
	this := new(FilterParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &FilterParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Filter.g4"

	return this
}

// FilterParser tokens.
const (
	FilterParserEOF        = antlr.TokenEOF
	FilterParserOPAR       = 1
	FilterParserCPAR       = 2
	FilterParserOR         = 3
	FilterParserAND        = 4
	FilterParserNOT        = 5
	FilterParserCOMPARISON = 6
	FilterParserGT         = 7
	FilterParserGE         = 8
	FilterParserLT         = 9
	FilterParserLE         = 10
	FilterParserEQ         = 11
	FilterParserNE         = 12
	FilterParserLK         = 13
	FilterParserVALUE      = 14
	FilterParserTRUE       = 15
	FilterParserFALSE      = 16
	FilterParserPROPERTY   = 17
	FilterParserSTRING     = 18
	FilterParserINT        = 19
	FilterParserFLOAT      = 20
	FilterParserSPACE      = 21
)

// FilterParser rules.
const (
	FilterParserRULE_filter = 0
	FilterParserRULE_expr   = 1
)

// IFilterContext is an interface to support dynamic dispatch.
type IFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext

	// IsFilterContext differentiates from other interfaces.
	IsFilterContext()
}

type FilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterContext() *FilterContext {
	var p = new(FilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FilterParserRULE_filter
	return p
}

func InitEmptyFilterContext(p *FilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FilterParserRULE_filter
}

func (*FilterContext) IsFilterContext() {}

func NewFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterContext {
	var p = new(FilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FilterParserRULE_filter

	return p
}

func (s *FilterContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterContext) EOF() antlr.TerminalNode {
	return s.GetToken(FilterParserEOF, 0)
}

func (s *FilterContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *FilterContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FilterListener); ok {
		listenerT.EnterFilter(s)
	}
}

func (s *FilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FilterListener); ok {
		listenerT.ExitFilter(s)
	}
}

func (s *FilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FilterVisitor:
		return t.VisitFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FilterParser) Filter() (localctx IFilterContext) {
	localctx = NewFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, FilterParserRULE_filter)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(5)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&131106) != 0) {
		{
			p.SetState(4)
			p.expr(0)
		}

		p.SetState(7)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(9)
		p.Match(FilterParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPAR() antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	CPAR() antlr.TerminalNode
	NOT() antlr.TerminalNode
	PROPERTY() antlr.TerminalNode
	COMPARISON() antlr.TerminalNode
	VALUE() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FilterParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FilterParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FilterParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) OPAR() antlr.TerminalNode {
	return s.GetToken(FilterParserOPAR, 0)
}

func (s *ExprContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprContext) CPAR() antlr.TerminalNode {
	return s.GetToken(FilterParserCPAR, 0)
}

func (s *ExprContext) NOT() antlr.TerminalNode {
	return s.GetToken(FilterParserNOT, 0)
}

func (s *ExprContext) PROPERTY() antlr.TerminalNode {
	return s.GetToken(FilterParserPROPERTY, 0)
}

func (s *ExprContext) COMPARISON() antlr.TerminalNode {
	return s.GetToken(FilterParserCOMPARISON, 0)
}

func (s *ExprContext) VALUE() antlr.TerminalNode {
	return s.GetToken(FilterParserVALUE, 0)
}

func (s *ExprContext) AND() antlr.TerminalNode {
	return s.GetToken(FilterParserAND, 0)
}

func (s *ExprContext) OR() antlr.TerminalNode {
	return s.GetToken(FilterParserOR, 0)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FilterListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FilterListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FilterVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FilterParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *FilterParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, FilterParserRULE_expr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FilterParserOPAR:
		{
			p.SetState(12)
			p.Match(FilterParserOPAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(13)
			p.expr(0)
		}
		{
			p.SetState(14)
			p.Match(FilterParserCPAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FilterParserNOT:
		{
			p.SetState(16)
			p.Match(FilterParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(17)
			p.expr(4)
		}

	case FilterParserPROPERTY:
		{
			p.SetState(18)
			p.Match(FilterParserPROPERTY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(19)
			p.Match(FilterParserCOMPARISON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(20)
			p.Match(FilterParserVALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(31)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(29)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FilterParserRULE_expr)
				p.SetState(23)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(24)
					p.Match(FilterParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(25)
					p.expr(3)
				}

			case 2:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, FilterParserRULE_expr)
				p.SetState(26)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(27)
					p.Match(FilterParserOR)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(28)
					p.expr(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(33)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *FilterParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FilterParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
