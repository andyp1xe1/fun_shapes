package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type Canvas struct {
	Dx, Dy  int
	Palette color.Palette
	Dc      *gg.Context
	Img     image.Image
	//mu      sync.Mutex
	SVGs []string
}

type ProcFunc func() Shape
type ProcConf struct {
	NumShapes      int
	PopulationSize int
	Fn             ProcFunc
	Ctx            *gg.Context
}

func NewCanvas(file multipart.File) *Canvas {
	img, _, _ := image.Decode(file)
	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	dx := int(300 * ratio)
	dy := 300
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
		SVGs:    []string{},
	}
}
func (c *Canvas) Draw(s Shape) {
	s.Draw(c.Dc)
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

func (c *Canvas) Process(conf ProcConf, ch FrameChan) {
	for i := 0; i < conf.NumShapes; i++ {
		wg := sync.WaitGroup{}
		shapes := make([]Shape, conf.PopulationSize)
		for j := 0; j < conf.PopulationSize; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				shapes[idx] = conf.Fn()
			}(j)
		}
		wg.Wait()

		sort.Slice(shapes, func(i, j int) bool {
			return shapes[i].GetScore() < shapes[j].GetScore()
		})

		bestShape := shapes[0]
		c.Draw(bestShape)
		svg := bestShape.ToSVG()
		c.SVGs = append(c.SVGs, svg)

		ch <- svg

		//select {
		//case ch <- svg:
		//default:
		//	fmt.Println("Skipping image update; channel is full")
		//}

	}
}

func (c *Canvas) ToBytes() []byte {
	var buf bytes.Buffer
	err := c.Dc.EncodePNG(&buf)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return nil
	}
	return buf.Bytes()
}

func (c *Canvas) Save() {
	//TODO
}

func genOutPath(originalPath string) string {
	filename := filepath.Base(originalPath)
	filenameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
	timestamp := time.Now().Format("20060102_150405")
	newFilename := fmt.Sprintf("%s_%s.png", filenameWithoutExt, timestamp)
	out := filepath.Join("./img_res", newFilename)
	println("Saving to", out)
	return out
}

func (c *Canvas) SaveSVGsToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, svg := range c.SVGs {
		//fmt.Printf("SVG :\n%s\n\n", svg)
		_, err := file.WriteString(svg + "\n")
		if err != nil {
			return err
		}
	}
	//_, err = file.Write(c.ToBytes())
	//if err != nil {
	//	return err
	//}
	fmt.Println("SVG data saved to", "./svg")
	return nil
}
