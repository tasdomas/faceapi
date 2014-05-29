package faceapi_test

import (
	"testing"

	"github.com/tasdomas/faceapi"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type TSuite struct{}

const (
	HAAR_CASCADE_FILE = "data/haarcascade_frontalface_alt.xml"
	IMAGE_FILE        = "data/lena.jpg"
)

var _ = Suite(&TSuite{})

// Test initializing detector with non-existant file fails.
func (t *TSuite) TestNonExistantFile(c *C) {
	detector, err := faceapi.NewDetector("nonexistantfile.xml")
	c.Assert(detector, IsNil)
	c.Assert(err, ErrorMatches, "Classifier file does not exist.")
}

// Test detector initialization.
func (t *TSuite) TestDetectorInit(c *C) {
	detector, err := faceapi.NewDetector(HAAR_CASCADE_FILE)
	c.Assert(detector, NotNil)
	c.Assert(err, IsNil)
}

// Test detector fails when image file does not exist.
func (t *TSuite) TestNoImage(c *C) {
	detector, err := faceapi.NewDetector(HAAR_CASCADE_FILE)
	c.Assert(detector, NotNil)
	c.Assert(err, IsNil)

	result, err := detector.FindInFile("nonexistantimage.jpg")
	c.Assert(result.Objects, HasLen, 0)
	c.Assert(err, ErrorMatches, "Failed to load image from file.")
}

// Test detection on image file.
func (t *TSuite) TestImageFile(c *C) {
	detector, err := faceapi.NewDetector(HAAR_CASCADE_FILE)
	c.Assert(detector, NotNil)
	c.Assert(err, IsNil)

	result, err := detector.FindInFile(IMAGE_FILE)
	c.Assert(result.Objects, HasLen, 1)
	c.Assert(err, IsNil)
}
