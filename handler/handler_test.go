package handler_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "launchpad.net/gocheck"

	"github.com/tasdomas/faceapi/handler"
)

func Test(t *testing.T) {
	TestingT(t)
}

type TSuite struct {
	handler      http.Handler
	multipartReq bytes.Buffer
	contentType  string
}

var _ = Suite(&TSuite{})

func (t *TSuite) SetUpSuite(c *C) {
	var err error
	t.handler, err = handler.NewDetectorHandler("../data/haarcascade_frontalface_alt.xml")
	c.Assert(err, IsNil)

	w := multipart.NewWriter(&t.multipartReq)
	// Add your image file
	f, err := os.Open("../data/lena.jpg")
	c.Assert(err, IsNil)

	fw, err := w.CreateFormFile("image", "lena.jpg")
	c.Assert(err, IsNil)

	_, err = io.Copy(fw, f)
	c.Assert(err, IsNil)
	w.Close()
	t.contentType = w.FormDataContentType()
}

// A request that is not post will be responded to
// with a Bad Request status code.
func (t *TSuite) TestInvalidRequest(c *C) {
	req, err := http.NewRequest("GET", "localhost", nil)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	t.handler.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, http.StatusBadRequest)
}

// A POST request that submits no image will be responded to
// with a Bad Request status code.
func (t *TSuite) TestInvalidPostRequest(c *C) {
	req, err := http.NewRequest("POST", "localhost", nil)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	t.handler.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, http.StatusBadRequest)
}

// A valid POST request gets a response with json encoded
// object rectangle data.
func (t *TSuite) TestValidPostRequest(c *C) {
	req, err := http.NewRequest("POST", "localhost", &t.multipartReq)
	c.Assert(err, IsNil)
	req.Header.Set("Content-Type", t.contentType)
	w := httptest.NewRecorder()
	t.handler.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "{\"Objects\":[{\"A\":{\"X\":214,\"Y\":201},\"B\":{\"X\":390,\"Y\":377}}]}\n")
}

func (t *TSuite) BenchmarkHandlerDetection(c *C) {
	w := httptest.NewRecorder()

	for i := 0; i < c.N; i++ {
		form := t.multipartReq
		req, err := http.NewRequest("POST", "localhost", &form)
		c.Assert(err, IsNil)
		req.Header.Set("Content-Type", t.contentType)
		t.handler.ServeHTTP(w, req)
		c.Assert(w.Code, Equals, http.StatusOK)
	}
}
