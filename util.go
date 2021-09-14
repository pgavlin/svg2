package svg

import (
	"io"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type cssToken struct {
	Type  css.TokenType
	Value string
}

func cssTokens(s string) ([]cssToken, error) {
	var tokens []cssToken

	l := css.NewLexer(parse.NewInput(strings.NewReader(s)))
	for {
		typ, value := l.Next()
		if typ == css.ErrorToken {
			if l.Err() == io.EOF {
				break
			}
			return nil, l.Err()
		}
		tokens = append(tokens, cssToken{Type: typ, Value: string(value)})
	}

	return tokens, nil
}

func hueToRGB(m1, m2, h float64) byte {
	switch {
	case h < 0:
		h += 1
	case h > 1:
		h -= 1
	}

	switch {
	case h*6 < 1:
		return byte(m1 + (m2-m1)*h*6*255)
	case h*2 < 1:
		return byte(m2 * 255)
	case h*3 < 2:
		return byte(m1 + (m2-m1)*(2/3-h)*6*255)
	}
	return byte(m1 * 255)
}

func hslToRGB(h, s, l byte) (r, g, b byte) {
	hf, sf, lf := float64(h)/255, float64(s)/255, float64(l)/255

	var m2 float64
	if lf <= 0.5 {
		m2 = lf * (sf + 1)
	} else {
		m2 = lf + sf - lf*sf
	}

	m1 := lf*2 - m2
	return hueToRGB(m1, m2, hf+1/3), hueToRGB(m1, m2, hf), hueToRGB(m1, m2, hf-1/3)
}
