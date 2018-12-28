package api

import (
	"net/http"
	"reflect"
	"testing"
)

func TestMiddlewareReturnUnauthorized(t *testing.T) {
	var h Handler = func(req Request) Response {
		return JSONResponse(200, map[string]interface{}{"hello": "world"})
	}
	var m middleware = func(next Handler) Handler {
		return func(req Request) Response {
			return UnauthorizedResp()
		}
	}
	resp := m(h)(NewDummyRequest())
	if resp.StatusCode() != http.StatusUnauthorized {
		t.Error("should return unauthorized")
	}
}

func TestMiddlewareRunOrder(t *testing.T) {
	order := []string{}
	run := false
	var m middleware = func(next Handler) Handler {
		return func(req Request) Response {
			order = append(order, "m-before")
			resp := next(req)
			order = append(order, "m-after")
			run = true
			return resp
		}
	}
	var h Handler = func(req Request) Response {
		order = append(order, "h")
		return &jsonResponse{}
	}

	m(h)(NewDummyRequest())
	if !run {
		t.Error("middleware does not run")
	}
	if order[0] != "m-before" || order[1] != "h" || order[2] != "m-after" {
		t.Error("invalid run order", order)
	}
}

func TestWithMiddleware(t *testing.T) {
	order := []string{}
	var m1 middleware = func(next Handler) Handler {
		return func(req Request) Response {
			order = append(order, "m1-before")
			resp := next(req)
			order = append(order, "m1-after")
			return resp
		}
	}
	var m2 middleware = func(next Handler) Handler {
		return func(req Request) Response {
			order = append(order, "m2-before")
			resp := next(req)
			order = append(order, "m2-after")
			return resp
		}
	}
	var m3 middleware = func(next Handler) Handler {
		return func(req Request) Response {
			order = append(order, "m3-before")
			resp := next(req)
			order = append(order, "m3-after")
			return resp
		}
	}
	var h Handler = func(req Request) Response {
		order = append(order, "h")
		return &jsonResponse{}
	}

	h = withMiddleware(h, m1, m2, m3)
	h(NewDummyRequest())
	expected := []string{
		"m3-before",
		"m2-before",
		"m1-before",
		"h",
		"m1-after",
		"m2-after",
		"m3-after",
	}
	if !reflect.DeepEqual(order, expected) {
		t.Error("invalid run order", order)
	}
}
