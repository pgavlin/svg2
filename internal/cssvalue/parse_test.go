package cssvalue

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "orphans", input: "<integer>"},
		{name: "text-align", input: "left | right | center | justify"},
		{name: "padding-top", input: "<length> | <percentage>"},
		{name: "outline-color", input: "<color> | invert"},
		{name: "text-decoration", input: "none | underline || overline || line-through || blink"},
		{name: "font-family", input: "[ <family-name> | <generic-family> ]#"},
		{name: "border-width", input: "[ <length> | thick | medium | thin ]{1,4}"},
		{name: "box-shadow", input: "[ inset? && <length>{2,4} && <color>? ]# | none"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := ParseGrammar(strings.NewReader(c.input))
			assert.NoError(t, err)
		})
	}
}
