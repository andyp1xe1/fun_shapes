package main

import (
	"fmt"
	gen "fun_shapes/server/generator"
	net "fun_shapes/server/network"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/fogleman/gg"
)

// TODO Opts chan. image recv post req. multiple alg instances
func main() {
	optCh := make(chan gen.Options)
	net.StartServer(optCh)
	//o := &gen.Opts
	o := <-optCh
	close(optCh)

	defer func() {
		//Remove the file after processing
		err := os.Remove("./img_test/testfile.png")
		if err != nil {
			fmt.Println("Error removing file:", err)
		}
	}()

	img, _ := gg.LoadImage(o.InPath)
	c := gen.NewCanvas(img)

	numMonteSamples := int(float64(c.Dx*c.Dy) * o.MonteDensity)

	// Random col. example
	// i := rand.Intn(len(c.Palette))
	// s.Color = c.Palette[i]
	procSteps := []procConf{
		{
			NumShapes:      o.NumMonteShapes,
			PopulationSize: o.PopulationSize,
			Ctx:            c.Dc,
			Fn: func() gen.Shape {
				shape := gen.NewShape(c, o.ShapeType)
				shape.SetColFrom(c)
				c.EvalScoreMonte(shape, numMonteSamples)
				return shape
			},
		}, {
			NumShapes:      o.NumSolidShapes,
			PopulationSize: o.PopulationSize,
			Ctx:            c.Dc,
			Fn: func() gen.Shape {
				shape := gen.NewShape(c, o.ShapeType)
				shape.SetColFrom(c)
				c.EvalScore(shape)
				return shape
			},
		}, {
			NumShapes:      o.NumOpaqueShapes,
			PopulationSize: o.PopulationSize,
			Ctx:            c.Dc,
			Fn: func() gen.Shape {
				shape := gen.NewShape(c, o.ShapeType)
				col := shape.SetColFrom(c)
				shape.SetCol(gen.ColAddOpacity(col, 0.8))
				c.EvalScore(shape)
				return shape
			},
		},
	}

	for _, conf := range procSteps {
		processor(conf)
	}
	c.Dc.SavePNG(genOutPath(o.InPath))

}

func processor(conf procConf) {
	for i := 0; i < conf.NumShapes; i++ {
		wg := sync.WaitGroup{}
		shapes := make([]gen.Shape, conf.PopulationSize)
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
		bestShape.Draw(conf.Ctx)
		//println(bestShape.GetScore())
		net.UpdateCurrentImg(conf.Ctx)

	}
}

type ProcFunc func() gen.Shape
type procConf struct {
	NumShapes      int
	PopulationSize int
	Fn             ProcFunc
	Ctx            *gg.Context
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
