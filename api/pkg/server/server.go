package server

import (
	"net/http"
	"time"
)

// Run runs http server.
// addr is ip and port where server should listen to.
func Run(addr string, isProduction bool) error {
	router := buildRouter(isProduction)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return srv.ListenAndServe()
}
