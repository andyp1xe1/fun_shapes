package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"mime/multipart"
	"sort"
	"sync"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type Canvas struct {
	Dx, Dy  int
	Palette color.Palette
	Dc      *gg.Context
	Img     image.Image
	//mu      sync.Mutex
}

type ProcFunc func() Shape
type ProcConf struct {
	NumShapes      int
	PopulationSize int
	Fn             ProcFunc
	Ctx            *gg.Context
	frameChan      chan []byte
}

func NewCanvas(file multipart.File) *Canvas {
	img, _, _ := image.Decode(file)
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

func (c *Canvas) Process(conf ProcConf) {
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
		//println(bestShape.GetScore())
		c.Draw(bestShape)

		select {
		case conf.frameChan <- c.toBytes():
		default:
			//fmt.Println("Skipping image update; channel is full")
		}
	}
}

func (c *Canvas) toBytes() []byte {
	var buf bytes.Buffer
	err := c.Dc.EncodePNG(&buf)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return nil
	}
	return buf.Bytes()
}
