package svg

import (
	"encoding/xml"
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

type SVGImage struct {
	doc *SVG
	ctx *gg.Context
}

func (i *SVGImage) SVG() *SVG {
	return i.doc
}

func (i *SVGImage) Context() *gg.Context {
	return i.ctx
}

func (i *SVGImage) ColorModel() color.Model {
	return i.ctx.Image().ColorModel()
}

func (i *SVGImage) Bounds() image.Rectangle {
	return i.ctx.Image().Bounds()
}

func (i *SVGImage) At(x, y int) color.Color {
	return i.ctx.Image().At(x, y)
}

func (i *SVGImage) Scale(factor float64) (*SVGImage, error) {
	ctx := NewScaledContext(i.doc, factor)
	if err := Render(ctx, i.doc); err != nil {
		return nil, err
	}
	return &SVGImage{
		doc: i.doc,
		ctx: ctx,
	}, nil
}

func Decode(r io.Reader) (image.Image, error) {
	var doc SVG
	if err := xml.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}

	ctx := NewContext(&doc)
	if err := Render(ctx, &doc); err != nil {
		return nil, err
	}

	return &SVGImage{
		doc: &doc,
		ctx: ctx,
	}, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	img, err := Decode(r)
	if err != nil {
		return image.Config{}, err
	}
	bounds := img.Bounds()
	cfg := image.Config{
		ColorModel: img.ColorModel(),
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
	}
	return cfg, nil
}

func init() {
	image.RegisterFormat("svg", "<svg", Decode, DecodeConfig)
}
