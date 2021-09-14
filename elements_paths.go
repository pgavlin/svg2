package svg

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
)

// PathData represents SVG path data.
type PathData struct {
	Commands []PathCommand
}

func (d *PathData) UnmarshalText(text []byte) error {
	commands, err := ParsePathCommands(string(text))
	if err != nil {
		return err
	}
	d.Commands = commands
	return nil
}

// Path represents an SVG `path` element.
type Path struct {
	ElementAttributes

	XMLName xml.Name `xml:"path"`

	PathLength float64  `xml:"pathLength,attr"`
	D          PathData `xml:"d,attr"`

	Children []any `xml:",any"`
}

func (Path) isElement() {}

// PathCommand represents an SVG path command.
type PathCommand interface {
	isPathCommand()
}

type Point struct {
	X float64
	Y float64
}

// MoveTo represents an SVG `moveto` command.
type MoveTo struct {
	IsAbsolute bool

	Points []Point
}

func (*MoveTo) isPathCommand() {}

// ClosePath represents an SVG `closepath` command.
type ClosePath struct{}

func (*ClosePath) isPathCommand() {}

// LineTo represents an SVG `lineto` command.
type LineTo struct {
	IsAbsolute bool

	Points []Point
}

func (*LineTo) isPathCommand() {}

// CubicBezierCoordinates represents a set of cubic Bezier curve coordinates.
type CubicBezierCoordinates struct {
	Point

	X1 float64
	Y1 float64
	X2 float64
	Y2 float64
}

// CubicBezier represents an SVG cubic Bezier curve command.
type CubicBezier struct {
	IsAbsolute bool
	IsSmooth   bool

	Coordinates []CubicBezierCoordinates
}

func (*CubicBezier) isPathCommand() {}

// QuadraticBezierCoordinates represents a set of cubic Bezier curve coordinates.
type QuadraticBezierCoordinates struct {
	Point

	X1 float64
	Y1 float64
}

// QuadraticBezier represents an SVG cubic Bezier curve command.
type QuadraticBezier struct {
	IsAbsolute bool

	Coordinates []QuadraticBezierCoordinates
}

func (*QuadraticBezier) isPathCommand() {}

// EllipticalArcCoordinates represents a set of elliptical arc coordinates.
type EllipticalArcCoordinates struct {
	Point

	Rx            float64
	Ry            float64
	XAxisRotation float64
	LargeArc      bool
	Sweep         bool
}

// EllipticalArc represents an SVG elliptical arc command.
type EllipticalArc struct {
	IsAbsolute bool

	Coordinates []EllipticalArcCoordinates
}

func (*EllipticalArc) isPathCommand() {}

