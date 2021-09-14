package cssvalue

import (
	"errors"
	"io"
	"math"
	"strconv"
)

// oneOf -> anyOf
//       |  oneOf "|" anyOf
//
// anyOf -> allOf
//       |  anyOf "||" allOf
//
// allOf -> seq
//       |  allOf "&&" seq
//
// seq -> multipliable
//     |  "?" ident multipliable
//     |  seq multipliable
//
// multiplier -> multipliable
//            |  multipliable "*"
//            |  multipliable "+"
//            |  multipliable "{" number "}"
//            |  multipliable "{" number "," number "}"
//            |  group "!"
//
// multipliable -> group
//              |  type
//
// group -> "[" oneOf "]"
//       |  "[" "?" ident oneOf "]"
//
// type -> "<" ident ">"
//      |  "<" ident "[" number "," number "]" ">"
//      |  "<" "'" ident "'" ">"
//      | ident
//      | literal

func ParseGrammar(r io.Reader) (Term, error) {
	lexer := lexer{r: r}
	if err := lexer.lex(); err != nil {
		return nil, err
	}
	term, err := parseOneOf(&lexer)
	if err != nil {
		return nil, err
	}
	if lexer.token != 0 {
		return nil, errors.New("expected EOF")
	}
	return term, nil
}

func parseOneOf(l *lexer) (Term, error) {
	var operands []Term
	for {
		operand, err := parseAnyOf(l)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)

		if l.token != '|' {
			break
		}
		if err := l.lex(); err != nil {
			return nil, err
		}
	}
	if len(operands) == 1 {
		return operands[0], nil
	}
	return &OneOf{Operands: operands}, nil
}

func parseAnyOf(l *lexer) (Term, error) {
	var operands []Term
	for {
		operand, err := parseAllOf(l)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)

		if l.token != anyOfToken {
			break
		}
		if err := l.lex(); err != nil {
			return nil, err
		}
	}
	if len(operands) == 1 {
		return operands[0], nil
	}
	return &AnyOf{Operands: operands}, nil
}

func parseAllOf(l *lexer) (Term, error) {
	var operands []Term
	for {
		operand, err := parseSeq(l)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)

		if l.token != allOfToken {
			break
		}
		if err := l.lex(); err != nil {
			return nil, err
		}
	}
	if len(operands) == 1 {
		return operands[0], nil
	}
	return &AllOf{Operands: operands}, nil
}

func parseSeq(l *lexer) (Term, error) {
	name := ""
	if l.token == '?' {
		if err := l.lex(); err != nil {
			return nil, err
		}
		if l.token != identToken {
			return nil, errors.New("expected ident")
		}
		name = l.value.String()
		if err := l.lex(); err != nil {
			return nil, err
		}
	}

	var operands []Term
	for {
		operand, err := parseMultiplier(l)
		if err != nil {
			return nil, err
		}
		if operand == nil {
			break
		}
		operands = append(operands, operand)

		if l.token == ']' || l.token == '|' || l.token == anyOfToken || l.token == allOfToken {
			break
		}
	}
	if len(operands) == 1 && name == "" {
		return operands[0], nil
	}
	return &Seq{Name: name, Operands: operands}, nil
}

func parseMultiplier(l *lexer) (Term, error) {
	operand, err := parseMultipliable(l)
	if err != nil || operand == nil {
		return nil, err
	}

	switch l.token {
	case '*':
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return &ZeroOrMore{Operand: operand}, nil
	case '+':
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return &OneOrMore{Operand: operand}, nil
	case '?':
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return &ZeroOrOne{Operand: operand}, nil
	case '!':
		group, ok := operand.(*Group)
		if !ok {
			return nil, errors.New("expected a group")
		}
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		group.Required = true
		return group, nil
	case '#', '{':
		// handled below
	default:
		return operand, nil
	}

	commas := false
	if l.token == '#' {
		commas = true
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}

		if l.token != '{' {
			return &Repeat{Min: 1, Max: int(math.MaxInt64), Commas: true, Operand: operand}, nil
		}
	}

	if err := l.lex(); err != nil {
		return nil, err
	}

	if l.token != numberToken {
		return nil, errors.New("expected a number")
	}
	min, err := strconv.ParseInt(l.value.String(), 10, 0)
	if err != nil {
		return nil, err
	}
	if err := l.lex(); err != nil {
		return nil, err
	}

	max := min
	if l.token == ',' {
		if err := l.lex(); err != nil {
			return nil, err
		}
		max = math.MaxInt64
		if l.token == numberToken {
			max, err = strconv.ParseInt(l.value.String(), 10, 0)
			if err != nil {
				return nil, err
			}
			if err := l.lex(); err != nil {
				return nil, err
			}
		}
	}

	if l.token != '}' {
		return nil, errors.New("expected '}'")
	}
	if err := l.lex(); err != nil && err != io.EOF {
		return nil, err
	}

	return &Repeat{Min: int(min), Max: int(max), Commas: commas, Operand: operand}, nil
}

