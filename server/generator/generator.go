package generator

import (
	"mime/multipart"
)

type FrameChan chan string

func NewFrameChan() FrameChan {
	return make(chan string, 128)
}

type FrameSelector interface {
	selectCh(id string) FrameChan
	registerCh(id string, frameChan FrameChan)
	FrameChan(id string) FrameChan
}

type Options struct {
	//Id              string
	InImage         multipart.File
	FrameCh         FrameChan
	NumSolidShapes  int
	NumOpaqueShapes int
	NumMonteShapes  int
	MonteDensity    float64
	PopulationSize  int
	ShapeType       ShapeType
	//InPath          string
}
