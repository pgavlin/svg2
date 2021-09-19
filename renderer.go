package svg

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// NewContext creates a new render context for an SVG document.
func NewContext(svg *SVG) *gg.Context {
	if svg.Width.Length.Value != 0 && svg.Height.Length.Value != 0 {
		return gg.NewContext(int(svg.Width.Length.Value), int(svg.Height.Length.Value))
	}
	return gg.NewContext(1024, 1024)
}

// NewScaledContext creates a new render context for an SVG document with the given scaling factor.
func NewScaledContext(svg *SVG, scale float64) *gg.Context {
	width, height := 1024, 1024
	if svg.Width.Length.Value != 0 && svg.Height.Length.Value != 0 {
		width, height = int(svg.Width.Length.Value), int(svg.Height.Length.Value)
	}

	ctx := gg.NewContext(int(float64(width)*scale), int(float64(height)*scale))
	ctx.Scale(scale, scale)
	return ctx
}

// Render renders an SVG document to the given context.
func Render(ctx *gg.Context, svg *SVG) error {
	// Graphics are painted and composited in rendering-tree order, subject to re-ordering based on the paint-order property. Note that elements that have no visual paint may still be in the rendering tree.
	//
	// shadow DOM elements, such as those generated by ‘use’ elements or by cross-references between paint servers;
	//
	// With the exception of the ‘foreignObject’ element, the back to front stacking order for a stacking context created by an SVG element is:
	//
	//   1. the background and borders of the element forming the stacking context, if any
	//   2. descendants, in tree order
	//
	// Since the ‘foreignObject’ element creates a "fixed position containing block" in CSS terms, the normative rules for the stacking order of the stacking context created by ‘foreignObject’ elements are the rules in Appendix E of CSS 2.1.
	//
	// Individual graphics elements are treated as if they are a non-isolated group, the components (fill, stroke, etc) that make up a graphic element (See Painting shapes and text) being members of that group.
	//
	// Grouping elements, such as the ‘g’ element (see container elements ) create a compositing group. Similarly, a ‘use’ element creates a compositing group for its shadow content. The Compositing and Blending specification normatively describes how to render compositing groups. In SVG, effects may be applied to a group. For example, opacity, filters or masking. These effects are applied to the rendered result of the group immediately before any transforms on the group are applied, which are applied immediately before the group is blended and composited with the group backdrop.
	//
	// Thus, rendering a compositing group follows the following steps:
	// If the group is isolated:
	//
	//   1. The initial backdrop is set to a new buffer initialised with rgba(0,0,0,0)
	//   2. The contents of the group that are graphics elements or ‘g’ elements are rendered in order, onto the initial backdrop
	//   3. filters and other effects that modify the group canvas are applied
	//      To provide for high quality rendering, filter primitives and other bitmap effects must be applied in the operating coordinate space.
	//   4. Group transforms are applied
	//   5. The group canvas is blended and composited with the group backdrop
	//
	// else (the group is not isolated):
	//
	//   1. The initial backdrop is set to the group backdrop
	//   2. The contents of the group that are graphics elements or ‘g’ elements are rendered in order, onto the initial backdrop. The group transforms are applied to each element as they are rendered.
	//

	// Collect elements by ID.
	//
	// TODO: duplicate IDs
	r := renderer{
		elements: map[string]Element{},
		fonts:    map[string]*fontFamily{},
	}
	walk(svg, func(e Element) {
		if id := e.id(); id != "" {
			r.elements[id] = e
		}
	})

	root := &Grouping{
		ElementAttributes: ElementAttributes{
			Color:      &Color{Value: color.Black},
			Fill:       &Paint{Color: color.Black},
			Stroke:     &Paint{Color: color.Transparent},
			FontFamily: &FontFamily{Values: []string{"sans-serif"}},
			FontSize:   &LengthPercentageNumberIdent{LengthPercentageNumber: LengthPercentageNumber{Number: 12}},
			FontStyle:  "normal",
			FontWeight: &NumberIdent{Ident: "normal"},
		},
	}

	r.push(root, float64(ctx.Width()), float64(ctx.Height()))

	return r.renderCompositingGroup(ctx, false, svg.Children)
}

type element struct {
	Element

	width, height float64
}

