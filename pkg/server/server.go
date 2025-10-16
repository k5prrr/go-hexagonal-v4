package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer      *http.Server
	ShutdownTimeout time.Duration
}

func New(addr string, router http.Handler, readHeaderTimeout time.Duration) *Server {
	if addr == "" {
		addr = ":8080"
	}
	if router == nil {
		router = http.DefaultServeMux
	}
	if readHeaderTimeout == 0 {
		readHeaderTimeout = time.Second * 10
	}

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: readHeaderTimeout,
			Handler:           router,
		},
		ShutdownTimeout: 5 * time.Second,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		log.Printf("Server is running on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown failed: %v\n", err)
		return err
	}

	log.Printf("http server stopped")

	return nil
}
