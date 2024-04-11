package generator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	X, Y, W, H, Th float64
	Err            float64
	Score          float64
	Color          color.Color
}

func (r *Rectangle) GenerateByNum(c *Canvas, i int) {
	// Set random position within the bounds of the image
	r.X = rand.Float64() * float64(c.Dx)
	r.Y = rand.Float64() * float64(c.Dy)

	maxWidth := float64(c.Dx) * 0.6
	maxHeight := float64(c.Dy) * 0.6
	minWidth := 10.0
	minHeight := 5.0

	scale := rand.Float64() + float64(Opts.NumShapes-i)/float64(Opts.NumShapes)
	r.W = (maxWidth-minWidth)*scale + minWidth
	r.H = (maxHeight-minHeight)*scale + minHeight

	r.Th = rand.Float64() * 2 * math.Pi

	col := rand.Intn(len(c.Palette))
	r.Color = c.Palette[col]
}

func (r *Rectangle) Generate(c *Canvas) {
	r.X = rand.Float64() * float64(c.Dx)
	r.Y = rand.Float64() * float64(c.Dy)
	r.W = rand.Float64()*100 + 2
	r.H = rand.Float64()*150 + 2
	r.Th = rand.Float64() * 2 * math.Pi

	i := rand.Intn(len(c.Palette))
	r.Color = c.Palette[i]

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

func (r *Rectangle) SetScore(score float64) {
	r.Score = score
}

func (r *Rectangle) GetScore() float64 {
	return r.Score
}
