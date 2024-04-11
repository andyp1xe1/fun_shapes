package generator

import (
	//"fmt"
	"github.com/fogleman/gg"
	"image"
	"math"
	"math/rand"
)

func old_fitness_monte(originalImg image.Image, ctx *gg.Context, numSamples int) float64 {
	generatedImg := ctx.Image()
	var total float64
	bounds := originalImg.Bounds()
	for range numSamples {
		x := bounds.Min.X + rand.Intn(bounds.Dx())
		y := bounds.Min.Y + rand.Intn(bounds.Dy())

		originalColor := originalImg.At(x, y)
		generatedColor := generatedImg.At(x, y)

		r1, g1, b1, _ := originalColor.RGBA()
		r2, g2, b2, _ := generatedColor.RGBA()

		total += math.Pow(float64(r2-r1), 2) + math.Pow(float64(g2-g1), 2) + math.Pow(float64(b2-b1), 2)

	}
	return total
	//return math.Sqrt(float64(total)/float64(numSamples*3)) / 225

}

func old_fitness(originalImg image.Image, generatedImg *gg.Context) float64 {
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

			println(r1, g1, b1)
			println(r2, g2, b2)

			println("dr: r1 - r2: ", r1-r2)
			println("dr to flotat: ", float64(r1-r2))
			println("dr^2: ", math.Pow(float64(r1-r2), 2))
			mse += math.Pow(float64(r1-r2), 2) + math.Pow(float64(g1-g2), 2) + math.Pow(float64(b1-b2), 2)
		}
	}

	// Normalize MSE
	numPixels := float64(bounds.Dx() * bounds.Dy())
	maxPossibleMSE := 3 * math.Pow(255, 2) // Maximum pixel value is 255
	normalizedMSE := mse / (numPixels * maxPossibleMSE)

	return normalizedMSE
}

func fitness_monte(originalImg image.Image, ctx *gg.Context, numSamples int) float64 {
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

	return math.Sqrt(float64(total)/float64(numSamples*3)) / 225
	//return math.Log(float64(total))
}

func fitness(originalImg image.Image, generatedImg *gg.Context) float64 {
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

	return math.Sqrt(float64(mse)/float64(generatedImg.Width()*generatedImg.Height()*3)) / 225
}
