package generator

import (
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	//"sync"
)

type Canvas struct {
	Dx, Dy  int
	Palette color.Palette
	Dc      *gg.Context
	Img     image.Image
	//mu      sync.Mutex
}

func NewCanvas(img image.Image) *Canvas {
	dx := img.Bounds().Dx() / 4
	dy := img.Bounds().Dy() / 4
	img = resize.Thumbnail(uint(dx), uint(dy), img, resize.Lanczos2)
	p := quantizeImage(img, 256)
	dc := gg.NewContext(dx, dy)
	dc.SetColor(averageColor(p))
	dc.Clear()

	return &Canvas{
		Dx:      dx,
		Dy:      dy,
		Palette: p,
		Dc:      dc,
		Img:     img,
	}
}

func (c *Canvas) EvalScore(s Shape) float64 {
	ctx := gg.NewContext(c.Dc.Width(), c.Dc.Height())
	ctx.DrawImage(c.Dc.Image(), 0, 0)
	s.Draw(ctx)
	score := fitness(c.Img, ctx)
	s.SetScore(score)
	return s.GetScore()
}
func (c *Canvas) EvalScoreMonte(s Shape, monteSamples int) float64 {
	ctx := gg.NewContext(c.Dc.Width(), c.Dc.Height())
	ctx.DrawImage(c.Dc.Image(), 0, 0)
	s.Draw(ctx)
	var score float64
	if monteSamples > 0 {
		score = old_fitness_monte(c.Img, ctx, monteSamples)
	} else {
		score = fitness(c.Img, ctx)
	}
	s.SetScore(score)
	return s.GetScore()
}
