package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nerock/urlshort/grpc"

	"github.com/nerock/urlshort/docs"
	"github.com/nerock/urlshort/server"
	"github.com/nerock/urlshort/url"
	urlgenerator "github.com/nerock/urlshort/url/generator"
	urlrouter "github.com/nerock/urlshort/url/router"
	urlstore "github.com/nerock/urlshort/url/store"
)

const (
	defaultPort     = 8080
	defaultGRPCPort = 50051
	defaultDomain   = "localhost"
	defaultDBConn   = "urlshort.db"
)

func main() {
	// Quit app signal notifier
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// DB Connection
	db, err := sql.Open("sqlite3", getDBConnection())
	if err != nil {
		log.Fatal("could not establish connection with sqlite db:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("could not properly close connection to sqlite db")
		}
	}()

	// Services startup
	urlStore, err := urlstore.NewURLStore(db)
	if err != nil {
		log.Fatal(err)
	}
	urlService := url.NewService(getDomain(), urlgenerator.URLGenerator{}, urlStore)
	urlGrpc := urlrouter.NewURLgRPC(urlService)
	urlRouter := urlrouter.NewURLRouter(urlService)

	docsRouter := docs.Router{}

	// Servers startup
	httpSrv := server.NewHTTPServer(getHttpPort())
	grpcSrv := grpc.NewGRPCServer(urlGrpc)
	go func() {
		if err := httpSrv.Run(urlRouter, docsRouter); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("error running HTTP server:", err)
		}
	}()

	go func() {
		if err := grpcSrv.RunServer(getGRPCPort()); err != nil {
			log.Fatal("error running gRPC server:", err)
		}
	}()

	// Wait for quit signal
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		grpcSrv.Shutdown()

		cancel()
	}()

	// Force quit if deadline exceeded
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Fatal("timeout gracefully shutting down")
	}
}

func getHttpPort() int {
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return defaultPort
}

func getGRPCPort() int {
	if portStr := os.Getenv("GRPC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return defaultGRPCPort
}

func getDomain() string {
	if domain := os.Getenv("DOMAIN"); domain != "" {
		return domain
	}

	return fmt.Sprintf("%s:%d/", defaultDomain, getHttpPort())
}

func getDBConnection() string {
	if dbConn := os.Getenv("DB_CONN"); dbConn != "" {
		return dbConn
	}

	return defaultDBConn
}
