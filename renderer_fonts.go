package svg

import (
	_ "github.com/flopp/go-findfont"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
)

type fontWeight struct {
	weight font.Weight
	normal interface{}
	italic interface{}
}

type fontFamily struct {
	newFaceFunc func(font interface{}, points float64, hints font.Hinting) (font.Face, error)

	// TODO: stretch

	weights []fontWeight
}

func (ff *fontFamily) newFace(weight font.Weight, style font.Style, points float64, hints font.Hinting) (font.Face, error) {
	var chosen *fontWeight
	for i := range ff.weights {
		w := &ff.weights[i]
		if w.weight >= weight {
			chosen = w
			break
		}
	}
	if chosen == nil {
		chosen = &ff.weights[len(ff.weights)-1]
	}

	if style == font.StyleItalic {
		return ff.newFaceFunc(chosen.italic, points, hints)
	}
	return ff.newFaceFunc(chosen.normal, points, hints)
}

func newOpentypeFontFamily(weights []fontWeight) *fontFamily {
	return &fontFamily{
		newFaceFunc: func(font interface{}, points float64, hints font.Hinting) (font.Face, error) {
			return opentype.NewFace(font.(*sfnt.Font), &opentype.FaceOptions{
				Size:    points,
				DPI:     72, // TODO: DPI
				Hinting: hints,
			})
		},
		weights: weights,
	}
}

func loadOpentypeFontFamily(weights []fontWeight) (*fontFamily, error) {
	for i := range weights {
		w := &weights[i]
		if w.normal != nil {
			f, err := opentype.Parse(w.normal.([]byte))
			if err != nil {
				return nil, err
			}
			w.normal = f
		}

		if w.italic != nil {
			f, err := opentype.Parse(w.italic.([]byte))
			if err != nil {
				return nil, err
			}
			w.italic = f
		}
	}

	return newOpentypeFontFamily(weights), nil
}

func mustOpentypeFontFamily(weights []fontWeight) *fontFamily {
	f, err := loadOpentypeFontFamily(weights)
	if err != nil {
		panic(err)
	}
	return f
}

var goProportional = mustOpentypeFontFamily([]fontWeight{
	{
		weight: font.WeightNormal,
		normal: goregular.TTF,
		italic: goitalic.TTF,
	},
	{
		weight: font.WeightMedium,
		normal: gomedium.TTF,
		italic: gomediumitalic.TTF,
	},
	{
		weight: font.WeightBold,
		normal: gobold.TTF,
		italic: gobolditalic.TTF,
	},
})

var goMonospace = mustOpentypeFontFamily([]fontWeight{
	{
		weight: font.WeightNormal,
		normal: gomono.TTF,
		italic: gomonoitalic.TTF,
	},
	{
		weight: font.WeightBold,
		normal: gomonobold.TTF,
		italic: gomonobolditalic.TTF,
	},
})

func defaultFonts() map[string]*fontFamily {
	return map[string]*fontFamily{
		"serif":         goProportional,
		"sans-serif":    goProportional,
		"monospace":     goMonospace,
		"cursive":       goProportional,
		"fantasy":       goProportional,
		"system-ui":     goProportional,
		"ui-serif":      goProportional,
		"ui-sans-serif": goProportional,
		"ui-monospace":  goMonospace,
		"ui-rounded":    goProportional,
		"math":          goProportional,
		"emoji":         goProportional,
		"fangsong":      goProportional,
	}
}

func (r *renderer) resolveFontFamily(family *FontFamily) (*fontFamily, error) {
	for _, v := range family.Values {
		if f, ok := r.fonts[v]; ok {
			return f, nil
		}
	}

	return goProportional, nil
}
