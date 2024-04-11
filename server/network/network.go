package network

import (
	"bytes"
	"fmt"
	"fun_shapes/server/generator"
	"net/http"
	"sync"
)

var (
	currentImage []byte
	mutex        sync.Mutex
)

func StartServer(o *generator.Options) {
	http.HandleFunc("/frame", getCurrentImage)
	http.HandleFunc("/", serveViewer) // Add handler to serve the viewer HTML page
	go func() {
		fmt.Println("Server started at http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
}

func getCurrentImage(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	if currentImage == nil {
		http.Error(w, "No image available", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(currentImage)
}

func serveViewer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func UpdateCurrentImg(c *generator.Canvas) {
	var buf bytes.Buffer
	err := c.Dc.EncodePNG(&buf)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	currentImage = buf.Bytes()
}
