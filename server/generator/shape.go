package generator

import (
	"github.com/fogleman/gg"
	"image/color"
)

type Shape interface {
	Draw(dc *gg.Context)
	//Mutate()
	SetScore(score float64)
	GetScore() float64
	SetCol(c color.Color)
	GetCol() color.Color
	SetColFrom(c *Canvas) color.Color
}

//type ShapeGenConf struct {
//	ShapeType string
//	C         *Canvas
//}

func NewShape(c *Canvas, ShapeType string) Shape {
	var s Shape
	switch ShapeType {
	case "Triangle":
		s = NewTrig(c.Dx, c.Dy)
	case "Rectangle":
		s = NewRect(c.Dx, c.Dy)
	default:
		panic("Unknown shape type")

	}
	return s

}
