package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

// WrapGin wraps a Handler and turns it into gin compatible handler
// This method should be called with a fresh ctx
func WrapGin(parent context.Context, h Handler) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		defer func() {
			if thing := recover(); thing != nil {
				logrus.WithError(fmt.Errorf("%v", thing)).Errorln("panic while handling request")
			}
		}()

		start := time.Now()

		// create request and run the handler
		var req = newGinRequest(gCtx)
		resp := h(parent, req)

		// get the body first
		body := resp.Body()

		// write header
		for k, v := range resp.Header() {
			for _, h := range v {
				gCtx.Writer.Header().Add(k, h)
			}
		}
		gCtx.Writer.Header().Add("content-type", resp.ContentType())
		gCtx.Writer.Header().Add("X-Request-ID", req.ID())

		// write body and status
		gCtx.Writer.Write(body)
		gCtx.Writer.WriteHeader(resp.StatusCode())

		// access log
		logrus.
			WithField("request_id", req.ID()).
			WithField("duration", time.Since(start)/time.Millisecond).
			WithField("method", gCtx.Request.Method).
			WithField("url", gCtx.Request.URL).
			WithField("headers", gCtx.Request.Header).
			WithField("status", resp.StatusCode()).
			WithField("resp_body_length", len(body)).
			Infoln("finished handling request")
	}
}

type ginRequest struct {
	gCtx *gin.Context
	id   string
}

func newGinRequest(gCtx *gin.Context) Request {
	return &ginRequest{
		gCtx: gCtx,
		id:   xid.New().String(),
	}
}

func (r *ginRequest) ID() string {
	return r.id
}

func (r *ginRequest) Bind(out interface{}) error {
	return r.gCtx.Bind(out)
}

func (r *ginRequest) Header() http.Header {
	return r.gCtx.Request.Header
}

func (r *ginRequest) ContentType() string {
	return r.gCtx.ContentType()
}

func (r *ginRequest) Raw() *http.Request {
	return r.gCtx.Request
}
