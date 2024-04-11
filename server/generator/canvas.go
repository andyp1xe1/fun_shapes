package generator

import (
	"github.com/ericpauley/go-quantize/quantize"
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
	dx := img.Bounds().Dx() / 5
	dy := img.Bounds().Dy() / 5
	img = resize.Thumbnail(uint(dx), uint(dy), img, resize.Lanczos2)
	p := quantizeImage(img, Opts.NumColors)
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

//func (c *Canvas) DrawSafe(s Shape) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	s.Draw(c.Dc)
//}

func (c *Canvas) EvalScore(s Shape, f1 bool) float64 {
	ctx := gg.NewContext(c.Dc.Width(), c.Dc.Height())
	ctx.DrawImage(c.Dc.Image(), 0, 0)
	s.Draw(ctx)
	var score float64
	if f1 {
		score = old_fitness_monte(c.Img, ctx, Opts.NumSamples)
	} else {
		score = fitness(c.Img, ctx)
	}
	s.SetScore(score)
	return s.GetScore()
}

func addOpacity(cols []color.Color, opacity float64) []color.Color {
	p := make([]color.Color, len(cols))
	for i, c := range cols {
		r, g, b, _ := c.RGBA()
		p[i] = color.RGBA{
			uint8(r >> 8),
			uint8(g >> 8),
			uint8(b >> 8),
			uint8(opacity * 255),
		}
	}
	return p
}

func quantizeImage(img image.Image, numColors int) color.Palette {
	palette := make(color.Palette, 0, numColors)
	q := quantize.MedianCutQuantizer{}
	return addOpacity(q.Quantize(palette, img), 0.8)
}

func averageColor(p color.Palette) color.Color {
	var totalR, totalG, totalB uint32
	for _, c := range p {
		r, g, b, _ := c.RGBA()
		totalR += r
		totalG += g
		totalB += b
	}
	n := uint32(len(p))
	avgR := totalR / n
	avgG := totalG / n
	avgB := totalB / n

	return color.RGBA{
		uint8(avgR >> 8),
		uint8(avgG >> 8),
		uint8(avgB >> 8),
		255, // Full opacity
	}
}
