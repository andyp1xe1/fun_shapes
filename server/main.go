package main

import (
	//"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/ericpauley/go-quantize/quantize"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type Options struct {
	NumColors, NumShapes, PopulationSize int
	numSamples                           int
}

type Canvas struct { // TODO fix this crap of code (data structuring)
	Dx, Dy  int
	Palette color.Palette
	dc      *gg.Context
	img     image.Image
}

type Shape interface {
	Generate()
	Draw()
	Mutate()
}

type Rectangle struct {
	x, y, w, h, th float64
	score          float64
	color          color.Color
}

func (r *Rectangle) Generate(c Canvas, o Options, i int) {
	// Set random position within the bounds of the image
	r.x = rand.Float64() * float64(c.Dx)
	r.y = rand.Float64() * float64(c.Dy)

	// Define practical proportions for width and height based on the shape number
	maxWidth := float64(c.Dx) * 0.6  // Set maximum width to 60% of the image width
	maxHeight := float64(c.Dy) * 0.6 // Set maximum height to 60% of the image height
	minWidth := 10.0                 // Minimum width
	minHeight := 5.0                 // Minimum height

	// Adjust the width and height based on the shape number
	scale := rand.Float64() + float64(o.NumShapes-i)/float64(o.NumShapes)
	r.w = (maxWidth-minWidth)*scale + minWidth
	r.h = (maxHeight-minHeight)*scale + minHeight

	// Set random rotation angle
	r.th = rand.Float64() * 2 * math.Pi

	// Randomly select a color from the provided palette
	col := rand.Intn(len(c.Palette))
	r.color = c.Palette[col]
}

func (r *Rectangle) Mutate() *Rectangle {
	// Adjust position randomly
	r.x += rand.Float64()*20 - 10 // Randomly adjust x within ±10 pixels
	r.y += rand.Float64()*20 - 10 // Randomly adjust y within ±10 pixels

	r.w += rand.Float64()*20 - 10 // Randomly adjust width within ±10 pixels
	r.h += rand.Float64()*20 - 10 // Randomly adjust height within ±10 pixels

	// Adjust rotation angle randomly
	maxRotation := 0.2                                 // Maximum rotation in radians
	r.th += rand.Float64()*maxRotation*2 - maxRotation // Randomly adjust rotation within ±maxRotation
	return r
}

func (r *Rectangle) Score(o Options, c Canvas) float64 {
	ctx := gg.NewContext(c.dc.Width(), c.dc.Height())
	ctx.DrawImage(c.dc.Image(), 0, 0)
	r.Draw(ctx)
	r.score = fitness2(c.img, ctx, o.numSamples)
	return r.score
}

func (r *Rectangle) Draw(dc *gg.Context) {
	// Save the current transformation matrix
	dc.Push()

	// Translate to the center of the rectangle
	centerX := r.x + r.w/2
	centerY := r.y + r.h/2
	dc.Translate(centerX, centerY)

	// Rotate around the center of the rectangle
	dc.RotateAbout(r.th, 0, 0)

	// Translate back to the top-left corner of the rectangle
	dc.Translate(-r.w/2, -r.h/2)

	// Set the color and draw the rectangle
	dc.SetColor(r.color)
	dc.DrawRectangle(0, 0, r.w, r.h)
	dc.Fill()

	// Restore the previous transformation matrix
	dc.Pop()
}

func main() {
	_img, _ := gg.LoadImage("./in.png")
	o := Options{
		NumColors:      100,
		NumShapes:      200,
		PopulationSize: 60,
	}

	c := Canvas{
		Dx: _img.Bounds().Dx() / 4,
		Dy: _img.Bounds().Dy() / 4,
	}

	c.img = resize.Thumbnail(uint(c.Dx), uint(c.Dy), _img, resize.Lanczos2)
	o.numSamples = c.Dx * c.Dy * 7 / 100

	c.Palette = quantizeImage(c.img, o.NumColors)
	c.dc = gg.NewContext(c.Dx, c.Dy)

	avgColor := averageColor(c.Palette)

	c.dc.SetColor(avgColor)
	c.dc.Clear()

	for i := range o.NumShapes {
		var bestRect *Rectangle
		//toDraw := false
		for range o.PopulationSize {
			ctx := gg.NewContext(c.dc.Width(), c.dc.Height())
			ctx.DrawImage(c.dc.Image(), 0, 0)

			rect := &Rectangle{}
			rect.Generate(c, o, i)
			rect.Score(o, c)

			rect.score = fitness2(c.img, ctx, o.numSamples)
			if bestRect == nil || rect.score < bestRect.score {
				mutRect := &Rectangle{
					x:     rect.x,
					y:     rect.y,
					w:     rect.w,
					th:    rect.th,
					color: rect.color,
				}
				mutRect.Mutate()
				mutRect.Score(o, c)
				if mutRect.score < rect.score {
					bestRect = mutRect
				} else {
					bestRect = rect
				}
			}
		}
		//if toDraw {
		bestRect.Draw(c.dc)
		println(bestRect.score)
		//i++
		//}

	}
	c.dc.SavePNG("out00.png")
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

// Quantize the image to a palette
func quantizeImage(img image.Image, numColors int) color.Palette {
	palette := make(color.Palette, 0, numColors)
	q := quantize.MedianCutQuantizer{}
	return addOpacity(q.Quantize(palette, img), 0.8)
}

func addOpacity(o []color.Color, opacity float64) []color.Color {
	p := make([]color.Color, len(o))
	for i, c := range o {
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

func fitness2(originalImg image.Image, ctx *gg.Context, numSamples int) float64 {
	generatedImg := ctx.Image()

	// Calculate Mean Squared Error (MSE)
	var total uint64
	bounds := originalImg.Bounds()
	for range numSamples {
		x := bounds.Min.X + rand.Intn(bounds.Dx())
		y := bounds.Min.Y + rand.Intn(bounds.Dy())

		originalColor := originalImg.At(x, y)
		generatedColor := generatedImg.At(x, y)

		r1, g1, b1, _ := originalColor.RGBA()
		r2, g2, b2, _ := generatedColor.RGBA()
		r := r1 - r2
		g := g1 - g2
		b := b1 - b2
		//a := a1 - a2
		total += uint64(r*r + g*g + b*b) // + a*a
	}

	//return math.Sqrt(float64(total)/float64(numSamples*3)) / 225
	return math.Log(float64(total))
}

func fitness(originalImg image.Image, generatedImg *gg.Context) uint64 {
	// Convert generated image to image.Image
	generatedImage := generatedImg.Image()

	// Calculate Mean Squared Error (MSE)
	var mse uint64
	bounds := originalImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := originalImg.At(x, y)
			generatedColor := generatedImage.At(x, y)

			r1, g1, b1, _ := originalColor.RGBA()
			r2, g2, b2, _ := generatedColor.RGBA()

			r := r1 - r2
			g := g1 - g2
			b := b1 - b2
			//a := a1 - a2
			mse += uint64(r*r + g*g + b*b)
		}
	}

	return mse
}