func (r *renderer) push(e Element, width, height float64) {
	r.stack = append(r.stack, &element{Element: e, width: width, height: height})
}

func (r *renderer) pop() {
	r.stack = r.stack[:len(r.stack)-1]
}

func (r *renderer) top() *element {
	return r.stack[len(r.stack)-1]
}

func (r *renderer) width() float64 {
	return r.top().width
}

func (r *renderer) height() float64 {
	return r.top().height
}

func (r *renderer) diag() float64 {
	w, h := r.width(), r.height()
	return math.Sqrt(w*w+h*h) / math.Sqrt2
}

type renderer struct {
	elements map[string]Element
	fonts    map[string]*fontFamily
	stack    []*element
}

func (r *renderer) computeNumberPercentage(parent float64, np *NumberPercentage) float64 {
	if np == nil {
		return parent
	}
	if np.Percentage != 0 {
		return np.Percentage * parent
	}
	return np.Number
}

func (r *renderer) computeLengthPercentage(parent float64, lp LengthPercentage) float64 {
	if lp.Percentage != 0 {
		return lp.Percentage * parent
	}
	if lp.Length.Units != "" {
		panic(errors.New("NYI: units"))
	}
	return lp.Length.Value
}

func (r *renderer) computeLengthPercentageNumber(parent float64, lp LengthPercentageNumber) float64 {
	if lp.Number != 0 {
		return lp.Number
	}
	return r.computeLengthPercentage(parent, lp.LengthPercentage)
}

func (r *renderer) computeBoxLengthPercentage(parent, auto float64, blp *BoxLengthPercentage) float64 {
	if blp == nil || blp.Value == "auto" {
		return auto
	}

	if blp.Value == "" {
		return r.computeLengthPercentage(parent, blp.LengthPercentage)
	}

	// TODO: validate value?

	return parent
}

func (r *renderer) computePattern(e Element, patternOpacity float64) (gg.Pattern, error) {
	switch e := e.(type) {
	case *LinearGradient:
		x1, y1 := r.computeLengthPercentage(r.width(), e.X1), r.computeLengthPercentage(r.height(), e.Y1)
		x2, y2 := r.computeLengthPercentage(r.width(), e.X2), r.computeLengthPercentage(r.height(), e.Y2)
		gradient := gg.NewLinearGradient(x1, y1, x2, y2)

		r.push(e, r.top().width, r.top().height)
		defer r.pop()

		dx, dy := x2-x1, y2-y1
		sz := math.Sqrt(dx*dx + dy*dy)

		for _, s := range e.Stops {
			offset := r.computeNumberPercentage(sz, &s.Offset)

			if s.Color.Value == nil {
				s.Color.Value = color.Black
			}
			sr, sg, sb, sa := s.Color.Value.RGBA()

			opacity := r.computeNumberPercentage(1.0, &s.Opacity)
			if sa == 0 {
				opacity = 0
			}

			gradient.AddColorStop(offset, color.NRGBA{
				R: byte(sr),
				G: byte(sg),
				B: byte(sb),
				A: byte(opacity * patternOpacity * 255),
			})
		}

		return gradient, nil
	case *RadialGradient:
		return nil, errors.New("NYI: radial gradients")
	case *Pattern:
		return nil, errors.New("NYI: pattern")
	default:
		return nil, errors.New("not a paint server element")
	}
}

func (r *renderer) computePaint(p *Paint, opacity float64) (gg.Pattern, error) {
	switch p.Context {
	case "context-fill":
		return nil, errors.New("NYI: context-fill")
	case "context-stroke":
		return nil, errors.New("NYI: context-stroke")
	}

	if p.URL != "" {
		if p.URL[0] != '#' {
			return nil, errors.New("NYI: non-fragment URLs")
		}
		id := p.URL[1:]
		e, ok := r.elements[id]
		if ok {
			if p, err := r.computePattern(e, opacity); err == nil {
				return p, nil
			}
		}
	}

	c := p.Color
	if c != nil {
		sr, sg, sb, sa := c.RGBA()
		if sa == 0 {
			opacity = 0
		}
		c = color.NRGBA{
			R: byte(sr),
			G: byte(sg),
			B: byte(sb),
			A: byte(opacity * 255),
		}
	} else {
		c = color.Transparent
	}
	return gg.NewSolidPattern(c), nil
}

