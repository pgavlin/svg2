# svg2

[![PkgGoDev](https://pkg.go.dev/badge/github.com/pgavlin/svg2)](https://pkg.go.dev/github.com/pgavlin/svg2)
[![codecov](https://codecov.io/gh/pgavlin/svg2/branch/master/graph/badge.svg)](https://codecov.io/gh/pgavlin/svg2)
[![Go Report Card](https://goreportcard.com/badge/github.com/pgavlin/svg2)](https://goreportcard.com/report/github.com/pgavlin/svg2)
[![Test](https://github.com/pgavlin/svg2/workflows/Test/badge.svg)](https://github.com/pgavlin/svg2/actions?query=workflow%3ATest)

A pure Go SVG renderer built on top of [gg](https://pkg.go.dev/github.com/fogleman/gg).

Note that this is _very much_ a work in progress, and many features of SVG are
not implemented. This includes (but is not limited to):
- length units
- radial gradients
- patterns
- context-fill
- context-stroke
- non-fragment URLs
- use elements
- switch elements
- quadratic bezier curves
- elliptical arcs
- ellipse elements
- line elements
- polyline elements
- tspan elements
- image elements
- foreignObject elements
- transforms

Which is really to say that pretty much the only SVG elements that _are_ supported are
paths, groups, linear gradients, and text.
