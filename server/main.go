package main

import (
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
	Dx, Dy                               int
	Palette                              color.Palette
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

func (r *Rectangle) Generate(o Options) {
	r.x = rand.Float64() * float64(o.Dx)
	r.y = rand.Float64() * float64(o.Dy)
	r.w = rand.Float64()*100 + 20
	r.h = rand.Float64()*100 + 20
	r.th = rand.Float64() * 2 * math.Pi

	i := rand.Intn(len(o.Palette))
	r.color = o.Palette[i]

}

func (r *Rectangle) Mutate() {
	//TODO
}

func (r *Rectangle) Draw(dc *gg.Context) {
	dc.Push()
	dc.SetColor(r.color)
	dc.Rotate(r.th)
	dc.DrawRectangle(r.x, r.y, r.w, r.h)
	dc.Fill()
	dc.Pop()
}

func main() {
	img, _ := gg.LoadImage("./in.png")
	options := Options{
		NumColors:      255,
		NumShapes:      200,
		PopulationSize: 30,
		Dx:             img.Bounds().Dx() / 4,
		Dy:             img.Bounds().Dy() / 4,
	}
	img = resize.Thumbnail(uint(options.Dx), uint(options.Dy), img, resize.Lanczos2)

	options.Palette = quantizeImage(img, options.NumColors)
	dc := gg.NewContext(options.Dx, options.Dy)

	avgColor := averageColor(options.Palette)

	dc.SetColor(avgColor)
	dc.Clear()

	for range options.NumShapes {
		//bestRect := &Rectangle{score: -math.MaxFloat64}
		var bestRect *Rectangle

		for range options.PopulationSize {
			ctx := gg.NewContext(dc.Width(), dc.Height())
			ctx.DrawImage(dc.Image(), 0, 0)
			rect := &Rectangle{}
			rect.Generate(options)
			rect.Draw(ctx)

			rect.score = fitness(img, ctx)
			if bestRect == nil || rect.score > bestRect.score {
				bestRect = rect
				//dc = ctx
			}
		}

		bestRect.Draw(dc)

		println(bestRect.score)
	}
	dc.SavePNG("out.png")
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

// vvvv  !!!!!!!!!!!!!!!!!!! CHATGPT CODE !!!!!!!!!!!!!!!!!!!   vvvv
func fitness(originalImg image.Image, generatedImg *gg.Context) float64 {
	// Convert generated image to image.Image
	generatedImage := generatedImg.Image()

	// Calculate Mean Squared Error (MSE)
	var mse float64
	bounds := originalImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := originalImg.At(x, y)
			generatedColor := generatedImage.At(x, y)

			r1, g1, b1, _ := originalColor.RGBA()
			r2, g2, b2, _ := generatedColor.RGBA()

			mse += math.Pow(float64(r1-r2), 2) + math.Pow(float64(g1-g2), 2) + math.Pow(float64(b1-b2), 2)
		}
	}

	// Normalize MSE
	numPixels := float64(bounds.Dx() * bounds.Dy())
	maxPossibleMSE := 3 * math.Pow(255, 2) // Maximum pixel value is 255
	normalizedMSE := mse / (numPixels * maxPossibleMSE)

	return normalizedMSE
} //   ^^^^  !!!!!!!!!!!!!!!!!!! CHATGPT CODE !!!!!!!!!!!!!!!!!!!   ^^^^