func parseMultipliable(l *lexer) (Term, error) {
	switch l.token {
	case '[':
		return parseGroup(l)
	case '<':
		return parseType(l)
	case identToken:
		v := l.value.String()
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return Keyword(v), nil
	case literalToken:
		v := l.value.String()
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return Literal(v), nil
	case ',':
		// Outside the context of a bracketed range or repeat clause, commas should be
		// treated as literals.
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}
		return Literal(string([]rune{l.token})), nil
	case 0:
		return nil, nil
	default:
		return nil, errors.New("expected a group, type, keyword, or literal")
	}
}

func parseGroup(l *lexer) (*Group, error) {
	// skip '['
	if err := l.lex(); err != nil {
		return nil, err
	}

	name := ""
	if l.token == '?' {
		if err := l.lex(); err != nil {
			return nil, err
		}
		if l.token != identToken {
			return nil, errors.New("expected ident")
		}
		name = l.value.String()
		if err := l.lex(); err != nil {
			return nil, err
		}
	}

	operand, err := parseOneOf(l)
	if err != nil {
		return nil, err
	}
	if l.token != ']' {
		return nil, errors.New("expected ']'")
	}
	if err := l.lex(); err != nil && err != io.EOF {
		return nil, err
	}

	return &Group{Name: name, Operand: operand}, nil
}

func parseType(l *lexer) (Term, error) {
	// skip '<'
	if err := l.lex(); err != nil {
		return nil, err
	}

	switch l.token {
	case identToken:
		name := l.value.String()
		if err := l.lex(); err != nil {
			return nil, err
		}

		isBasicType, allowsRange := false, false
		switch name {
		case "custom-ident", "dashed-ident", "string", "url", "color":
			isBasicType = true
		case "integer", "number", "dimension", "percentage", "ratio",
			"length", "angle", "time", "frequency", "resolution":
			isBasicType, allowsRange = true, true
		}

		var r *Range
		if allowsRange {
			if l.token == '[' {
				rng, err := parseRange(l)
				if err != nil {
					return nil, err
				}
				r = rng
			}
		}

		if l.token != '>' {
			return nil, errors.New("expected '>'")
		}
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}

		if isBasicType {
			return &BasicType{Name: name, Range: r}, nil
		}
		return NonTerminal(name), nil
	case '\'':
		if err := l.lex(); err != nil {
			return nil, err
		}

		if l.token != identToken {
			return nil, errors.New("expected an identifier")
		}
		name := l.value.String()
		if err := l.lex(); err != nil {
			return nil, err
		}

		if l.token != '\'' {
			return nil, errors.New("expected '")
		}
		if err := l.lex(); err != nil {
			return nil, err
		}

		if l.token != '>' {
			return nil, errors.New("expected '>'")
		}
		if err := l.lex(); err != nil && err != io.EOF {
			return nil, err
		}

		return PropertyType(name), nil
	default:
		return nil, errors.New("expected an identifier or '")
	}
}

func parseRange(l *lexer) (*Range, error) {
	// skip '['
	if err := l.lex(); err != nil {
		return nil, err
	}

	if l.token != numberToken {
		return nil, errors.New("expected a number")
	}
	min, err := strconv.ParseFloat(l.value.String(), 64)
	if err != nil {
		return nil, err
	}
	if err := l.lex(); err != nil {
		return nil, err
	}

	if l.token != ',' {
		return nil, errors.New("expected ','")
	}
	if err := l.lex(); err != nil {
		return nil, err
	}

	if l.token != numberToken {
		return nil, errors.New("expected a number")
	}
	max, err := strconv.ParseFloat(l.value.String(), 64)
	if err != nil {
		return nil, err
	}
	if err := l.lex(); err != nil {
		return nil, err
	}

	if l.token != ']' {
		return nil, errors.New("expected ']'")
	}
	if err := l.lex(); err != nil {
		return nil, err
	}

	return &Range{Min: min, Max: max}, nil
}
