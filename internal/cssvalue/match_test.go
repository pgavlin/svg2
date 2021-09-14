package cssvalue

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/parse/v2/css"
)

func TestMatch(t *testing.T) {
	lengthPercentage := &OneOf{
		Operands: []Term{
			&BasicType{Name: "length"},
			&BasicType{Name: "percentage"},
		},
	}

	ctx := Context{
		NonTerminals: map[string]Term{
			"family-name": &OneOf{
				Operands: []Term{
					&BasicType{Name: "string"},
					&BasicType{Name: "custom-ident"},
				},
			},
			"generic-family": &OneOf{
				Operands: []Term{
					Keyword("serif"),
					Keyword("sans-serif"),
					Keyword("cursive"),
					Keyword("fantasy"),
					Keyword("monospace"),
				},
			},
		},
		Functions: map[string]*Function{
			"rgba": &Function{
				Params: []Term{
					lengthPercentage,
					lengthPercentage,
					lengthPercentage,
					lengthPercentage,
				},
				Type: &BasicType{
					Name: "color",
				},
			},
		},
	}

	cases := []struct {
		name  string
		input string
		value string
	}{
		{name: "orphans", input: "<integer>", value: "3"},
		{name: "text-align", input: "left | right | center | justify", value: "center"},
		{name: "padding-top", input: "<length> | <percentage>", value: "5%"},
		{name: "outline-color", input: "<color> | invert", value: "#fefefe"},
		{name: "text-decoration", input: "none | underline || overline || line-through || blink", value: "overline underline"},
		{name: "font-family", input: "[ <family-name> | <generic-family> ]#", value: `"Gill Sans", Futura, sans-serif`},
		{name: "border-width", input: "[ <length> | thick | medium | thin ]{1,4}", value: "2px medium 4px"},
		{name: "box-shadow", input: "[ inset? && <length>{2,4} && <color>? ]# | none", value: "3px 3px rgba(50%, 50%, 50%, 50%), lemonchiffon 0 0 4px inset"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			term, err := ParseGrammar(strings.NewReader(c.input))
			require.NoError(t, err)

			captures := Match(&ctx, term, strings.NewReader(c.value))
			assert.NotEqual(t, 0, len(captures))
		})
	}
}

func TestCaptures(t *testing.T) {
	ctx := Context{
		NonTerminals: map[string]Term{
			"min-x":  &BasicType{Name: "number"},
			"min-y":  &BasicType{Name: "number"},
			"width":  &BasicType{Name: "number"},
			"height": &BasicType{Name: "number"},
			"family-name": &OneOf{
				Operands: []Term{
					&BasicType{Name: "string"},
					&BasicType{Name: "custom-ident"},
				},
			},
			"generic-family": &OneOf{
				Operands: []Term{
					Keyword("serif"),
					Keyword("sans-serif"),
					Keyword("cursive"),
					Keyword("fantasy"),
					Keyword("monospace"),
				},
			},
		},
	}

	cases := []struct {
		name     string
		input    string
		value    string
		captures []Capture
	}{
		{
			name:  "font-family",
			input: "[ <family-name> | <generic-family> ]#",
			value: `"Gill Sans", Futura, sans-serif`,
			captures: []Capture{
				{Name: "family-name", Values: []Token{{Type: css.StringToken, Value: `"Gill Sans"`}}},
				{Name: "family-name", Values: []Token{{Type: css.IdentToken, Value: "Futura"}}},
				{Name: "family-name", Values: []Token{{Type: css.IdentToken, Value: "sans-serif"}}}, // TODO: this ought to be generic-family, but custom-ident is matching sans-serif
			},
		},
		{
			name:  "viewBox",
			input: "[<min-x>,? <min-y>,? <width>,? <height>]",
			value: "0 0 200 300",
			captures: []Capture{
				{Name: "min-x", Values: []Token{{Type: css.NumberToken, Value: "0"}}},
				{Name: "min-y", Values: []Token{{Type: css.NumberToken, Value: "0"}}},
				{Name: "width", Values: []Token{{Type: css.NumberToken, Value: "200"}}},
				{Name: "height", Values: []Token{{Type: css.NumberToken, Value: "300"}}},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			term, err := ParseGrammar(strings.NewReader(c.input))
			require.NoError(t, err)

			captures := Match(&ctx, term, strings.NewReader(c.value))
			assert.Equal(t, c.captures, captures[:len(captures)-1])
		})
	}
}
