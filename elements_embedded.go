package svg

import (
	"encoding/xml"
)

// Image represents an SVG `image` element.
type Image struct {
	ElementAttributes

	XMLName xml.Name `xml:"image"`

	PreserveAspectRatio string `xml:"preserveAspectRatio,attr"`
	Href                string `xml:"href,attr"`
	CrossOrigin         string `xml:"crossOrigin,attr"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`
}

func (Image) isElement() {}

// ForeignObject represents an SVG `foreignObject` element.
type ForeignObject struct {
	ElementAttributes

	XMLName xml.Name `xml:"foreignObject"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`

	InnerXML []byte `xml:",innerxml"`
}

func (ForeignObject) isElement() {}
