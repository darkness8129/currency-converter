package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	server   *http.Server
	router   *gin.Engine
	notifyCh chan error
}

type Options struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func NewGinHTTPServer(opt Options) *HTTPServer {
	router := gin.New()

	httpServer := &http.Server{
		Handler:      router,
		Addr:         opt.Addr,
		WriteTimeout: opt.WriteTimeout,
		ReadTimeout:  opt.ReadTimeout,
	}

	return &HTTPServer{
		server:   httpServer,
		router:   router,
		notifyCh: make(chan error, 1),
	}
}

func (s *HTTPServer) Start() {
	go func() {
		defer close(s.notifyCh)
		s.notifyCh <- s.server.ListenAndServe()
	}()
}

func (s *HTTPServer) Notify() <-chan error {
	return s.notifyCh
}

func (s *HTTPServer) Router() *gin.Engine {
	return s.router
}

func (s *HTTPServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
