package svg

func walk(svg *SVG, visitor func(e Element)) {
	// Document-order walk.
	walkElements(svg.Children, visitor)
}

func walkElements(es []any, visitor func(e Element)) {
	for _, e := range es {
		walkElement(e.X, visitor)
	}
}

func walkElement(e Element, visitor func(e Element)) {
	visitor(e)

	switch e := e.(type) {
	case *Grouping:
		walkElements(e.Children, visitor)
	case *Defs:
		walkElements(e.Children, visitor)
	case *Symbol:
		walkElements(e.Children, visitor)
	case *Switch:
		walkElements(e.Children, visitor)
	case *Marker:
		walkElements(e.Children, visitor)
	case *Pattern:
		walkElements(e.Children, visitor)
	case *Path:
		walkElements(e.Children, visitor)
	case *Rect:
		walkElements(e.Children, visitor)
	case *Circle:
		walkElements(e.Children, visitor)
	case *Ellipse:
		walkElements(e.Children, visitor)
	case *Line:
		walkElements(e.Children, visitor)
	case *Polyline:
		walkElements(e.Children, visitor)
	case *Polygon:
		walkElements(e.Children, visitor)
	}
}
