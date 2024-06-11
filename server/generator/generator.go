package generator

import (
	"mime/multipart"
)

type Options struct {
	Id              string
	InImage         multipart.File
	FrameCh         chan []byte
	NumSolidShapes  int
	NumOpaqueShapes int
	NumMonteShapes  int
	MonteDensity    float64
	PopulationSize  int
	ShapeType       ShapeType
	//InPath          string
}

//var Opts = Options{
//	NumSolidShapes:  200, // max 500
//	NumOpaqueShapes: 30,  // max 100
//	NumMonteShapes:  80,  // max 50
//	MonteDensity:    0.4, //0 - 1
//	PopulationSize:  150,
//
//	ShapeType: "4",
//}