func (r *renderer) setPaints(ctx *gg.Context) error {
	// Compute the fill (TODO: style)
	fillOpacity := r.computeNumberPercentage(1.0, r.getFillOpacity())
	fill, err := r.computePaint(r.getFill(), fillOpacity)
	if err != nil {
		return err
	}

	// Compute the stroke (TODO: style)
	strokeOpacity := r.computeNumberPercentage(1.0, r.getStrokeOpacity())
	stroke, err := r.computePaint(r.getStroke(), strokeOpacity)
	if err != nil {
		return err
	}

	// Handle line caps
	switch r.getStrokeLinecap() {
	case "butt":
		ctx.SetLineCap(gg.LineCapButt)
	case "round":
		ctx.SetLineCap(gg.LineCapRound)
	case "square":
		ctx.SetLineCap(gg.LineCapSquare)
	}

	// Handle line width
	strokeWidth := 1.0
	if sw := r.getStrokeWidth(); sw != nil {
		strokeWidth = r.computeLengthPercentage(r.diag(), *sw)
	}
	ctx.SetLineWidth(strokeWidth)

	ctx.SetFillStyle(fill)
	ctx.SetStrokeStyle(stroke)
	return nil
}

func (r *renderer) renderCompositingGroup(ctx *gg.Context, isolated bool, elements []any) error {
	for _, e := range elements {
		if err := r.renderElement(ctx, e.X); err != nil {
			return err
		}
	}
	return nil
}

func (r *renderer) renderElement(ctx *gg.Context, e Element) error {
	switch e := e.(type) {
	case *Grouping:
		return r.renderGrouping(ctx, e)
	case *Use:
		return r.renderUse(ctx, e)
	case *Switch:
		return r.renderSwitch(ctx, e)
	case *Path:
		return r.renderPath(ctx, e)
	case *Rect:
		return r.renderRect(ctx, e)
	case *Circle:
		return r.renderCircle(ctx, e)
	case *Ellipse:
		return r.renderEllipse(ctx, e)
	case *Line:
		return r.renderLine(ctx, e)
	case *Polyline:
		return r.renderPolyline(ctx, e)
	case *Polygon:
		return r.renderPolygon(ctx, e)
	case *Text:
		return r.renderText(ctx, e)
	case *TSpan:
		return r.renderTSpan(ctx, e)
	case *Image:
		return r.renderImage(ctx, e)
	case *ForeignObject:
		return r.renderForeignObject(ctx, e)
	case *Defs, *Marker, *Symbol, *LinearGradient, *RadialGradient, *Pattern:
		// Never rendered
		return nil
	default:
		panic(fmt.Errorf("unexpected element type %T", e))
	}
}

func (r *renderer) renderGrouping(ctx *gg.Context, e *Grouping) error {
	r.push(e, r.width(), r.height())
	defer r.pop()

	for _, e := range e.Children {
		if err := r.renderElement(ctx, e.X); err != nil {
			return err
		}
	}

	return nil
}

func (r *renderer) renderUse(ctx *gg.Context, e *Use) error {
	return errors.New("NYI: use")
}

func (r *renderer) renderSwitch(ctx *gg.Context, e *Switch) error {
	return errors.New("NYI: switch")
}

