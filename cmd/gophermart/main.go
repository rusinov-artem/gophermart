package main

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
)

func main() {
	fmt.Println("Hello World")

	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	addr := ":7777"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("unable to listen: %w", err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		err := srv.Serve(ln)
		fmt.Println("server exited:", err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Got signal:", <-c)

	ctx, closeFN := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeFN()

	err = srv.Shutdown(ctx)
	if err != nil {
		fmt.Println("error while shutdown:", err)
	}

	wg.Wait()
}
