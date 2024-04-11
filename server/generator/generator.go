package generator

type Options struct {
	NumColors, NumShapes, PopulationSize, NumSamples int
	InPath, OutPath                                  string
}

var Opts = Options{
	NumColors:      256,
	NumShapes:      600,
	PopulationSize: 150,
	NumSamples:     0,
	InPath:         "./in2.jpg",
	OutPath:        "out2.png",
}