// ParsePathCommands parses a sequence of SVG path commands according to the SVG path
// grammar (reproduced below).
//
// svg_path::= wsp* moveto? (moveto drawto_command*)?
//
// drawto_command::=
//     moveto
//     | closepath
//     | lineto
//     | horizontal_lineto
//     | vertical_lineto
//     | curveto
//     | smooth_curveto
//     | quadratic_bezier_curveto
//     | smooth_quadratic_bezier_curveto
//     | elliptical_arc
//
// moveto::=
//     ( "M" | "m" ) wsp* coordinate_pair_sequence
//
// closepath::=
//     ("Z" | "z")
//
// lineto::=
//     ("L"|"l") wsp* coordinate_pair_sequence
//
// horizontal_lineto::=
//     ("H"|"h") wsp* coordinate_sequence
//
// vertical_lineto::=
//     ("V"|"v") wsp* coordinate_sequence
//
// curveto::=
//     ("C"|"c") wsp* curveto_coordinate_sequence
//
// curveto_coordinate_sequence::=
//     coordinate_pair_triplet
//     | (coordinate_pair_triplet comma_wsp? curveto_coordinate_sequence)
//
// smooth_curveto::=
//     ("S"|"s") wsp* smooth_curveto_coordinate_sequence
//
// smooth_curveto_coordinate_sequence::=
//     coordinate_pair_double
//     | (coordinate_pair_double comma_wsp? smooth_curveto_coordinate_sequence)
//
// quadratic_bezier_curveto::=
//     ("Q"|"q") wsp* quadratic_bezier_curveto_coordinate_sequence
//
// quadratic_bezier_curveto_coordinate_sequence::=
//     coordinate_pair_double
//     | (coordinate_pair_double comma_wsp? quadratic_bezier_curveto_coordinate_sequence)
//
// smooth_quadratic_bezier_curveto::=
//     ("T"|"t") wsp* coordinate_pair_sequence
//
// elliptical_arc::=
//     ( "A" | "a" ) wsp* elliptical_arc_argument_sequence
//
// elliptical_arc_argument_sequence::=
//     elliptical_arc_argument
//     | (elliptical_arc_argument comma_wsp? elliptical_arc_argument_sequence)
//
// elliptical_arc_argument::=
//     number comma_wsp? number comma_wsp? number comma_wsp
//     flag comma_wsp? flag comma_wsp? coordinate_pair
//
// coordinate_pair_double::=
//     coordinate_pair comma_wsp? coordinate_pair
//
// coordinate_pair_triplet::=
//     coordinate_pair comma_wsp? coordinate_pair comma_wsp? coordinate_pair
//
// coordinate_pair_sequence::=
//     coordinate_pair | (coordinate_pair comma_wsp? coordinate_pair_sequence)
//
// coordinate_sequence::=
//     coordinate | (coordinate comma_wsp? coordinate_sequence)
//
// coordinate_pair::= coordinate comma_wsp? coordinate
//
// coordinate::= sign? number
//
// sign::= "+"|"-"
// number ::= ([0-9])+
// flag::=("0"|"1")
// comma_wsp::=(wsp+ ","? wsp*) | ("," wsp*)
// wsp ::= (#x9 | #x20 | #xA | #xC | #xD)
func ParsePathCommands(commands string) ([]PathCommand, error) {
	r := bufio.NewReader(strings.NewReader(commands))

	if err := skipWhitespace(r); err != nil {
		return nil, err
	}

	var pathCommands []PathCommand
	for {
		next, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch next {
		case 'Z', 'z':
			pathCommands = append(pathCommands, &ClosePath{})
			continue
		}

		if err = skipWhitespace(r); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var command PathCommand
		switch next {
		case 'M', 'm':
			command, err = parseMoveTo(r, next == 'M')
		case 'L', 'l':
			command, err = parseLineTo(r, next == 'L', false, false)
		case 'H', 'h':
			command, err = parseLineTo(r, next == 'H', true, false)
		case 'V', 'v':
			command, err = parseLineTo(r, next == 'V', false, true)
		case 'C', 'c':
			command, err = parseCubicBezier(r, next == 'C', false)
		case 'S', 's':
			command, err = parseCubicBezier(r, next == 'S', true)
		case 'Q', 'q':
			command, err = parseQuadraticBezier(r, next == 'Q', false)
		case 'T', 't':
			command, err = parseQuadraticBezier(r, next == 'T', true)
		case 'A', 'a':
			command, err = parseEllipticalArc(r, next == 'A')
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		pathCommands = append(pathCommands, command)
	}

	return pathCommands, nil
}

func isWhitespace(b byte) bool {
	switch b {
	case 0x09, 0x0a, 0x0c, 0x0d, 0x20:
		return true
	}
	return false
}

func skipWhitespace(r *bufio.Reader) error {
	for {
		next, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !isWhitespace(next) {
			return r.UnreadByte()
		}
	}
}

func parseSign(r *bufio.Reader) (float64, error) {
	next, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	switch next {
	case '+':
		return 1.0, nil
	case '-':
		return -1.0, nil
	}
	return 1.0, r.UnreadByte()
}

func startsCoordinate(b byte) bool {
	return b == '-' || b == '+' || b >= '0' && b <= '9'
}

func parseCoordinate(r *bufio.Reader) (float64, error) {
	v, err := parseSign(r)
	if err != nil {
		return 0, err
	}

	var b strings.Builder
	for {
		next, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		if (next < '0' || next > '9') && next != '.' {
			if err = r.UnreadByte(); err != nil {
				return 0, err
			}
			break
		}
		b.WriteByte(next)
	}

	f, err := strconv.ParseFloat(b.String(), 64)
	if err != nil {
		return 0, err
	}
	return v * f, nil
}

func parseOptionalComma(r *bufio.Reader) (bool, error) {
	next, err := r.ReadByte()
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	if isWhitespace(next) {
		if err = skipWhitespace(r); err != nil {
			return false, err
		}
		next, err = r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
	}

	switch {
	case next == ',':
		return true, skipWhitespace(r)
	case startsCoordinate(next):
		return true, r.UnreadByte()
	}
	return false, r.UnreadByte()
}

func parseCoordinatePair(r *bufio.Reader) (Point, error) {
	x, err := parseCoordinate(r)
	if err != nil {
		return Point{}, err
	}

	if _, err = parseOptionalComma(r); err != nil {
		return Point{}, err
	}

	y, err := parseCoordinate(r)
	if err != nil {
		return Point{}, err
	}
	return Point{X: x, Y: y}, nil
}

func parseCoordinateSequence(r *bufio.Reader) ([]float64, error) {
	var coords []float64

	for {
		c, err := parseCoordinate(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, c)

		more, err := parseOptionalComma(r)
		if err == io.EOF || !more {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return coords, nil
}

func parseCoordinatePairSequence(r *bufio.Reader) ([]Point, error) {
	var coords []Point

	for {
		p, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, p)

		more, err := parseOptionalComma(r)
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
	}

	return coords, nil
}

func parseCoordinatePairTripletSequence(r *bufio.Reader) ([]Point, error) {
	var coords []Point

	for {
		a, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, a)

		if _, err = parseOptionalComma(r); err != nil {
			return nil, err
		}

		b, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, b)

		if _, err = parseOptionalComma(r); err != nil {
			return nil, err
		}

		c, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, c)

		more, err := parseOptionalComma(r)
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
	}

	return coords, nil
}

