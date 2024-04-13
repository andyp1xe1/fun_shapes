package generator

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

type Triangle struct {
	X1, Y1, X2, Y2, X3, Y3 float64
	Score                  float64
	Color                  color.Color
}

func NewTrig(dx, dy int) *Triangle {
	return &Triangle{
		X1: rand.Float64() * float64(dx),
		Y1: rand.Float64() * float64(dy),
		X2: rand.Float64() * float64(dx),
		Y2: rand.Float64() * float64(dy),
		X3: rand.Float64() * float64(dx),
		Y3: rand.Float64() * float64(dy),
	}
}

func (t *Triangle) Mutate() *Triangle {
	return t
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
