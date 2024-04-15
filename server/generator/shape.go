package generator

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
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
	case "Ellipse":
		s = NewEllipse(c.Dx, c.Dy)
	case "All":
		vari := rand.Float64()
		if vari <= 0.3 {
			s = NewTrig(c.Dx, c.Dy)
		}
		if vari < 0.6 && vari > 0.3 {
			s = NewRect(c.Dx, c.Dy)
		}
		if vari >= 0.6 {
			s = NewEllipse(c.Dx, c.Dy)
		}
	default:
		panic("Unknown shape type")

	}
	return s

}
