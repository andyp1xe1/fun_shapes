package generator

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Triangle struct {
	X1, Y1, X2, Y2, X3, Y3, Th float64
	Err                        float64
	Score                      float64
	Color                      color.Color
}

func (t *Triangle) GenerateByNum(c *Canvas, i int) {
	// Your implementation for generating triangle by number of shapes
}

func (t *Triangle) Generate(c *Canvas) {
	// Generate random coordinates for triangle vertices within canvas bounds
	t.X1 = rand.Float64() * float64(c.Dx)
	t.Y1 = rand.Float64() * float64(c.Dy)
	t.X2 = rand.Float64() * float64(c.Dx)
	t.Y2 = rand.Float64() * float64(c.Dy)
	t.X3 = rand.Float64() * float64(c.Dx)
	t.Y3 = rand.Float64() * float64(c.Dy)

	// Generate random rotation angle for the triangle
	t.Th = rand.Float64() * 2 * math.Pi

	// Generate random color for the triangle from canvas palette
	col := rand.Intn(len(c.Palette))
	t.Color = c.Palette[col]
}

func (t *Triangle) Mutate() *Triangle {
	return t
}

func (t *Triangle) Draw(dc *gg.Context) {
	dc.Push()

	// Move to the first vertex of the triangle
	dc.MoveTo(t.X1, t.Y1)

	// Draw lines to the other vertices
	dc.LineTo(t.X2, t.Y2)
	dc.LineTo(t.X3, t.Y3)

	// Close the path
	dc.ClosePath()

	// Translate and rotate the triangle
	dc.Translate(t.X1, t.Y1)
	dc.RotateAbout(t.Th, 0, 0)
	dc.Translate(-t.X1, -t.Y1)

	// Fill the triangle with the specified color
	dc.SetColor(t.Color)
	dc.Fill()

	dc.Pop()
}

func (t *Triangle) SetScore(score float64) {
	t.Score = score
}

func (t *Triangle) GetScore() float64 {
	return t.Score
}
