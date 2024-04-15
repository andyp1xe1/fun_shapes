package generator

type Options struct {
	NumSolidShapes  int
	NumOpaqueShapes int
	NumMonteShapes  int
	MonteDensity    float64
	PopulationSize  int
	InPath          string
	ShapeType       string
	//NumColors         int
}

var Opts = Options{
	//NumColors:       256,
	NumSolidShapes:  200,
	NumOpaqueShapes: 30,
	NumMonteShapes:  80,
	MonteDensity:    0.4,
	PopulationSize:  150,
	InPath:          "./img_test/in.png",
	ShapeType:       "All",
}

//func UpdateOptions(newOpts *Options) {
//	currentOpts = *newOpts
//}
//
//func GetOptions() *Options {
//	return &currentOpts
//}
