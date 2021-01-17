package swf

import (
	"github.com/dexter3k/dash/com"
)

type FillStyle interface{}

type SolidFillStyle struct {
	Color com.RGBA
}

type GradientPoint struct {
	Ratio uint8
	Color com.RGBA
}

type GradientSpread int8
const (
	SpreadPad     GradientSpread = iota
	SpreadReflect
	SpreadRepeat
)

type GradientFillStyle struct {
	Matrix     com.Matrix
	Radial     bool
	Focal      bool
	Spread     GradientSpread
	LinearRGB  bool
	FocalPoint com.I8_8
	Points     []*GradientPoint
}

type BitmapFillStyle struct {
	Id      uint16
	Matrix  com.Matrix
	Clipped bool
	Smooth  bool
}


type LineStyle interface{}

type LineStyleSimple struct {
	Width uint16
	Color com.RGBA
}

type LineStyleFancy struct {
	// not implemented
}


type Edge struct {
	Target  com.Vector2
	Control com.Vector2
}

type Shape struct {
	Position   com.Vector2
	FillStyle0 FillStyle
	FillStyle1 FillStyle
	LineStyle  LineStyle
	Edges      []*Edge
}

type DefineShape struct {
	Id     uint16
	Bounds com.Rect
	Shapes []*Shape
}
