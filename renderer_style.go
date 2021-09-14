package svg

func (r *renderer) getAttr(get func(Element) bool) {
	for i := len(r.stack) - 1; i >= 0; i-- {
		if get(r.stack[i].Element) {
			return
		}
	}
}

func (r *renderer) getAlignmentBaseline() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().AlignmentBaseline; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getBaselineShift() *LengthPercentageIdent {
	var v *LengthPercentageIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().BaselineShift; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getClipPath() *ClipPath {
	var v *ClipPath
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ClipPath; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getClipRule() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ClipRule; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getColor() *Color {
	var v *Color
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Color; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getColorInterpolation() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ColorInterpolation; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getColorInterpolationFilters() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ColorInterpolationFilters; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getColorRendering() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ColorRendering; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getCursor() *Cursor {
	var v *Cursor
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Cursor; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getDirection() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Direction; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getDisplay() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Display; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getDominantBaseline() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().DominantBaseline; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFill() *Paint {
	var v *Paint
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Fill; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFillOpacity() *NumberPercentage {
	var v *NumberPercentage
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FillOpacity; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFillRule() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FillRule; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFilter() *FilterList {
	var v *FilterList
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Filter; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFloodColor() *Color {
	var v *Color
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FloodColor; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFloodOpacity() *NumberPercentage {
	var v *NumberPercentage
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FloodOpacity; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontFamily() *FontFamily {
	var v *FontFamily
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontFamily; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontSize() *LengthPercentageNumberIdent {
	var v *LengthPercentageNumberIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontSize; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontSizeAdjust() *NumberIdent {
	var v *NumberIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontSizeAdjust; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontStretch() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontStretch; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontStyle() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontStyle; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontVariant() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontVariant; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getFontWeight() *NumberIdent {
	var v *NumberIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().FontWeight; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getImageRendering() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ImageRendering; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getLetterSpacing() *LengthIdent {
	var v *LengthIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().LetterSpacing; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getLightingColor() *Color {
	var v *Color
	r.getAttr(func(e Element) bool {
		if i := e.attrs().LightingColor; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getMarkerEnd() *URLIdent {
	var v *URLIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().MarkerEnd; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getMarkerMid() *URLIdent {
	var v *URLIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().MarkerMid; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getMarkerStart() *URLIdent {
	var v *URLIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().MarkerStart; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getMask() *Mask {
	var v *Mask
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Mask; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getOpacity() *NumberIdent {
	var v *NumberIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Opacity; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getOverflow() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Overflow; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getPaintOrder() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().PaintOrder; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getPointerEvents() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().PointerEvents; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getShapeRendering() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().ShapeRendering; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStroke() *Paint {
	var v *Paint
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Stroke; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeDasharray() *DashArray {
	var v *DashArray
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeDasharray; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeDashoffset() *LengthPercentage {
	var v *LengthPercentage
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeDashoffset; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeLinecap() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeLinecap; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeLinejoin() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeLinejoin; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeMiterlimit() *float64 {
	var v *float64
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeMiterlimit; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeOpacity() *NumberPercentage {
	var v *NumberPercentage
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeOpacity; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getStrokeWidth() *LengthPercentage {
	var v *LengthPercentage
	r.getAttr(func(e Element) bool {
		if i := e.attrs().StrokeWidth; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getTextAnchor() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().TextAnchor; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getTextDecoration() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().TextDecoration; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getTextOverflow() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().TextOverflow; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getTextRendering() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().TextRendering; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getTransform() []Transform {
	var v []Transform
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Transform; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getUnicodeBidi() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().UnicodeBidi; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getVectorEffect() *VectorEffect {
	var v *VectorEffect
	r.getAttr(func(e Element) bool {
		if i := e.attrs().VectorEffect; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getVisibility() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().Visibility; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getWhiteSpace() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().WhiteSpace; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getWordSpacing() *LengthIdent {
	var v *LengthIdent
	r.getAttr(func(e Element) bool {
		if i := e.attrs().WordSpacing; i != nil {
			v = i
			return true
		}
		return false
	})
	return v
}

func (r *renderer) getWritingMode() Ident {
	var v Ident
	r.getAttr(func(e Element) bool {
		if i := e.attrs().WritingMode; i != "" {
			v = i
			return true
		}
		return false
	})
	return v
}
