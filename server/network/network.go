package network

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"net/http"
	//"sync"
)

var imgChan = make(chan []byte, 10)

func StartServer() {
	//http.HandleFunc("/conf", setOpts)
	http.HandleFunc("/frame", getCurrentImage)
	http.HandleFunc("/", serveViewer)
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
		fmt.Println("Skipping image update; channel is full")
	}
}

func serveViewer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
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
