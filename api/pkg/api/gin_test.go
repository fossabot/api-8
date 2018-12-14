package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/logger"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/testhelper"
)

func init() {
	logger.SurpressLog()
}

type errResponse struct {
}

func (r *errResponse) StatusCode() int {
	return http.StatusInternalServerError
}

func (r *errResponse) Body() []byte {
	return nil
}

func (r *errResponse) Header() http.Header {
	return nil
}

func (r *errResponse) ContentType() string {
	return ContentTypeJSON
}

func TestWrapGin(t *testing.T) {
	tcs := []struct {
		Request          Request
		Handler          Handler
		ExpectedResponse Response
	}{
		{
			// normal request - response
			Request: NewDummyRequest().SetContentType(ContentTypeJSON),
			Handler: func(ctx context.Context, req Request) Response {
				resp := JSONResponse(http.StatusOK, map[string]interface{}{
					"hello": "world",
				})
				resp.Header().Add("first_key", "first_val1")
				resp.Header().Add("first_key", "first_val2")
				return resp
			},
			ExpectedResponse: &dummyResponse{
				statusCode: http.StatusOK,
				body:       []byte(`{"hello":"world"}`),
				header: map[string][]string{
					"first_key": []string{"first_val1", "first_val2"},
				},
				contentType: ContentTypeJSON,
			},
		},
	}

	for _, tc := range tcs {
		body, header, statusCode, err := runGin(WrapGin(testhelper.NewContext(), tc.Handler), tc.Request.Raw())
		if err != nil {
			t.Error("there was an error when running the handler", err)
		}
		exBody := tc.ExpectedResponse.Body()
		if strings.TrimSpace(string(body)) != strings.TrimSpace(string(exBody)) {
			t.Error("does not return expected response body", "got", string(body), "want", string(exBody))
		}
		if statusCode != tc.ExpectedResponse.StatusCode() {
			t.Error("does not return expected status code")
		}
		if reflect.DeepEqual(header, tc.ExpectedResponse.Header()) {
			t.Error("does not return expected header", "got", header, "want", tc.ExpectedResponse.Header())
		}
		if !strings.Contains(header.Get("content-type"), tc.ExpectedResponse.ContentType()) {
			t.Error("does not return expected content type, got", header.Get("content-type"))
		}
	}
}

func TestHandlePostJson(t *testing.T) {
	p := struct {
		Value string `json:"value"`
	}{}
	handler := func(ctx context.Context, req Request) Response {
		req.Bind(&p)
		return OKResp(nil)
	}
	req := NewDummyRequest().SetJSONBody(map[string]interface{}{"value": "something"})
	runGin(WrapGin(testhelper.NewContext(), handler), req.Raw())

	if p.Value != "something" {
		t.Error("does not work")
	}
}

func runGin(h gin.HandlerFunc, req *http.Request) ([]byte, http.Header, int, error) {
	port := getUnusedPort()
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	gin.SetMode("test")
	r := gin.New()
	if strings.ToUpper(req.Method) == "POST" {
		r.POST("/", h)
	} else {
		r.GET("/", h)
	}
	go func() {
		r.Run(addr)
	}()
	time.Sleep(100 * time.Millisecond) // wait until the server ready

	req.URL, _ = url.Parse("http://" + addr) // override the url
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, err
	}
	return body, resp.Header, resp.StatusCode, nil
}

func getUnusedPort() string {
	h, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	h.Close()
	return strings.Split(h.Addr().String(), ":")[1]
}
