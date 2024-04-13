package generator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	X, Y, W, H, Th float64
	Score          float64
	Color          color.Color
}

func NewRect(dx, dy int) *Rectangle {
	r := &Rectangle{}
	r.X = rand.Float64() * float64(dx)
	r.Y = rand.Float64() * float64(dy)
	r.W = rand.Float64()*100 + 2
	r.H = rand.Float64()*150 + 2
	r.Th = rand.Float64() * 2 * math.Pi
	return r
}

func (r *Rectangle) Mutate() *Rectangle {
	// Adjust position randomly
	r.X += rand.Float64()*20 - 10
	r.Y += rand.Float64()*20 - 10

	r.W += rand.Float64()*20 - 10
	r.H += rand.Float64()*20 - 10

	// Adjust rotation angle randomly
	maxRotation := 0.2
	r.Th += rand.Float64()*maxRotation*2 - maxRotation
	return r
}
func (r *Rectangle) Draw(dc *gg.Context) {
	dc.Push()

	centerX := r.X + r.W/2
	centerY := r.Y + r.H/2

	dc.Translate(centerX, centerY)
	dc.RotateAbout(r.Th, 0, 0)
	dc.Translate(-r.W/2, -r.H/2)

	dc.SetColor(r.Color)
	dc.DrawRectangle(0, 0, r.W, r.H)
	dc.Fill()

	dc.Pop()
}

func (r *Rectangle) SetColFrom(c *Canvas) color.Color {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	r.Color = c.ColorAt(cx, cy)
	return r.Color
}

func (r *Rectangle) SetScore(score float64) {
	r.Score = score
}

func (r *Rectangle) GetScore() float64 {
	return r.Score

}
func (r *Rectangle) SetCol(c color.Color) {
	r.Color = c
}

func (r *Rectangle) GetCol() color.Color {
	return r.Color
}
