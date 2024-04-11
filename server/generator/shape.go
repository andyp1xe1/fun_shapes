package generator

import "github.com/fogleman/gg"

type Shape interface {
	Generate(c *Canvas)
	Draw(dc *gg.Context)
	//Mutate()
	SetScore(score float64)
	GetScore() float64
}

func NewShape(shapeType string, c *Canvas) Shape {
	var s Shape
	switch shapeType {
	case "Triangle":
		s = &Triangle{}
	case "Rectangle":
		s = &Rectangle{}
	default:
		panic("Unknown shape type")
	}
	s.Generate(c)
	return s
}
