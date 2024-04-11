package main

import (
	gen "fun_shapes/server/generator"
	net "fun_shapes/server/network"
	"github.com/fogleman/gg"
	"sort"
	"sync"
)

func main() {
	o := &gen.Opts
	net.StartServer(o)

	img, _ := gg.LoadImage(o.InPath)
	c := gen.NewCanvas(img)

	o.NumSamples = c.Dx * c.Dy * 4 / 100

	for i := 0; i < o.NumShapes; i++ {
		var wg sync.WaitGroup
		shapes := make([]gen.Shape, o.PopulationSize)

		for j := 0; j < o.PopulationSize; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				shape := gen.NewShape("Rectangle", c)
				c.EvalScore(shape, i < o.PopulationSize/3)
				shapes[idx] = shape
			}(j)
		}
		wg.Wait()

		sort.Slice(shapes, func(i, j int) bool {
			return shapes[i].GetScore() < shapes[j].GetScore()
		})

		bestShape := shapes[0]
		bestShape.Draw(c.Dc)
		println(bestShape.GetScore())
		net.UpdateCurrentImg(c)

	}
	c.Dc.SavePNG(o.OutPath)
}