func (r *renderer) renderPath(ctx *gg.Context, e *Path) error {
	ctx.Push()
	defer ctx.Pop()

	r.push(e, r.width(), r.height())
	defer r.pop()

	r.setPaints(ctx)

	// TODO: path length

	p, _ := ctx.GetCurrentPoint()
	x, y := p.X, p.Y

	active, subpath := false, false
	ctx.ClearPath()
	for i, c := range e.D.Commands {
		switch c := c.(type) {
		case *MoveTo:
			if active {
				ctx.NewSubPath()
				subpath = true
			}
			active = true

			if !c.IsAbsolute {
				x, y = x+c.Points[0].X, y+c.Points[0].Y
			} else {
				x, y = c.Points[0].X, c.Points[0].Y
			}
			ctx.MoveTo(x, y)

			for _, p := range c.Points[1:] {
				if !c.IsAbsolute {
					x, y = x+p.X, y+p.Y
				} else {
					x, y = p.X, p.Y
				}
				ctx.LineTo(x, y)
			}
		case *ClosePath:
			ctx.ClosePath()
			if subpath {
				ctx.ClipPreserve()
			}
			active = false
		case *LineTo:
			for _, p := range c.Points {
				if !c.IsAbsolute {
					if !math.IsNaN(p.X) {
						x += p.X
					}
					if !math.IsNaN(p.Y) {
						y += p.Y
					}
				} else {
					if !math.IsNaN(p.X) {
						x = p.X
					}
					if !math.IsNaN(p.Y) {
						y = p.Y
					}
				}
				ctx.LineTo(x, y)
			}
		case *CubicBezier:

			// If the current point is (curx, cury) and the final control point of the
			// previous path segment is (oldx2, oldy2), then the reflected point (i.e.,
			// (newx1, newy1), the first control point of the current path segment) is:
			//
			// (newx1, newy1) = (curx - (oldx2 - curx), cury - (oldy2 - cury))
			//                = (2*curx - oldx2, 2*cury - oldy2)

			x1, y1, x2, y2 := 0.0, 0.0, 0.0, 0.0
			for _, p := range c.Coordinates {
				if c.IsSmooth {
					hasPreviousPoint := false
					if i > 0 {
						_, hasPreviousPoint = e.D.Commands[i-1].(*CubicBezier)
					}
					if !hasPreviousPoint {
						x1, y1 = x, y
					} else {
						x1, y1 = 2*x-x2, 2*y-y2
					}
				}

				if !c.IsAbsolute {
					if !c.IsSmooth {
						x1, y1 = x+p.X1, y+p.Y1
					}
					x2, y2 = x+p.X2, y+p.Y2
					x, y = x+p.X, y+p.Y
				} else {
					if !c.IsSmooth {
						x1, y1 = p.X1, p.Y1
					}
					x2, y2 = p.X2, p.Y2
					x, y = p.X, p.Y
				}
				ctx.CubicTo(x1, y1, x2, y2, x, y)
			}
		case *QuadraticBezier:
			return errors.New("NYI: quadratic bezier")
		case *EllipticalArc:
			return errors.New("NYI: elliptical arc")
		}
	}
	ctx.FillPreserve()
	ctx.StrokePreserve()
	ctx.ClearPath()

	return nil
}

func (r *renderer) renderRect(ctx *gg.Context, e *Rect) error {
	ctx.Push()
	defer ctx.Pop()

	x0, y0 := r.computeLengthPercentage(r.width(), e.X), r.computeLengthPercentage(r.height(), e.Y)
	w, h := r.computeBoxLengthPercentage(r.width(), r.width(), &e.Width), r.computeBoxLengthPercentage(r.height(), r.height(), &e.Height)

	cssRx := e.Rx
	if cssRx == nil {
		cssRx = e.Ry
	}
	rx := r.computeBoxLengthPercentage(r.width(), 0, e.Rx)
	ry := r.computeBoxLengthPercentage(r.height(), rx, e.Ry)

	x1, y1 := x0+rx, y0+ry
	x2, y2 := x0+w-rx, y0+h-ry
	x3, y3 := x0+w, y0+h

	r.push(e, w, h)
	defer r.pop()

	r.setPaints(ctx)

	ctx.ClearPath()
	ctx.MoveTo(x1, y0)
	ctx.LineTo(x2, y0)
	ctx.DrawEllipticalArc(x2, y1, rx, ry, gg.Radians(270), gg.Radians(360))
	ctx.LineTo(x3, y2)
	ctx.DrawEllipticalArc(x2, y2, rx, ry, gg.Radians(0), gg.Radians(90))
	ctx.LineTo(x1, y3)
	ctx.DrawEllipticalArc(x1, y2, rx, ry, gg.Radians(90), gg.Radians(180))
	ctx.LineTo(x0, y1)
	ctx.DrawEllipticalArc(x1, y1, rx, ry, gg.Radians(180), gg.Radians(270))
	ctx.FillPreserve()
	ctx.StrokePreserve()
	ctx.ClosePath()

	return nil
}

