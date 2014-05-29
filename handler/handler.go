package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tasdomas/faceapi"
)

// DetectorHandler detects objects in images posted to it.
// It responds with coordinates of detected object rectangles.
type DetectorHandler struct {
	detector *faceapi.Detector
}

// Create a DetectorHandler with specified Haar classifier file.
func NewDetectorHandler(classifierDefinition string) (*DetectorHandler, error) {
	detector, err := faceapi.NewDetector(classifierDefinition)
	if err != nil {
		return nil, err
	}
	return &DetectorHandler{detector}, nil
}

// Handle a posted image with created Haar classifier.
func (h DetectorHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// check the request method
	if req.Method != "POST" {
		log.Print("Invalid request.")
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(resp, "Expecting post request.")
		return
	}

	// retrieve uploaded image
	image, _, err := req.FormFile("image")
	if err != nil {
		log.Printf("Invalid request: %v.", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	// save temp file
	tmp, err := ioutil.TempFile(os.TempDir(), "facedetector")
	if err != nil {
		log.Printf("Failed to create temp file: %v.", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tmp.Close()
	defer image.Close()
	_, err = io.Copy(tmp, image)
	if err != nil {
		log.Printf("Failed to save image: %v.", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	results, err := h.detector.FindInFile(tmp.Name())
	if err != nil {
		log.Printf("Failed to analyze image: %v.", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(resp)
	err = out.Encode(results)
	if err != nil {
		log.Printf("Failed to send results: %v.", err)
		return
	}

}
