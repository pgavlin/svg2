package svg

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/tdewolff/parse/v2/css"
)

const (
	UserSpaceOnUse    Units = "userSpaceOnUse"
	ObjectBoundingBox Units = "objectBoundingBox"
)

type Ident string

func (i *Ident) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	if tokens[0].Type != css.IdentToken {
		return errors.New("expected an identifier")
	}

	*i = Ident(tokens[0].Value)
	return nil
}

type Units string

func (u *Units) UnmarshalText(text []byte) error {
	switch Units(text) {
	case UserSpaceOnUse, ObjectBoundingBox:
		*u = Units(text)
		return nil
	default:
		return errors.New("units must be one of \"userSpaceOnUse\" or \"objectBoundingBox\"")
	}
}

type URLIdent struct {
	URL   string
	Ident string
}

func (ui *URLIdent) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.URLToken:
		ui.URL = token.Value[len("url(") : len(token.Value)-1]
		return nil
	case css.IdentToken:
		ui.Ident = token.Value
		return nil
	default:
		return errors.New("expected a URL or identifier")
	}
}

// NumberPercentage represents an XML number or percentage.
type NumberPercentage struct {
	Number     float64
	Percentage float64
}

func (np *NumberPercentage) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.NumberToken:
		n, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return err
		}
		np.Number = n
		return nil
	case css.PercentageToken:
		p, err := strconv.ParseFloat(token.Value[:len(token.Value)-1], 64)
		if err != nil {
			return err
		}
		np.Percentage = p / 100
		return nil
	default:
		return errors.New("expected a number or percentage")
	}
}

// NumberIdent represents an XML number or identifier.
type NumberIdent struct {
	Number float64
	Ident  string
}

func (ni *NumberIdent) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.NumberToken:
		n, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return err
		}
		ni.Number = n
		return nil
	case css.IdentToken:
		ni.Ident = token.Value
		return nil
	default:
		return errors.New("expected a number or identifier")
	}
}

// Length represents a CSS length.
type Length struct {
	Value float64
	Units string
}

func parseLength(token cssToken) (Length, error) {
	switch token.Type {
	case css.NumberToken:
		if token.Value == "0" {
			return Length{}, nil
		}
	case css.DimensionToken:
		// OK
	default:
		return Length{}, errors.New("expceted a dimension")
	}

	// Snip off the units
	v, units := token.Value, ""
	for i := len(v) - 1; i >= 0; i-- {
		c := v[i]
		if c >= '0' && c <= '9' {
			v, units = v[:i+1], v[i+1:]
			break
		}
	}

	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return Length{}, err
	}

	return Length{Value: n, Units: units}, nil
}

func (l *Length) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	*l, err = parseLength(tokens[0])
	return err
}

type LengthIdent struct {
	Length Length
	Ident  string
}

func (lp *LengthIdent) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.NumberToken, css.DimensionToken:
		l, err := parseLength(token)
		if err != nil {
			return err
		}
		lp.Length = l
		return nil
	case css.IdentToken:
		lp.Ident = token.Value
		return nil
	default:
		return errors.New("expected a length or identifier")
	}
}

type LengthPercentage struct {
	Length     Length
	Percentage float64
}

func parseLengthPercentage(token cssToken) (LengthPercentage, error) {
	switch token.Type {
	case css.NumberToken, css.DimensionToken:
		l, err := parseLength(token)
		if err != nil {
			return LengthPercentage{}, err
		}
		return LengthPercentage{Length: l}, nil
	case css.PercentageToken:
		p, err := strconv.ParseFloat(token.Value[:len(token.Value)-1], 64)
		if err != nil {
			return LengthPercentage{}, err
		}
		return LengthPercentage{Percentage: p / 100}, nil
	default:
		return LengthPercentage{}, errors.New("expected a length or percentage")
	}
}

func (lp *LengthPercentage) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	*lp, err = parseLengthPercentage(tokens[0])
	return err
}

type BoxLengthPercentage struct {
	LengthPercentage

	Value string
}

func (blp *BoxLengthPercentage) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.IdentToken:
		switch token.Value {
		case "auto", "inherit":
			blp.Value = token.Value
		default:
			return errors.New("expected 'auto' or 'inherit'")
		}
		return nil
	case css.NumberToken, css.DimensionToken, css.PercentageToken:
		lp, err := parseLengthPercentage(token)
		if err != nil {
			return err
		}
		blp.LengthPercentage = lp
		return nil
	default:
		return errors.New("expected a length, percentage, 'auto', or 'inherit'")
	}
}

