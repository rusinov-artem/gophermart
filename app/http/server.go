package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const defaultTimeout = 5 * time.Second
const shutdownTimout = 4 * time.Second
const exitCode = 2

type Server struct {
	s      *http.Server
	logger *zap.Logger
}

func NewServer(address string, mux http.Handler, logger *zap.Logger) *Server {
	s := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		IdleTimeout:       defaultTimeout,
	}
	return &Server{
		s:      s,
		logger: logger,
	}
}

func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.s.Addr)
	if err != nil {
		err := fmt.Errorf("unable to listen: %w", err)
		s.logger.Error(err.Error(), zap.Error(err))
		os.Exit(exitCode)
	}
	s.logger.Info(fmt.Sprintf("Listening %s", s.s.Addr))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		err := s.s.Serve(ln)
		fmt.Println("server exited:", err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Got signal:", <-c)

	ctx, closeFN := context.WithTimeout(context.Background(), shutdownTimout)
	defer closeFN()

	err = s.s.Shutdown(ctx)
	if err != nil {
		fmt.Println("error while shutdown:", err)
	}

	wg.Wait()
}