func parseCoordinatePairDoubleSequence(r *bufio.Reader) ([]Point, error) {
	var coords []Point

	for {
		a, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, a)

		if _, err = parseOptionalComma(r); err != nil {
			return nil, err
		}

		b, err := parseCoordinatePair(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, b)

		more, err := parseOptionalComma(r)
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
	}

	return coords, nil
}

func parseEllipticalArcArgument(r *bufio.Reader) (EllipticalArcCoordinates, error) {
	rp, err := parseCoordinatePair(r)
	if err != nil {
		return EllipticalArcCoordinates{}, err
	}

	if _, err = parseOptionalComma(r); err != nil {
		return EllipticalArcCoordinates{}, err
	}

	rot, err := parseCoordinate(r)
	if err != nil {
		return EllipticalArcCoordinates{}, err
	}

	if _, err = parseOptionalComma(r); err != nil {
		return EllipticalArcCoordinates{}, err
	}

	next, err := r.ReadByte()
	if err != nil {
		return EllipticalArcCoordinates{}, err
	}
	if next != '0' && next != '1' {
		return EllipticalArcCoordinates{}, errors.New("expected a flag")
	}
	largeArc := next == '1'

	if _, err = parseOptionalComma(r); err != nil {
		return EllipticalArcCoordinates{}, err
	}

	next, err = r.ReadByte()
	if err != nil {
		return EllipticalArcCoordinates{}, err
	}
	if next != '0' && next != '1' {
		return EllipticalArcCoordinates{}, errors.New("expected a flag")
	}
	sweep := next == '1'

	if _, err = parseOptionalComma(r); err != nil {
		return EllipticalArcCoordinates{}, err
	}

	point, err := parseCoordinatePair(r)
	if err != nil {
		return EllipticalArcCoordinates{}, err
	}

	return EllipticalArcCoordinates{
		Point:         point,
		Rx:            rp.X,
		Ry:            rp.Y,
		XAxisRotation: rot,
		LargeArc:      largeArc,
		Sweep:         sweep,
	}, nil
}

