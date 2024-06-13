package main

import (
	"fmt"
	gen "fun_shapes/server/generator"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type FrameSelector interface {
	FrameChan(addr string) <-chan []byte
	Register(addr string, frameChan chan []byte)
}

type FrameStore struct {
	mutex sync.Mutex
	chans map[string](chan []byte)
}

func newFrameStore() *FrameStore {
	return &FrameStore{
		mutex: sync.Mutex{},
		chans: make(map[string](chan []byte)),
	}
}

func (m *FrameStore) FrameChan(addr string) <-chan []byte {
	m.mutex.Lock()
	ch := m.chans[addr]
	m.mutex.Unlock()
	fmt.Printf("querying user %s\n", addr)
	return ch
}

func (m *FrameStore) Register(addr string, frameChan chan []byte) {
	m.mutex.Lock()
	m.chans[addr] = frameChan
	m.mutex.Unlock()
	fmt.Printf("registered user %s\n", addr)
}

func frameHandler(selector FrameSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case frame := <-selector.FrameChan(r.RemoteAddr):
			w.Header().Set("Content-Type", "image/png")
			w.Write(frame)
		default:
			http.Error(w, "No image available", http.StatusNotFound)
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func submitHandler(optsChan chan gen.Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		shapeType := r.FormValue("shape")
		strSolidnr := r.FormValue("solidshapes")
		strOpaquenr := r.FormValue("opaqueshapes")
		strMontenr := r.FormValue("monteshapes")
		strMontedensity := r.FormValue("montedensity")

		solidnr, _ := strconv.Atoi(strSolidnr)
		opaquenr, _ := strconv.Atoi(strOpaquenr)
		montenr, _ := strconv.Atoi(strMontenr)
		montedensity, _ := strconv.ParseFloat(strMontedensity, 64)

		image, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received shape: %s\n", shapeType)
		fmt.Printf("Values: %d %d %d %f\n", solidnr, opaquenr, montenr, montedensity)

		optsChan <- gen.Options{
			Id:              r.RemoteAddr,
			FrameCh:         make(chan []byte, 64),
			PopulationSize:  150,
			NumSolidShapes:  solidnr,
			NumOpaqueShapes: opaquenr,
			NumMonteShapes:  montenr,
			MonteDensity:    montedensity,
			ShapeType:       gen.ShapeType(shapeType),
			InImage:         image,
		}
	}
}

// TODO make canvas implement `save` method
func main() {
	frames := newFrameStore()
	optsChan := make(chan gen.Options, 100)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/submit", submitHandler(optsChan))
	http.HandleFunc("/frame", frameHandler(frames))

	go func() {
		fmt.Println("Server started at http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	wg := sync.WaitGroup{}
	for opts := range optsChan {
		go func() {
			wg.Add(1)
			handleOpts(&opts)
			//c.SavePNG() //TODO
			wg.Done()
		}()
		frames.Register(opts.Id, opts.FrameCh)
	}
	wg.Wait()
	close(optsChan)
}

func handleOpts(o *gen.Options) {
	c := gen.NewCanvas(o.InImage)

	// Random col. example
	// i := rand.Intn(len(c.Palette))
	// s.Color = c.Palette[i]
	procSteps := []gen.ProcConf{
		{
			NumShapes:      o.NumMonteShapes,
			PopulationSize: o.PopulationSize,
			Fn: func() gen.Shape {
				shape := gen.NewShape(c, o.ShapeType)
				shape.SetColFrom(c)
				c.EvalScoreMonte(shape, int(float64(c.Dx*c.Dy)*o.MonteDensity))
				return shape
			},
		}, {
			NumShapes:      o.NumSolidShapes,
			PopulationSize: o.PopulationSize,
			Fn: func() gen.Shape {
				shape := gen.NewShape(c, o.ShapeType)
				shape.SetColFrom(c)
				c.EvalScore(shape)
				return shape
			},
		}, {
			NumShapes:      o.NumOpaqueShapes,
			PopulationSize: o.PopulationSize,
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
		c.Process(conf, o.FrameCh)
		err := c.SaveSVGsToFile("output.txt")
		if err != nil {
			fmt.Println("Error saving SVG data:", err)
		}
	}
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
