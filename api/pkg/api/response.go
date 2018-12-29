package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
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
		logrus.WithError(err).Errorln("failed to parse response body")
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

type codeResponse struct {
	code int
}

func CodeOnlyResp(code int) Response {
	return &codeResponse{
		code: code,
	}
}

func (r *codeResponse) StatusCode() int {
	return r.code
}

func (r *codeResponse) Body() []byte {
	return nil
}

func (r *codeResponse) Header() http.Header {
	return nil
}

func (r *codeResponse) ContentType() string {
	return ""
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
