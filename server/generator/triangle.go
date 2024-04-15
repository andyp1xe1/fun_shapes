package generator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Triangle struct {
	X1, Y1, X2, Y2, X3, Y3 float64
	Score                  float64
	Color                  color.Color
}

func NewTrig(dx, dy int) *Triangle {
	for {

		t := &Triangle{
			X1: rand.Float64() * float64(dx),
			Y1: rand.Float64() * float64(dy),
			X2: rand.Float64() * float64(dx),
			Y2: rand.Float64() * float64(dy),
			X3: rand.Float64() * float64(dx),
			Y3: rand.Float64() * float64(dy),
		}
		l1 := math.Sqrt((t.X2-t.X1)*(t.X2-t.X1) + (t.Y2-t.Y1)*(t.Y2-t.Y1))
		l2 := math.Sqrt((t.X3-t.X2)*(t.X3-t.X2) + (t.Y3-t.Y2)*(t.Y3-t.Y2))
		l3 := math.Sqrt((t.X1-t.X3)*(t.X1-t.X3) + (t.Y1-t.Y3)*(t.Y1-t.Y3))

		angleA := math.Acos((l2*l2+l3*l3-l1*l1)/(2*l2*l3)) * 180 / math.Pi
		angleB := math.Acos((l1*l1+l3*l3-l2*l2)/(2*l1*l3)) * 180 / math.Pi
		angleC := 180 - angleA - angleB
		if trianglevalidation(l1, l2, l3, angleA, angleB, angleC) {
			return t
		}
	}
}

func (t *Triangle) Mutate() *Triangle {
	return t
}

func trianglevalidation(l1, l2, l3, angleA, angleB, angleC float64) bool {
	if angleA < 15 || angleB < 15 || angleC < 15 {
		return false
	}
	if l1 > 200 || l2 > 200 || l3 > 200 {
		return false
	}
	return true
}

func (t *Triangle) Draw(dc *gg.Context) {
	dc.Push()
	dc.MoveTo(t.X1, t.Y1)
	dc.LineTo(t.X2, t.Y2)
	dc.LineTo(t.X3, t.Y3)
	dc.ClosePath()
	dc.SetColor(t.Color)
	dc.Fill()
	dc.Pop()
}

func (t *Triangle) SetColFrom(c *Canvas) color.Color {
	cx := (t.X1 + t.X2 + t.X3) / 3
	cy := (t.Y1 + t.Y2 + t.Y3) / 3
	t.Color = c.ColorAt(cx, cy)
	return t.Color
}

func (t *Triangle) SetScore(score float64) {
	t.Score = score
}

func (t *Triangle) GetScore() float64 {
	return t.Score
}

func (t *Triangle) SetCol(c color.Color) {
	t.Color = c
}

func (t *Triangle) GetCol() color.Color {
	return t.Color
}