type LengthPercentageNumber struct {
	LengthPercentage

	Number float64
}

func parseLengthPercentageNumber(token cssToken) (LengthPercentageNumber, error) {
	switch token.Type {
	case css.NumberToken:
		n, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return LengthPercentageNumber{}, err
		}
		return LengthPercentageNumber{Number: n}, nil
	case css.DimensionToken, css.PercentageToken:
		lp, err := parseLengthPercentage(token)
		if err != nil {
			return LengthPercentageNumber{}, err
		}
		return LengthPercentageNumber{LengthPercentage: lp}, nil
	default:
		return LengthPercentageNumber{}, errors.New("expected a length, percentage, or number")
	}
}

func (lpn *LengthPercentageNumber) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	*lpn, err = parseLengthPercentageNumber(tokens[0])
	return err
}

type LengthPercentageNumbers struct {
	Values []LengthPercentageNumber
}

func (ns *LengthPercentageNumbers) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}

	var values []LengthPercentageNumber
	for len(tokens) > 0 {
		v, err := parseLengthPercentageNumber(tokens[0])
		if err != nil {
			return err
		}
		values, tokens = append(values, v), tokens[1:]

		if len(tokens) > 0 && tokens[0].Type == css.CommaToken {
			tokens = tokens[1:]
		}
	}

	ns.Values = values
	return nil
}

type LengthPercentageIdent struct {
	LengthPercentage

	Ident string
}

func parseLengthPercentageIdent(token cssToken) (LengthPercentageIdent, error) {
	switch token.Type {
	case css.IdentToken:
		return LengthPercentageIdent{Ident: token.Value}, nil
	case css.DimensionToken, css.PercentageToken:
		lp, err := parseLengthPercentage(token)
		if err != nil {
			return LengthPercentageIdent{}, err
		}
		return LengthPercentageIdent{LengthPercentage: lp}, nil
	default:
		return LengthPercentageIdent{}, errors.New("expected a length, percentage, or identifier")
	}
}

func (lpi *LengthPercentageIdent) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	*lpi, err = parseLengthPercentageIdent(tokens[0])
	return err
}

type LengthPercentageNumberIdent struct {
	LengthPercentageNumber

	Ident string
}

func (lpni *LengthPercentageNumberIdent) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}
	if len(tokens) != 1 {
		return errors.New("unexpected token")
	}

	token := tokens[0]
	switch token.Type {
	case css.IdentToken:
		lpni.Ident = token.Value
		return nil
	case css.NumberToken, css.DimensionToken, css.PercentageToken:
		lpn, err := parseLengthPercentageNumber(token)
		if err != nil {
			return err
		}
		lpni.LengthPercentageNumber = lpn
		return nil
	default:
		return errors.New("expected a length, percentage, number, or identifier")
	}
}

// Color represents a color.
type Color struct {
	Value color.Color
}

func parseColorFunction(tokens []cssToken) (color.Color, error) {
	fn, arity := tokens[0].Value, 0
	switch fn {
	case "rgb(", "hsl(":
		arity = 3
	case "rgba(", "hsla(":
		arity = 4
	default:
		return nil, fmt.Errorf("unknown color function %v", tokens[0].Value)
	}

	tokens = tokens[1:]
	if len(tokens) == 0 {
		return nil, errors.New("expected a number or ')'")
	}

	args := make([]byte, 0, arity)
	if tokens[0].Type == css.RightParenthesisToken {
		tokens = tokens[1:]
	} else {
		for {
			switch tokens[0].Type {
			case css.NumberToken:
				n, err := strconv.ParseUint(tokens[0].Value, 10, 8)
				if err != nil {
					return nil, err
				}
				args, tokens = append(args, byte(n)), tokens[1:]
			case css.PercentageToken:
				n, err := strconv.ParseUint(tokens[0].Value, 10, 8)
				if err != nil {
					return nil, err
				}
				if n > 100 {
					return nil, fmt.Errorf("percentage %v%% is out of range", n)
				}
				args, tokens = append(args, byte(255*100/n)), tokens[1:]
			default:
				return nil, errors.New("expected a number or percentage")
			}

			if len(tokens) == 0 {
				return nil, errors.New("expected ',' or ')'")
			}

			if tokens[0].Type == css.RightParenthesisToken {
				tokens = tokens[1:]
				break
			}
			if tokens[0].Type != css.CommaToken {
				return nil, errors.New("expected ','")
			}
			tokens = tokens[1:]
		}
	}

	if len(tokens) != 0 {
		return nil, errors.New("garbage after function call")
	}

	if len(args) != arity {
		return nil, fmt.Errorf("%v requires %v arguments", fn, arity)
	}

	var r, g, b, a byte
	switch fn {
	case "rgb(":
		r, g, b, a = args[0], args[1], args[2], 255
	case "rgba(":
		r, g, b, a = args[0], args[1], args[2], args[3]
	case "hsl(":
		r, g, b = hslToRGB(args[0], args[1], args[2])
		a = 255
	case "hsla(":
		r, g, b = hslToRGB(args[0], args[1], args[2])
		a = args[3]
	}
	return &color.RGBA{R: r, G: g, B: b, A: a}, nil
}

