package svg

import (
	"encoding/xml"
	"fmt"
)

// Element represents an SVG element.
type Element interface {
	id() string
	attrs() *ElementAttributes

	isElement()
}

type any struct {
	X Element
}

func Any(x Element) any {
	return any{X: x}
}

func (a *any) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(a.X, start)
}

func (a *any) UnmarshalXML(d *xml.Decoder, s xml.StartElement) error {
	switch s.Name.Local {
	case "g":
		a.X = &Grouping{}
	case "defs":
		a.X = &Defs{}
	case "symbol":
		a.X = &Symbol{}
	case "use":
		a.X = &Use{}
	case "switch":
		a.X = &Switch{}
	case "marker":
		a.X = &Marker{}
	case "linearGradient":
		a.X = &LinearGradient{}
	case "radialGradient":
		a.X = &RadialGradient{}
	case "pattern":
		a.X = &Pattern{}
	case "path":
		a.X = &Path{}
	case "rect":
		a.X = &Rect{}
	case "circle":
		a.X = &Circle{}
	case "ellipse":
		a.X = &Ellipse{}
	case "line":
		a.X = &Line{}
	case "polyline":
		a.X = &Polyline{}
	case "polygon":
		a.X = &Polygon{}
	case "text":
		a.X = &Text{}
	case "tspan":
		a.X = &TSpan{}
	case "image":
		a.X = &Image{}
	case "foreignObject":
		a.X = &ForeignObject{}
	default:
		return fmt.Errorf("unrecognized element %v:%v", s.Name.Space, s.Name.Local)
	}

	return d.DecodeElement(a.X, &s)
}

// ElementAttributes contains standard SVG element attributes.
type ElementAttributes struct {
	ID string `xml:"id,attr"`

	AlignmentBaseline         Ident                        `xml:"alignment-baseline,attr"`
	BaselineShift             *LengthPercentageIdent       `xml:"baseline-shift,attr"`
	ClipPath                  *ClipPath                    `xml:"clip-path,attr"`
	ClipRule                  Ident                        `xml:"clip-rule,attr"`
	Color                     *Color                       `xml:"color,attr"`
	ColorInterpolation        Ident                        `xml:"color-interpolation,attr"`
	ColorInterpolationFilters Ident                        `xml:"color-interpolation-filters,attr"`
	ColorRendering            Ident                        `xml:"color-rendering,attr"`
	Cursor                    *Cursor                      `xml:"cursor,attr"`
	Direction                 Ident                        `xml:"direction,attr"`
	Display                   Ident                        `xml:"display,attr"`
	DominantBaseline          Ident                        `xml:"dominant-baseline,attr"`
	Fill                      *Paint                       `xml:"fill,attr"`
	FillOpacity               *NumberPercentage            `xml:"fill-opacity,attr"`
	FillRule                  Ident                        `xml:"fill-rule,attr"`
	Filter                    *FilterList                  `xml:"filter,attr"`
	FloodColor                *Color                       `xml:"flood-color,attr"`
	FloodOpacity              *NumberPercentage            `xml:"flood-opacity,attr"`
	FontFamily                *FontFamily                  `xml:"font-family,attr"`
	FontSize                  *LengthPercentageNumberIdent `xml:"font-size,attr"`
	FontSizeAdjust            *NumberIdent                 `xml:"font-size-adjust,attr"`
	FontStretch               Ident                        `xml:"font-stretch,attr"`
	FontStyle                 Ident                        `xml:"font-style,attr"`
	FontVariant               Ident                        `xml:"font-variant,attr"`
	FontWeight                *NumberIdent                 `xml:"font-weight,attr"`
	ImageRendering            Ident                        `xml:"image-rendering,attr"`
	LetterSpacing             *LengthIdent                 `xml:"letter-spacing,attr"`
	LightingColor             *Color                       `xml:"lighting-color,attr"`
	MarkerEnd                 *URLIdent                    `xml:"marker-end,attr"`
	MarkerMid                 *URLIdent                    `xml:"marker-mid,attr"`
	MarkerStart               *URLIdent                    `xml:"marker-start,attr"`
	Mask                      *Mask                        `xml:"mask,attr"`
	Opacity                   *NumberIdent                 `xml:"opacity,attr"`
	Overflow                  Ident                        `xml:"overflow,attr"`
	PaintOrder                Ident                        `xml:"paint-order,attr"`
	PointerEvents             Ident                        `xml:"pointer-events,attr"`
	ShapeRendering            Ident                        `xml:"shape-rendering,attr"`
	Stroke                    *Paint                       `xml:"stroke,attr"`
	StrokeDasharray           *DashArray                   `xml:"stroke-dasharray,attr"`
	StrokeDashoffset          *LengthPercentage            `xml:"stroke-dashoffset,attr"`
	StrokeLinecap             Ident                        `xml:"stroke-linecap,attr"`
	StrokeLinejoin            Ident                        `xml:"stroke-linejoin,attr"`
	StrokeMiterlimit          *float64                     `xml:"stroke-miterlimit,attr"`
	StrokeOpacity             *NumberPercentage            `xml:"stroke-opacity,attr"`
	StrokeWidth               *LengthPercentage            `xml:"stroke-width,attr"`
	TextAnchor                Ident                        `xml:"text-anchor,attr"`
	TextDecoration            Ident                        `xml:"text-decoration,attr"`
	TextOverflow              Ident                        `xml:"text-overflow,attr"`
	TextRendering             Ident                        `xml:"text-rendering,attr"`
	Transform                 []Transform                  `xml:"transform,attr"`
	UnicodeBidi               Ident                        `xml:"unicode-bidi,attr"`
	VectorEffect              *VectorEffect                `xml:"vector-effect,attr"`
	Visibility                Ident                        `xml:"visibility,attr"`
	WhiteSpace                Ident                        `xml:"white-space,attr"`
	WordSpacing               *LengthIdent                 `xml:"word-spacing,attr"`
	WritingMode               Ident                        `xml:"writing-mode,attr"`
}

