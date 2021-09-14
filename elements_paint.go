package svg

import (
	"encoding/xml"
	"image/color"
)

type PaintElement interface {
	isPaintElement()
}

// Gradient holds fields common to the `linearGradient` and `radialGradient` elements.
type Gradient struct {
	ElementAttributes

	GradientUnits     Units       `xml:"gradientUnits,attr"`
	GradientTransform []Transform `xml:"gradientTransform,attr"`
	SpreadMethod      string      `xml:"spreadMethod,attr"`
	Href              string      `xml:"href,attr"`

	Stops []GradientStop `xml:"stop"`
}

func (Gradient) isElement()      {}
func (Gradient) isPaintElement() {}

// GradientStop represents an SVG `stop` element.
type GradientStop struct {
	ElementAttributes

	XMLName xml.Name `xml:"stop"`

	Offset  NumberPercentage `xml:"offset,attr"`
	Color   Color            `xml:"stop-color,attr"`
	Opacity NumberPercentage `xml:"stop-opacity,attr"`
}

// LinearGradient represents an SVG `linearGradient` element.
type LinearGradient struct {
	Gradient

	XMLName xml.Name `xml:"linearGradient"`

	X1 LengthPercentage `xml:"x1,attr"`
	Y1 LengthPercentage `xml:"y1,attr"`
	X2 LengthPercentage `xml:"x2,attr"`
	Y2 LengthPercentage `xml:"y2,attr"`
}

// RadialGradient represents an SVG `radialGradient` element.
type RadialGradient struct {
	Gradient

	XMLName xml.Name `xml:"radialGradient"`

	Cx LengthPercentage `xml:"cx,attr"`
	Cy LengthPercentage `xml:"cy,attr"`
	R  LengthPercentage `xml:"r,attr"`
	Fx LengthPercentage `xml:"fx,attr"`
	Fy LengthPercentage `xml:"fy,attr"`
	Fr LengthPercentage `xml:"fr,attr"`
}

// Pattern represents an SVG `pattern` element.
type Pattern struct {
	ElementAttributes

	XMLName xml.Name `xml:"pattern"`

	ViewBox             string `xml:"viewBox,attr"`
	PreserveAspectRatio string `xml:"preserveAspectRatio,attr"`

	X      Length `xml:"x,attr"`
	Y      Length `xml:"y,attr"`
	Width  Length `xml:"width,attr"`
	Height Length `xml:"height,attr"`

	PatternUnits        Units `xml:"patternUnits,attr"`
	PatternContentUnits Units `xml:"patternContentUnits,attr"`

	PatternTransform []Transform `xml:"patternTransform,attr"`

	Href string `xml:"href,attr"`

	Children []any `xml:",any"`
}

func (*Pattern) isElement()      {}
func (*Pattern) isPaintElement() {}

// SolidColor represents a solid paint color.
type SolidColor struct {
	color.Color
}

func (SolidColor) isPaintElement() {}
