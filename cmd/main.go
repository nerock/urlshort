package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nerock/urlshort/doc"
	"github.com/nerock/urlshort/server"
)

const (
	defaultPort = 8080
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	httpSrv := server.NewHTTPServer(getHttpPort())

	docsRouter := doc.Router{}

	go func() {
		if err := httpSrv.Run(docsRouter); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("error running HTTP server:", err)
		}
	}()

	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		cancel()
	}()

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Fatal("timeout gracefully shutting down")
	}
}

func getHttpPort() int {
	if portStr := os.Getenv("port"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return defaultPort
}
