package main

import (
	gen "fun_shapes/server/generator"

	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// type FrameChan chan []byte

type Server struct {
	mutex sync.Mutex
	wg    sync.WaitGroup
	chans map[string]gen.FrameChan
}

func newServer() *Server {
	return &Server{
		mutex: sync.Mutex{},
		wg:    sync.WaitGroup{},
		chans: make(map[string]gen.FrameChan),
	}
}

func (m *Server) selectCh(id string) gen.FrameChan {
	m.mutex.Lock()
	ch := m.chans[id]
	m.mutex.Unlock()
	fmt.Printf("querying user %s\n", id)
	return ch
}

func (m *Server) registerCh(id string, ch gen.FrameChan) {
	m.mutex.Lock()
	m.chans[id] = ch
	m.mutex.Unlock()
	fmt.Printf("registered user %s\n", id)
}

func (s *Server) FrameChan(id string) (ch gen.FrameChan) {
	ch = s.selectCh(id)
	if ch == nil {
		ch = gen.NewFrameChan()
		s.registerCh(id, ch)
	}
	return
}

func (s *Server) handleWs(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	n, err := ws.Read(buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println("read error: ", err)
		return
	}
	uuid := string(buf[:n])
	ch := s.FrameChan(uuid)
	fmt.Println("ws remote uuid: ", uuid)
	for frame := range ch {
		ws.Write([]byte(frame))
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (s *Server) submitHandler() http.HandlerFunc {
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

		id := r.FormValue("uuid")
		if id == "" {
			id = uuid.NewString()
			w.Write([]byte(id))
		}

		opts := gen.Options{
			FrameCh:         s.FrameChan(id),
			PopulationSize:  150,
			NumSolidShapes:  solidnr,
			NumOpaqueShapes: opaquenr,
			NumMonteShapes:  montenr,
			MonteDensity:    montedensity,
			ShapeType:       gen.ShapeType(shapeType),
			InImage:         image,
		}

		s.wg.Add(1)
		go handleOpts(&opts)
	}

}

func main() {
	server := newServer()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/submit", server.submitHandler())
	//http.HandleFunc("/frame", frameHandler(frames))
	http.Handle("/ws", websocket.Handler(server.handleWs))

	func() {
		fmt.Println("Server started at http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	server.wg.Wait()
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
	}
	err := c.SaveSVGsToFile("output.txt")
	if err != nil {
		fmt.Println("Error saving SVG data:", err)
	}
}
