package faceapi

import (
	"fmt"
	"os"

	"github.com/hybridgroup/go-opencv/opencv"
)

// Detector allows detecting objects in images.
type Detector struct {
	cascade *opencv.HaarCascade
}

// Create new detector from supplied Haar classifier file.
func NewDetector(classifierFile string) (*Detector, error) {
	// check if file exists
	if st, err := os.Stat(classifierFile); os.IsNotExist(err) || st.IsDir() {
		return nil, fmt.Errorf("Classifier file does not exist.")
	}

	cascade := opencv.LoadHaarClassifierCascade(classifierFile)
	return &Detector{cascade: cascade}, nil
}

// Single point.
type Point struct {
	X, Y int
}

// A rectangle.
type Rect struct {
	A, B Point
}

type DetectorResult struct {
	Objects []Rect
}

func resultFromRects(rects []*opencv.Rect) DetectorResult {
	result := DetectorResult{}
	for _, rect := range rects {
		result.Objects = append(result.Objects,
			Rect{
				A: Point{rect.X(),
					rect.Y()},
				B: Point{rect.X() + rect.Width(),
					rect.Y() + rect.Height()},
			})
	}
	return result
}

func (d Detector) FindInFile(imageFile string) (DetectorResult, error) {
	image := opencv.LoadImage(imageFile)
	if image == nil {
		return DetectorResult{}, fmt.Errorf("Failed to load image from file.")
	}
	defer image.Release()

	rects := d.cascade.DetectObjects(image)
	return resultFromRects(rects), nil
}
