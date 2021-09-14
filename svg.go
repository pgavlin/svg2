package svg

import "encoding/xml"

// SVG represents an SVG document.
type SVG struct {
	XMLName xml.Name `xml:"svg"`

	Style *Style `xml:"style"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`

	Children []any `xml:",any"`
}

// Style represents an SVG `style` element.
type Style struct {
	XMLName xml.Name `xml:"style"`

	Type  string `xml:"type,attr"`
	Media string `xml:"media,attr"`
	Title string `xml:"title,attr`

	Style string `xml:",chardata"`
}

func (*Style) isElement() {}