func (r *renderer) renderCircle(ctx *gg.Context, e *Circle) error {
	ctx.Push()
	defer ctx.Pop()

	cx, cy := r.computeLengthPercentage(r.width(), e.Cx), r.computeLengthPercentage(r.height(), e.Cy)
	rr := r.computeLengthPercentage(r.diag(), e.R)

	r.push(e, rr, rr)
	defer r.pop()

	r.setPaints(ctx)

	ctx.ClearPath()
	ctx.DrawCircle(cx, cy, rr)
	ctx.FillPreserve()
	ctx.StrokePreserve()

	return nil
}

func (r *renderer) renderEllipse(ctx *gg.Context, e *Ellipse) error {
	return errors.New("NYI: ellipse")
}

func (r *renderer) renderLine(ctx *gg.Context, e *Line) error {
	return errors.New("NYI: line")
}

func (r *renderer) renderPolyline(ctx *gg.Context, e *Polyline) error {
	return errors.New("NYI: polyline")
}

func (r *renderer) renderPolygon(ctx *gg.Context, e *Polygon) error {
	return errors.New("NYI: polygon")
}

func (r *renderer) renderText(ctx *gg.Context, e *Text) error {
	ctx.Push()
	defer ctx.Pop()

	cssStyle := r.getFontStyle()
	style := font.StyleNormal
	if cssStyle == "italic" {
		style = font.StyleItalic
	}

	cssWeight := r.getFontWeight()
	weight := font.WeightNormal
	switch cssWeight.Ident {
	case "":
		switch cssWeight.Number {
		case 100:
			weight = font.WeightThin
		case 200:
			weight = font.WeightExtraLight
		case 300:
			weight = font.WeightLight
		case 400:
			weight = font.WeightNormal
		case 500:
			weight = font.WeightMedium
		case 600:
			weight = font.WeightSemiBold
		case 700:
			weight = font.WeightBold
		case 800:
			weight = font.WeightExtraBold
		case 900:
			weight = font.WeightBlack
		default:
			weight = font.Weight(cssWeight.Number/100 - 4)
		}
	case "normal":
		weight = font.WeightNormal
	case "bold":
		weight = font.WeightBold
	case "bolder", "lighter":
		return errors.New("NYI: bolder/lighter")
	}

	cssSize := r.getFontSize()
	size := 0.0
	if cssSize.Ident != "" {
		return errors.New("NYI: named font sizes")
	}
	switch {
	case cssSize.Length.Value != 0:
		if cssSize.Length.Units != "pt" {
			return errors.New("NYI: non-pt font sizes")
		}
		size = cssSize.Length.Value
	case cssSize.Percentage != 0:
		return errors.New("NYI: percentage font size")
	default:
		size = cssSize.Number
	}

	// TODO: stretch

	ax, ay := 0.0, 0.0
	switch r.getTextAnchor() {
	case "middle":
		ax, ay = 0.5, 0
	case "end":
		ax, ay = 1.0, 0
	}

	r.push(e, r.width(), r.height())
	defer r.pop()

	r.setPaints(ctx)
	ctx.ClearPath()

	fontFamily, err := r.resolveFontFamily(r.getFontFamily())
	if err != nil {
		return err
	}
	face, err := fontFamily.newFace(weight, style, size, font.HintingNone)
	if err != nil {
		return err
	}
	ctx.SetFontFace(face)

	if len(e.X.Values) != 1 || len(e.Y.Values) != 1 {
		return errors.New("NYI: x/y lists; dx/dy")
	}

	x, y := r.computeLengthPercentageNumber(r.width(), e.X.Values[0]), r.computeLengthPercentageNumber(r.height(), e.Y.Values[0])

	if len(e.Rotate.Values) != 0 {
		return errors.New("NYI: rotate")
	}

	ctx.DrawStringAnchored(e.Value, x, y, ax, ay)
	return nil
}

func (r *renderer) renderTSpan(ctx *gg.Context, e *TSpan) error {
	return errors.New("NYI: tspan")
}

func (r *renderer) renderImage(ctx *gg.Context, e *Image) error {
	return errors.New("NYI: image")
}

func (r *renderer) renderForeignObject(ctx *gg.Context, e *ForeignObject) error {
	return errors.New("NYI: foreignObject")
}
