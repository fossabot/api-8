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
		start := time.Now()
		var req = newGinRequest(gCtx)

		defer func(requestID string) {
			if thing := recover(); thing != nil {
				logrus.
					WithError(fmt.Errorf("%v", thing)).
					WithField("method", gCtx.Request.Method).
					WithField("path", gCtx.Request.URL).
					WithField("request_id", requestID).
					Errorln("panic while handling request")
			}
		}(req.ID())

		resp := h(parent, req)
		body := resp.Body()

		for k, v := range resp.Header() {
			for _, h := range v {
				gCtx.Writer.Header().Add(k, h)
			}
		}
		gCtx.Writer.Header().Add("content-type", resp.ContentType())
		gCtx.Writer.Header().Add("X-Request-ID", req.ID())

		gCtx.Writer.WriteHeader(resp.StatusCode())
		gCtx.Writer.Write(body)

		logrus.
			WithField("request_id", req.ID()).
			WithField("duration", time.Since(start)/time.Millisecond).
			WithField("method", gCtx.Request.Method).
			WithField("path", gCtx.Request.URL).
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

func (r *ginRequest) ClientIP() string {
	return r.gCtx.ClientIP()
}
