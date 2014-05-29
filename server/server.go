// A simple face detector server - analyzes images posted to it
// and responds with the rectangles containing matching objects
// (in json).
//
// The server can be started using:
//     go run server/server.go -haar=data/haarcascade_frontalface_alt.xml
//
// To test it, simply perform:
//     curl -i -F image=@data/lena.jpg localhost:8080
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tasdomas/faceapi/handler"
)

func main() {
	haar := flag.String("haar", "", "Haar classifier file")
	port := flag.String("http:", ":8080", "HTTP service address.")
	flag.Parse()

	handler, err := handler.NewDetectorHandler(*haar)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(*port, nil))
}
