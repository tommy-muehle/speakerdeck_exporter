package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	isHealthy int32
)

// Server provides methods to run an HTTP server with custom handlers.
type Server struct {
	addr   string
	router *http.ServeMux
	quit   chan os.Signal
}

// NewServer creates a Server.
func NewServer(addr string) *Server {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &Server{
		addr:   addr,
		router: http.NewServeMux(),
		quit:   quit,
	}
}

// AddHandler adds an HTTP Handler for a given route.
func (s *Server) AddHandler(route string, handler http.Handler) {
	s.router.Handle(route, handler)
}

// ListenAndServe is running the HTTP Server.
func (s *Server) ListenAndServe() {
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	done := make(chan bool)

	go func() {
		<-s.quit
		atomic.StoreInt32(&isHealthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		close(done)
	}()

	atomic.StoreInt32(&isHealthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-done
}

// Shutdown the HTTP Server.
func (s *Server) Shutdown() {
	s.quit <- syscall.SIGINT
}
