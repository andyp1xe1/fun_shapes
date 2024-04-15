package generator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Ellipse struct {
	X, Y, Rx, Ry, Th float64
	Score            float64
	Color            color.Color
}

func NewEllipse(dx, dy int) *Ellipse {
	for {
		r := &Ellipse{}
		r.X = rand.Float64() * float64(dx)
		r.Y = rand.Float64() * float64(dy)
		r.Rx = rand.Float64()*50 + 2
		r.Ry = rand.Float64()*75 + 2
		r.Th = rand.Float64() * 2 * math.Pi
		if r.Rx < 2*r.Ry && r.Ry < 2*r.Rx {
			return r
		}
	}
}

func (r *Ellipse) Mutate() *Ellipse {
	// Adjust position randomly
	r.X += rand.Float64()*20 - 10
	r.Y += rand.Float64()*20 - 10

	r.Rx += rand.Float64()*20 - 10
	r.Ry += rand.Float64()*20 - 10

	// Adjust rotation angle randomly
	maxRotation := 0.2
	r.Th += rand.Float64()*maxRotation*2 - maxRotation
	return r
}
func (r *Ellipse) Draw(dc *gg.Context) {
	dc.Push()
	dc.Translate(r.X, r.Y)
	//dc.RotateAbout(r.Th, 0, 0)
	//dc.Translate(-r.X, -r.Y)

	dc.SetColor(r.Color)
	dc.DrawEllipse(0, 0, r.Rx, r.Ry)
	dc.Fill()

	dc.Pop()
}

func (r *Ellipse) SetColFrom(c *Canvas) color.Color {
	cx := r.X
	cy := r.Y
	r.Color = c.ColorAt(cx, cy)
	return r.Color
}

func (r *Ellipse) SetScore(score float64) {
	r.Score = score
}

func (r *Ellipse) GetScore() float64 {
	return r.Score

}
func (r *Ellipse) SetCol(c color.Color) {
	r.Color = c
}

func (r *Ellipse) GetCol() color.Color {
	return r.Color
}