func parseHexColor(v string) (color.Color, error) {
	switch len(v) {
	case 3:
		v = string([]byte{v[0], v[0], v[1], v[1], v[2], v[2]})
	case 6:
		//  OK
	default:
		return nil, fmt.Errorf("invalid hex color %v", v)
	}

	bytes, err := hex.DecodeString(v)
	if err != nil {
		return nil, err
	}

	return &color.RGBA{R: bytes[0], G: bytes[1], B: bytes[2], A: 255}, nil
}

func parseColor(tokens []cssToken) (color.Color, error) {
	if tokens[0].Type == css.FunctionToken {
		return parseColorFunction(tokens)
	}

	if len(tokens) != 1 {
		return nil, errors.New("unexpected token")
	}

	switch tokens[0].Type {
	case css.IdentToken:
		ident := tokens[0].Value
		color, ok := cssColors[ident]
		if !ok {
			return nil, fmt.Errorf("unknown color %v", ident)
		}
		return color, nil
	case css.HashToken:
		return parseHexColor(tokens[0].Value[1:])
	default:
		return nil, errors.New("expected an identifier or hex color")
	}
}

func (c *Color) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}

	color, err := parseColor(tokens)
	if err != nil {
		return err
	}
	c.Value = color
	return nil
}

// Paint represents a paint color.
type Paint struct {
	Context string
	URL     string
	Color   color.Color
}

func (p *Paint) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		p.Color = color.Black
		return nil
	}

	if tokens[0].Type == css.URLToken {
		url := tokens[0].Value
		p.URL = url[len("url(") : len(url)-1]

		tokens = tokens[1:]
		if len(tokens) == 0 {
			return nil
		}
	}

	if tokens[0].Type == css.IdentToken {
		if len(tokens) != 1 {
			return errors.New("unexpected token")
		}

		ident := tokens[0].Value
		switch ident {
		case "context-fill", "context-stroke":
			p.Context = ident
			return nil
		case "none":
			p.Color = color.Transparent
			return nil
		}
	}

	color, err := parseColor(tokens)
	if err != nil {
		return err
	}
	p.Color = color
	return nil
}

type FontFamily struct {
	Values []string
}

func (ff *FontFamily) UnmarshalText(text []byte) error {
	tokens, err := cssTokens(string(text))
	if err != nil {
		return err
	}

	var values []string
	for len(tokens) > 0 {
		switch tokens[0].Type {
		case css.StringToken:
			values, tokens = append(values, tokens[0].Value), tokens[1:]
		case css.IdentToken:
			var f strings.Builder
			for len(tokens) > 0 {
				if tokens[0].Type == css.CommaToken {
					break
				}
				switch tokens[0].Type {
				case css.WhitespaceToken, css.IdentToken:
					f.WriteString(tokens[0].Value)
					tokens = tokens[1:]
				default:
					return errors.New("expected an identifier or ','")
				}
			}
			values = append(values, f.String())
		default:
			return errors.New("expected a string or identifier")
		}

		if len(tokens) == 0 {
			break
		}
		if tokens[0].Type != css.CommaToken {
			return errors.New("expected a ','")
		}
		tokens = tokens[1:]
	}

	ff.Values = values
	return nil
}

// TODO

type ClipPath string
type Cursor string
type DashArray string
type FilterList string
type Mask string
type Transform string
type VectorEffect string
