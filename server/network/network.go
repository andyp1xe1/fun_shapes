package network

import (
	"bytes"
	"fun_shapes/server/generator"
	"io"
	"os"
	"strconv"

	//"strconv"
	//"sync"

	//"encoding/json"
	"fmt"
	"net/http"

	"github.com/fogleman/gg"
)

var imgChan = make(chan []byte, 10)

func StartServer(optChan chan generator.Options) {
	//http.HandleFunc("/conf", setOpts)
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/images", http.FileServer(http.Dir("."))) // Serve static files (e.g., HTML, CSS, JS)
	http.HandleFunc("/frame", getCurrentImage)
	http.HandleFunc("/", serveViewer)
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		submitHandler(w, r, optChan)
	})

	go func() {
		fmt.Println("Server started at http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
}

func getCurrentImage(w http.ResponseWriter, r *http.Request) {
	select {
	case img := <-imgChan:
		w.Header().Set("Content-Type", "image/png")
		w.Write(img)
	default:
		http.Error(w, "No image available", http.StatusNotFound)
	}
}

func UpdateCurrentImg(dc *gg.Context) {
	var buf bytes.Buffer
	err := dc.EncodePNG(&buf)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}
	select {
	case imgChan <- buf.Bytes():
	default:
		//fmt.Println("Skipping image update; channel is full")
	}
}

func serveViewer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form containing the file
	err := r.ParseMultipartForm(10 << 20) // 10 MB maximum
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from the form data
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// For demonstration, let's just copy the file to a new location
	out, err := os.Create("./img_test/testfile.png")
	if err != nil {
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to copy file", http.StatusInternalServerError)
		return
	}
}

// Modify submitHandler to accept and process the integer value sent via AJAX
func submitHandler(w http.ResponseWriter, r *http.Request, optChan chan generator.Options) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the integer value from the AJAX request
	shapeType := r.FormValue("shape")
	strSolidnr := r.FormValue("solidshapes")
	strOpaquenr := r.FormValue("opaqueshapes")
	strMontenr := r.FormValue("monteshapes")
	strMontedensity := r.FormValue("montedensity")

	solidnr, _ := strconv.Atoi(strSolidnr)
	opaquenr, _ := strconv.Atoi(strOpaquenr)
	montenr, _ := strconv.Atoi(strMontenr)
	montedensity, _ := strconv.ParseFloat(strMontedensity, 64)

	fmt.Printf("Received shape: %s", shapeType)
	fmt.Printf("Values: %d %d %d %f", solidnr, opaquenr, montenr, montedensity)
	//.. create Opts obj
	newOpts := generator.Options{
		//NumColors:       256,

		PopulationSize: 150,
		InPath:         "./img_test/testfile.png",
	}
	newOpts.NumSolidShapes = solidnr
	newOpts.NumOpaqueShapes = opaquenr
	newOpts.NumMonteShapes = montenr
	newOpts.MonteDensity = montedensity
	newOpts.ShapeType = generator.ShapeType(shapeType)
	optChan <- newOpts
}

//func setOpts(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var newOpts generator.Options
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&newOpts); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	generator.UpdateOptions(&newOpts)
//}
