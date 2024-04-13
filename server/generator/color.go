package generator

import (
	"github.com/ericpauley/go-quantize/quantize"
	"image"
	"image/color"
)

func (c *Canvas) ColorAt(x, y float64) color.Color {
	bounds := c.Img.Bounds()
	ix, iy := int(x), int(y)
	if ix < bounds.Min.X || ix >= bounds.Max.X || iy < bounds.Min.Y || iy >= bounds.Max.Y {
		return color.RGBA{0, 0, 0, 255}
	}
	return c.Img.At(ix, iy)
	//idx := c.Palette.Index(c.Img.At(ix, iy))
	//return c.Palette[idx]
}

func ColAddOpacity(col color.Color, opacity float64) color.Color {
	r, g, b, _ := col.RGBA()
	new_c := color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(opacity * 255),
	}
	return new_c
}

func PaletteAddOpacity(cols []color.Color, opacity float64) []color.Color {
	p := make([]color.Color, len(cols))
	for i, c := range cols {
		p[i] = ColAddOpacity(c, opacity)
	}
	return p
}

func quantizeImage(img image.Image, numColors int) color.Palette {
	palette := make(color.Palette, 0, numColors)
	q := quantize.MedianCutQuantizer{}
	return PaletteAddOpacity(q.Quantize(palette, img), 0.8)
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

// Unused
func findMissingColors(original, generated color.Palette) color.Palette {
	missing := make(color.Palette, 0)
	found := false

	for _, oColor := range original {
		found = false
		for _, gColor := range generated {
			if oColor == gColor {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, oColor)
		}
	}
	return PaletteAddOpacity(missing, 0.8)
}
