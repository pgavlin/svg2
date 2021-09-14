package cssvalue

import (
	"errors"
	"io"
	"math/bits"
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokens: ident, number, literal, ||, &&

const (
	identToken = unicode.MaxRune + iota
	numberToken
	literalToken
	anyOfToken
	allOfToken
)

type lexer struct {
	r io.Reader

	buf [1]rune
	err error
	nb  int

	token rune
	value strings.Builder
}

func (l *lexer) readRune() (rune, error) {
	var buf [4]byte

	if _, err := l.r.Read(buf[:1]); err != nil {
		return 0, err
	}

	if !utf8.RuneStart(buf[0]) {
		return utf8.RuneError, nil
	}

	sz := bits.LeadingZeros8(^buf[0])
	if sz == 0 {
		return rune(buf[0]), nil
	}

	if _, err := l.r.Read(buf[1 : sz-1]); err != nil {
		return 0, err
	}

	c, _ := utf8.DecodeRune(buf[:sz])
	return c, nil
}

func (l *lexer) peek() rune {
	if l.nb != 1 {
		l.buf[0], l.err = l.readRune()
		l.nb = 1
	}
	return l.buf[0]
}

func (l *lexer) chomp() rune {
	if l.nb != 1 {
		panic("expected a buffered rune")
	}
	l.nb = 0
	return l.buf[0]
}

func (l *lexer) next() (rune, error) {
	if l.nb == 1 {
		l.nb = 0
		return l.buf[0], l.err
	}
	return l.readRune()
}

func startsIdent(c rune) bool {
	return c == '-' || c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (l *lexer) lexNext() (rune, error) {
	l.value.Reset()

	for {
		c := l.peek()
		switch c {
		case 0:
			return 0, l.err
		case '-':
			return l.lexIdentOrNumberOrToken()
		case '+':
			return l.lexNumberOrPlus()
		case '|':
			return l.lexOneOrAnyOf()
		case '&':
			return l.lexAllOfOrRune()
		case '<', '>', '[', ']', '{', '}', '\'', ',', '*', '?', '#', '!':
			return l.chomp(), nil
		default:
			switch {
			case c == utf8.RuneError || unicode.IsSpace(c):
				l.chomp()
			case startsIdent(c):
				return l.lexIdent(0)
			case c == '∞' || isDigit(c):
				return l.lexNumber(0)
			case c == '"':
				return l.lexLiteral()
			default:
				l.value.WriteRune(l.chomp())
				return literalToken, nil
			}
		}
	}
}

func (l *lexer) lex() error {
	t, err := l.lexNext()
	l.token = t
	return err
}

func (l *lexer) lexIdentOrNumberOrToken() (rune, error) {
	first, err := l.readRune()
	if err != nil {
		return 0, err
	}

	next := l.peek()
	if next >= '0' && next <= '9' || next == '∞' {
		return l.lexNumber(first)
	}
	if startsIdent(next) {
		return l.lexIdent(first)
	}
	return first, nil
}

func (l *lexer) lexNumberOrPlus() (rune, error) {
	first, next := l.chomp(), l.peek()
	if next >= '0' && next <= '9' || next == '∞' {
		return l.lexNumber(first)
	}
	return first, nil
}

func (l *lexer) lexOneOrAnyOf() (rune, error) {
	_, next := l.chomp(), l.peek()
	if next == '|' {
		l.chomp()
		return anyOfToken, nil
	}
	return '|', nil
}

func (l *lexer) lexAllOfOrRune() (rune, error) {
	_, next := l.chomp(), l.peek()
	if next == '&' {
		l.chomp()
		return allOfToken, nil
	}
	return '&', nil
}

func (l *lexer) lexIdent(first rune) (rune, error) {
	if first != 0 {
		l.value.WriteRune(first)
	}
	for startsIdent(l.peek()) {
		l.value.WriteRune(l.chomp())
	}
	return identToken, nil
}

func (l *lexer) lexNumber(first rune) (rune, error) {
	if first != 0 {
		l.value.WriteRune(first)
	}

	// ∞ or digits
	c := l.peek()
	if c == '∞' {
		l.value.WriteRune(l.chomp())
		return numberToken, nil
	}
	for isDigit(c) {
		l.value.WriteRune(l.chomp())
		c = l.peek()
	}

	// decimal
	if c == '.' {
		l.value.WriteRune(l.chomp())

		c = l.peek()
		for isDigit(c) {
			l.value.WriteRune(l.chomp())
			c = l.peek()
		}
	}
	if c != 'e' && c != 'E' {
		return numberToken, nil
	}

	// exponent
	l.value.WriteRune(l.chomp())
	c = l.peek()

	if c == '-' || c == '+' {
		l.value.WriteRune(l.chomp())
		c = l.peek()
	}

	if !isDigit(c) {
		return 0, errors.New("invalid number token")
	}
	for isDigit(c) {
		l.value.WriteRune(l.chomp())
		c = l.peek()
	}

	return numberToken, nil
}

func (l *lexer) lexLiteral() (rune, error) {
	l.chomp()

	for l.peek() != '"' {
		l.value.WriteRune(l.chomp())
	}
	l.chomp()

	return literalToken, nil
}
