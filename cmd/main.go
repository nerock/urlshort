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

	"github.com/nerock/urlshort/docs"
	"github.com/nerock/urlshort/server"
	"github.com/nerock/urlshort/url"
	urlgenerator "github.com/nerock/urlshort/url/generator"
	urlrouter "github.com/nerock/urlshort/url/router"
	urlstore "github.com/nerock/urlshort/url/store"
)

const (
	defaultPort = 8080
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	db, err := sql.Open("sqlite3", getDBConnection())
	if err != nil {
		log.Fatal("could not establish connection with sqlite db:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("could not properly close connection to sqlite db")
		}
	}()

	httpSrv := server.NewHTTPServer(getHttpPort())

	urlStore, err := urlstore.NewURLStore(db)
	if err != nil {
		log.Fatal(err)
	}
	urlService := url.NewService(getDomain(), urlgenerator.URLGenerator{}, urlStore)
	urlRouter := urlrouter.NewURLRouter(urlService)

	docsRouter := docs.Router{}

	go func() {
		if err := httpSrv.Run(urlRouter, docsRouter); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return defaultPort
}

func getDomain() string {
	if domain := os.Getenv("DOMAIN"); domain != "" {
		return domain
	}

	return fmt.Sprintf("localhost:%d/", getHttpPort())
}

func getDBConnection() string {
	if dbConn := os.Getenv("DBCONN"); dbConn != "" {
		return dbConn
	}

	return "urlshort.db"
}
