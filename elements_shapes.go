package svg

import (
	"bufio"
	"bytes"
	"encoding/xml"
)

// Rect represents an SVG `rect` element.
type Rect struct {
	ElementAttributes

	XMLName xml.Name `xml:"rect"`

	PathLength float64 `xml:"pathLength,attr"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`

	Rx *BoxLengthPercentage `xml:"rx,attr"`
	Ry *BoxLengthPercentage `xml:"ry,attr"`

	Children []any `xml:",any"`
}

func (Rect) isElement() {}

// Circle represents an SVG `circle` element.
type Circle struct {
	ElementAttributes

	XMLName xml.Name `xml:"circle"`

	PathLength float64 `xml:"pathLength,attr"`

	Cx LengthPercentage `xml:"cx,attr"`
	Cy LengthPercentage `xml:"cy,attr"`

	R LengthPercentage `xml:"r,attr"`

	Children []any `xml:",any"`
}

func (Circle) isElement() {}

// Ellipse represents an SVG `ellipse` element.
type Ellipse struct {
	ElementAttributes

	XMLName xml.Name `xml:"ellipse"`

	PathLength float64 `xml:"pathLength,attr"`

	Cx LengthPercentage `xml:"cx,attr"`
	Cy LengthPercentage `xml:"cy,attr"`

	Rx BoxLengthPercentage `xml:"rx,attr"`
	Ry BoxLengthPercentage `xml:"ry,attr"`

	Children []any `xml:",any"`
}

func (Ellipse) isElement() {}

// Line represents an SVG `line` element.
type Line struct {
	ElementAttributes

	XMLName xml.Name `xml:"line"`

	PathLength float64 `xml:"pathLength,attr"`

	X1 LengthPercentageNumber `xml:"x1,attr"`
	Y1 LengthPercentageNumber `xml:"y1,attr"`
	X2 LengthPercentageNumber `xml:"x2,attr"`
	Y2 LengthPercentageNumber `xml:"y2,attr"`

	Children []any `xml:",any"`
}

func (Line) isElement() {}

type PolyPoints []Point

func (p *PolyPoints) UnmarshalText(text []byte) error {
	points, err := parseCoordinatePairSequence(bufio.NewReader(bytes.NewReader(text)))
	if err != nil {
		return err
	}
	*p = points
	return nil
}

// Polyline represents an SVG `polyline` element.
type Polyline struct {
	ElementAttributes

	XMLName xml.Name `xml:"polyline"`

	Points PolyPoints `xml:"points,attr"`

	Children []any `xml:",any"`
}

func (Polyline) isElement() {}

// Polygon represents an SVG `polygon` element.
type Polygon struct {
	ElementAttributes

	XMLName xml.Name `xml:"polygon"`

	Points PolyPoints `xml:"points,attr"`

	Children []any `xml:",any"`
}

func (Polygon) isElement() {}
