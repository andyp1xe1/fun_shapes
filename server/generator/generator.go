package generator

type Options struct {
	NumSolidShapes  int
	NumOpaqueShapes int
	NumMonteShapes  int
	MonteDensity    float64
	PopulationSize  int
	InPath          string
	ShapeType       ShapeType
	//NumColors         int
}

var Opts = Options{
	//NumColors:       256,
	NumSolidShapes:  200, // max 500
	NumOpaqueShapes: 30,  // max 100
	NumMonteShapes:  80,  // max 50
	MonteDensity:    0.4, //0 - 1
	PopulationSize:  150,
	InPath:          "./img_test/in.png",
	ShapeType:       "4",
}

//func printoption()

//func UpdateOptions(newOpts *Options) {
//	currentOpts = *newOpts
//}
//
//func GetOptions() *Options {
//	return &currentOpts
//}
