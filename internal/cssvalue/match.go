package cssvalue

import (
	"fmt"
	"io"
	"math"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type Token struct {
	Type  css.TokenType
	Value string
}

type Function struct {
	Params []Term
	Type   *BasicType
}

type Context struct {
	Properties   map[string]Term
	NonTerminals map[string]Term
	Functions    map[string]*Function
}

type Capture struct {
	Name   string
	Values []Token
}

func Match(ctx *Context, term Term, r io.Reader) []Capture {
	matcher := matcher{
		context: ctx,
		lexer:   css.NewLexer(parse.NewInput(r)),
	}
	if !matcher.match(term) {
		return nil
	}
	matcher.captures = append(matcher.captures, Capture{
		Name:   "",
		Values: matcher.result,
	})
	return matcher.captures
}

type matcher struct {
	context *Context

	buffer []Token
	offset int
	lexer  *css.Lexer

	result   []Token
	captures []Capture
}

func (m *matcher) peek(n int) Token {
	for len(m.buffer[m.offset:]) < n {
		typ, value := m.lexer.Next()
		if typ == css.ErrorToken {
			return Token{Type: typ, Value: string(value)}
		}
		if typ != css.WhitespaceToken {
			m.buffer = append(m.buffer, Token{Type: typ, Value: string(value)})
		}
	}

	return m.buffer[m.offset+n-1]
}

func (m *matcher) chomp() {
	if len(m.buffer) == 0 {
		panic("buffer underflow")
	}
	m.result = append(m.result, m.buffer[m.offset])
	m.offset++
}

func (m *matcher) match(t Term) bool {
	so, sr, sc := m.offset, len(m.result), len(m.captures)

	matched := false
	switch t := t.(type) {
	case Keyword:
		matched = m.matchKeyword(t)
	case *BasicType:
		matched = m.matchBasicType(t)
	case PropertyType:
		matched = m.matchPropertyType(t)
	case NonTerminal:
		matched = m.matchNonTerminal(t)
	case Literal:
		matched = m.matchLiteral(t)
	case *Seq:
		matched = m.matchSeq(t)
	case *AllOf:
		matched = m.matchAllOf(t)
	case *AnyOf:
		matched = m.matchAnyOf(t)
	case *OneOf:
		matched = m.matchOneOf(t.Operands)
	case *Group:
		matched = m.matchGroup(t)
	case *ZeroOrMore:
		matched = m.matchZeroOrMore(t)
	case *OneOrMore:
		matched = m.matchOneOrMore(t)
	case *ZeroOrOne:
		matched = m.matchZeroOrOne(t)
	case *Repeat:
		matched = m.matchRepeat(t)
	default:
		panic(fmt.Errorf("unexpected term %T", t))
	}

	if !matched {
		m.offset, m.result, m.captures = so, m.result[:sr], m.captures[:sc]
	}
	return matched
}

func (m *matcher) matchKeyword(t Keyword) bool {
	next := m.peek(1)
	if next.Type != css.IdentToken || next.Value != string(t) {
		return false
	}
	m.chomp()
	return true
}

// <number [0,∞]> [ / <number [0,∞]> ]?
var ratioTerm = &Seq{
	Operands: []Term{
		&BasicType{
			Name: "number",
			Range: &Range{
				Min: 0,
				Max: math.Inf(1),
			},
		},
		&ZeroOrOne{
			Operand: &Group{
				Operand: &Seq{
					Operands: []Term{
						Literal("/"),
						&BasicType{
							Name: "number",
							Range: &Range{
								Min: 0,
								Max: math.Inf(1),
							},
						},
					},
				},
			},
		},
	},
}

func (m *matcher) matchBasicType(t *BasicType) bool {
	if m.peek(1).Type == css.FunctionToken {
		f, ok := m.matchFunction()
		if !ok {
			return false
		}
		return f.Type.Name == t.Name
	}

	switch t.Name {
	case "custom-ident", "dashed-ident":
		next := m.peek(1)
		if next.Type != css.IdentToken {
			return false
		}
		m.chomp()
		return true
	case "string":
		next := m.peek(1)
		if next.Type != css.StringToken {
			return false
		}
		m.chomp()
		return true
	case "url":
		next := m.peek(1)
		if next.Type != css.URLToken {
			return false
		}
		m.chomp()
		return true
	case "color":
		next := m.peek(1)
		if next.Type != css.HashToken && next.Type != css.IdentToken {
			return false
		}
		m.chomp()
		return true
	case "integer", "number":
		next := m.peek(1)
		if next.Type != css.NumberToken {
			return false
		}
		m.chomp()
		return true
	case "percentage":
		next := m.peek(1)
		if next.Type != css.PercentageToken {
			return false
		}
		m.chomp()
		return true
	case "dimension", "length", "angle", "time", "frequency", "resolution":
		next := m.peek(1)
		if next.Type != css.DimensionToken {
			return false
		}
		m.chomp()
		return true
	case "ratio":
		return m.match(ratioTerm)
	default:
		panic(fmt.Errorf("unexpected basic type %v", t.Name))
	}
}

func (m *matcher) matchPropertyType(t PropertyType) bool {
	if term, ok := m.context.Properties[string(t)]; ok {
		c := len(m.result)
		if !m.match(term) {
			return false
		}
		m.captures = append(m.captures, Capture{
			Name:   string(t),
			Values: m.result[c:],
		})
	}
	return false
}

func (m *matcher) matchNonTerminal(t NonTerminal) bool {
	if term, ok := m.context.NonTerminals[string(t)]; ok {
		c := len(m.result)
		if !m.match(term) {
			return false
		}
		m.captures = append(m.captures, Capture{
			Name:   string(t),
			Values: m.result[c:],
		})
		return true
	}
	return false
}

func (m *matcher) matchLiteral(t Literal) bool {
	expect := css.DelimToken
	switch t {
	case ":":
		expect = css.ColonToken
	case "{":
		expect = css.LeftBraceToken
	case "[":
		expect = css.LeftBracketToken
	case "(":
		expect = css.LeftParenthesisToken
	case "}":
		expect = css.RightBraceToken
	case "]":
		expect = css.RightBracketToken
	case ")":
		expect = css.RightParenthesisToken
	case ";":
		expect = css.SemicolonToken
	}

	next := m.peek(1)
	if next.Type != expect || next.Value != string(t) {
		return false
	}
	m.chomp()
	return true
}

func (m *matcher) matchSeq(t *Seq) bool {
	c := len(m.result)
	for _, t := range t.Operands {
		if !m.match(t) {
			return false
		}
	}
	if t.Name != "" {
		m.captures = append(m.captures, Capture{
			Values: m.result[c:],
			Name:   t.Name,
		})
	}
	return true
}

func (m *matcher) matchAllOf(t *AllOf) bool {
	pending := append([]Term{}, t.Operands...)

	for len(pending) != 0 {
		if m.peek(1).Type == css.ErrorToken {
			return false
		}

		any := false
		for i := range pending {
			if m.match(pending[i]) {
				pending, any = append(pending[:i], pending[i+1:]...), true
				break
			}
		}
		if !any {
			return false
		}
	}

	return true
}

func (m *matcher) matchAnyOf(t *AnyOf) bool {
	if !m.matchOneOf(t.Operands) {
		return false
	}
	for m.matchOneOf(t.Operands) {
	}
	return true
}

func (m *matcher) matchOneOf(operands []Term) bool {
	for _, t := range operands {
		if m.match(t) {
			return true
		}
	}
	return false
}

func (m *matcher) matchGroup(t *Group) bool {
	// TODO: required group?
	c := len(m.result)
	if !m.match(t.Operand) {
		return false
	}
	if t.Name != "" {
		m.captures = append(m.captures, Capture{
			Values: m.result[c:],
			Name:   t.Name,
		})
	}
	return true
}

func (m *matcher) matchZeroOrMore(t *ZeroOrMore) bool {
	for m.match(t.Operand) {
	}
	return true
}

func (m *matcher) matchOneOrMore(t *OneOrMore) bool {
	if !m.match(t.Operand) {
		return false
	}
	for m.match(t.Operand) {
	}
	return true
}

func (m *matcher) matchZeroOrOne(t *ZeroOrOne) bool {
	m.match(t.Operand)
	return true
}

func (m *matcher) matchRepeat(t *Repeat) bool {
	n := 0
	for n < t.Max {
		next := m.peek(1)
		if next.Type == css.ErrorToken {
			break
		}
		if n > 0 && t.Commas {
			if next.Type != css.CommaToken {
				return false
			}
			m.chomp()
		}
		if !m.match(t.Operand) {
			break
		}
		n++
	}
	return n >= t.Min
}

func (m *matcher) matchFunction() (*Function, bool) {
	s := m.offset

	token := m.peek(1)
	fn, ok := m.context.Functions[token.Value[:len(token.Value)-1]]
	if !ok {
		return nil, false
	}

	// Consume the function token.
	m.chomp()

	for i, t := range fn.Params {
		if !m.match(t) {
			m.offset = s
			return nil, false
		}

		// Note that this is a bit more permissive than the actual spec, as it does not
		// require that either _all_ or _no_ operands are separated with commas.
		if i < len(fn.Params)-1 && m.peek(1).Type == css.CommaToken {
			m.chomp()
		}
	}

	if m.peek(1).Type != css.RightParenthesisToken {
		m.offset = s
		return nil, false
	}

	// Consume the right parenthesis
	m.chomp()

	return fn, true
}
