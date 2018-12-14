package api

import (
	"encoding/json"
	"net/http"

	raven "github.com/getsentry/raven-go"
)

// Response represents an api response
type Response interface {
	StatusCode() int
	Body() []byte
	Header() http.Header
	ContentType() string
}

type jsonResponse struct {
	data       interface{}
	header     http.Header
	statusCode int
}

// JSONResponse creates a json response
func JSONResponse(code int, data interface{}) Response {
	return &jsonResponse{
		statusCode: code,
		data:       data,
		header:     http.Header{},
	}
}

func (r *jsonResponse) StatusCode() int {
	return r.statusCode
}

func (r *jsonResponse) Body() []byte {
	b, err := json.Marshal(r.data)
	if err != nil {
		raven.CaptureError(err, nil)
		return nil
	}
	return b
}

func (r *jsonResponse) Header() http.Header {
	return r.header
}

func (r *jsonResponse) ContentType() string {
	return ContentTypeJSON
}

type dummyResponse struct {
	statusCode  int
	body        []byte
	header      http.Header
	contentType string
}

func (r *dummyResponse) StatusCode() int {
	return r.statusCode
}

func (r *dummyResponse) Body() []byte {
	return r.body
}

func (r *dummyResponse) Header() http.Header {
	return r.header
}

func (r *dummyResponse) ContentType() string {
	return r.contentType
}
