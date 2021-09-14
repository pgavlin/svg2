package svg

import (
	"encoding/xml"
)

// Text represents an SVG `text` element.
type Text struct {
	ElementAttributes

	XMLName xml.Name `xml:"text"`

	X  LengthPercentageNumbers `xml:"x,attr"`
	Y  LengthPercentageNumbers `xml:"y,attr"`
	Dx LengthPercentageNumbers `xml:"dx,attr"`
	Dy LengthPercentageNumbers `xml:"dy,attr"`

	Rotate LengthPercentageNumbers `xml:"rotate,attr"`

	TextLength LengthPercentageNumber `xml:"textLength,attr"`

	Value string `xml:",chardata"`
}

func (Text) isElement() {}

// TSpan represents an SVG `tspan` element.
type TSpan struct {
	ElementAttributes

	XMLName xml.Name `xml:"tspan"`

	X  LengthPercentageNumbers `xml:"x,attr"`
	Y  LengthPercentageNumbers `xml:"y,attr"`
	Dx LengthPercentageNumbers `xml:"dx,attr"`
	Dy LengthPercentageNumbers `xml:"dy,attr"`

	Rotate LengthPercentageNumbers `xml:"rotate,attr"`

	TextLength   LengthPercentageNumber `xml:"textLength,attr"`
	LengthAdjust string                 `xml:"lengthAdjust,attr"`

	Value string `xml:",chardata"`
}

func (TSpan) isElement() {}

// TextPath represents an SVG `textPath` element.
type TextPath struct {
	ElementAttributes

	XMLName xml.Name `xml:"textPath"`

	Path        PathData         `xml:"path,attr"`
	Href        string           `xml:"href,attr"`
	StartOffset LengthPercentage `xml:"startOffset,attr"`
	Method      string           `xml:"method,attr"`
	Spacing     string           `xml:"spacing,attr"`
	Side        string           `xml:"side,attr"`

	TextLength   LengthPercentageNumber `xml:"textLength,attr"`
	LengthAdjust string                 `xml:"lengthAdjust,attr"`

	Value string `xml:",chardata"`
}

func (TextPath) isElement() {}