func parseMoveTo(r *bufio.Reader, isAbsolute bool) (*MoveTo, error) {
	points, err := parseCoordinatePairSequence(r)
	if err != nil {
		return nil, err
	}
	return &MoveTo{IsAbsolute: isAbsolute, Points: points}, nil
}

func parseLineTo(r *bufio.Reader, isAbsolute, isHoriz, isVert bool) (*LineTo, error) {
	switch {
	case isHoriz:
		xs, err := parseCoordinateSequence(r)
		if err != nil {
			return nil, err
		}
		points := make([]Point, len(xs))
		for i, x := range xs {
			points[i] = Point{X: x, Y: math.NaN()}
		}
		return &LineTo{IsAbsolute: isAbsolute, Points: points}, nil
	case isVert:
		ys, err := parseCoordinateSequence(r)
		if err != nil {
			return nil, err
		}
		points := make([]Point, len(ys))
		for i, y := range ys {
			points[i] = Point{X: math.NaN(), Y: y}
		}
		return &LineTo{IsAbsolute: isAbsolute, Points: points}, nil
	}

	points, err := parseCoordinatePairSequence(r)
	if err != nil {
		return nil, err
	}
	return &LineTo{IsAbsolute: isAbsolute, Points: points}, nil
}

func parseCubicBezier(r *bufio.Reader, isAbsolute, isSmooth bool) (*CubicBezier, error) {
	var coords []CubicBezierCoordinates

	if !isSmooth {
		points, err := parseCoordinatePairTripletSequence(r)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(points); i += 3 {
			c1, c2, c := points[i], points[i+1], points[i+2]
			coords = append(coords, CubicBezierCoordinates{
				Point: c,
				X1:    c1.X,
				Y1:    c1.Y,
				X2:    c2.X,
				Y2:    c2.Y,
			})
		}
	} else {
		points, err := parseCoordinatePairDoubleSequence(r)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(points); i += 2 {
			c2, c := points[i], points[i+1]
			coords = append(coords, CubicBezierCoordinates{
				Point: c,
				X2:    c2.X,
				Y2:    c2.Y,
			})
		}
	}

	return &CubicBezier{IsAbsolute: isAbsolute, IsSmooth: isSmooth, Coordinates: coords}, nil
}

func parseQuadraticBezier(r *bufio.Reader, isAbsolute, isSmooth bool) (*QuadraticBezier, error) {
	if isSmooth {
		return nil, errors.New("NYI: smooth quadratic bezier curve")
	}

	var coords []QuadraticBezierCoordinates

	points, err := parseCoordinatePairDoubleSequence(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(points); i += 2 {
		c1, c := points[i], points[i+1]
		coords = append(coords, QuadraticBezierCoordinates{
			Point: c,
			X1:    c1.X,
			Y1:    c1.Y,
		})
	}

	return &QuadraticBezier{IsAbsolute: isAbsolute, Coordinates: coords}, nil
}

func parseEllipticalArc(r *bufio.Reader, isAbsolute bool) (*EllipticalArc, error) {
	var coords []EllipticalArcCoordinates

	for {
		c, err := parseEllipticalArcArgument(r)
		if err != nil {
			return nil, err
		}
		coords = append(coords, c)

		more, err := parseOptionalComma(r)
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
	}

	return &EllipticalArc{IsAbsolute: isAbsolute, Coordinates: coords}, nil
}
