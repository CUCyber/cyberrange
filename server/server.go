package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

type controller struct {
	logger        *log.Logger
	nextRequestID func() string
	healthy       int64
}

var c *controller

func currentTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func disableDirList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *controller) shutdown(ctx context.Context, server *http.Server) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		atomic.StoreInt64(&c.healthy, 0)
		server.ErrorLog.Printf("Server is shutting down...\n")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %s\n", err)
		}
	}()

	return ctx
}

func serve(listenAddr string, logPath string) {
	writer, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("Could not create log file at path: %s\n", logPath)
		return
	}

	logger := log.New(writer, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	c = &controller{logger: logger, nextRequestID: currentTime}

	router := http.NewServeMux()
	router.HandleFunc("/", c.index)
	router.HandleFunc("/health", c.health)

	fs := http.FileServer(http.Dir("web/static"))
	fs = http.StripPrefix("/static/", disableDirList(fs))
	router.Handle("/static/", fs)

	server := &http.Server{
		Addr: listenAddr,
		Handler: (middlewares{
			c.tracing,
			c.logging,
			c.restore}).apply(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	ctx := c.shutdown(context.Background(), server)

	logger.Printf("Server is ready to handle requests at %q\n", listenAddr)
	atomic.StoreInt64(&c.healthy, time.Now().UnixNano())

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %q: %s\n", listenAddr, err)
	}
	<-ctx.Done()
	logger.Printf("Server stopped\n")
}
