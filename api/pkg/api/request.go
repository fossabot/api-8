package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

// Request represents an api request
type Request interface {
	ID() string
	ContentType() string
	Bind(out interface{}) error
	Raw() *http.Request
	Header() http.Header
	ClientIP() string
}

type DummyRequest struct {
	encodedBody []byte
	req         *http.Request
	id          string
}

// NewDummyRequest creates a new dummy request
func NewDummyRequest() *DummyRequest {
	req, _ := http.NewRequest("GET", "/", nil)
	return &DummyRequest{
		req: req,
		id:  "randomstring",
	}
}

func (r *DummyRequest) ID() string {
	return r.id
}

func (r *DummyRequest) ContentType() string {
	return r.req.Header.Get("content-type")
}

func (r *DummyRequest) Bind(out interface{}) error {
	return binding.Default(r.req.Method, r.ContentType()).Bind(r.req, out)
}

func (r *DummyRequest) Raw() *http.Request {
	return r.req
}

func (r *DummyRequest) SetMethod(method string) *DummyRequest {
	r.req.Method = method
	return r
}

func (r *DummyRequest) SetContentType(contentType string) *DummyRequest {
	r.req.Header.Set("content-type", contentType)
	return r
}

func (r *DummyRequest) AddHeader(key, val string) *DummyRequest {
	r.req.Header.Add(key, val)
	return r
}

func (r *DummyRequest) Header() http.Header {
	return r.req.Header
}

func (r *DummyRequest) ClientIP() string {
	return "127.0.0.1"
}

func (r *DummyRequest) SetJSONBody(p interface{}) *DummyRequest {
	c, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	r.updateBody(c)
	r.req.Method = "POST"
	r.req.Header.Set("content-type", ContentTypeJSON)
	return r
}

func (r *DummyRequest) updateBody(b []byte) {
	r.encodedBody = b
	r.req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
}