func (ea *ElementAttributes) id() string {
	return ea.ID
}

func (ea *ElementAttributes) attrs() *ElementAttributes {
	return ea
}

// Grouping represents an SVG `g` element.
type Grouping struct {
	ElementAttributes

	XMLName xml.Name `xml:"g"`

	Children []any `xml:",any"`
}

func (Grouping) isElement() {}

// Defs represents an SVG `defs` element.
type Defs struct {
	ElementAttributes

	XMLName xml.Name `xml:"defs"`

	Children []any `xml:",any"`
}

func (Defs) isElement() {}

// Symbol represents an SVG `symbol` element.
type Symbol struct {
	ElementAttributes

	XMLName xml.Name `xml:"symbol"`

	PreserveAspectRatio string `xml:"preserveAspectRatio,attr"`
	ViewBox             string `xml:"viewBox,attr"`

	RefX Length `xml:"refX,attr"`
	RefY Length `xml:"refY,attr"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`

	Children []any `xml:",any"`
}

func (Symbol) isElement() {}

// Use represents an SVG `use` element.
type Use struct {
	ElementAttributes

	XMLName xml.Name `xml:"use"`

	Href string `xml:"href,attr"`

	X      LengthPercentage    `xml:"x,attr"`
	Y      LengthPercentage    `xml:"y,attr"`
	Width  BoxLengthPercentage `xml:"width,attr"`
	Height BoxLengthPercentage `xml:"height,attr"`
}

func (Use) isElement() {}

// Switch represents an SVG `switch` element.
type Switch struct {
	ElementAttributes

	XMLName xml.Name `xml:"switch"`

	Children []any `xml:",any"`
}

func (Switch) isElement() {}

// Marker represents an SVG `marker` element.
type Marker struct {
	ElementAttributes

	XMLName xml.Name `xml:"marker"`

	PreserveAspectRatio string `xml:"preserveAspectRatio,attr"`
	ViewBox             string `xml:"viewBox,attr"`

	RefX Length `xml:"refX,attr"`
	RefY Length `xml:"refY,attr"`

	MarkerUnits  string                 `xml:"markerUnits,attr"`
	MarkerWidth  LengthPercentageNumber `xml:"markerWidth,attr"`
	MarkerHeight LengthPercentageNumber `xml:"markerHeight,attr"`

	Orient string `xml:"orient,attr"`

	Children []any `xml:",any"`
}

func (Marker) isElement() {}
