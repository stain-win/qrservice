package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"qrservice/healthcheckhandler"
	"qrservice/qrhandlers"
	"time"
)

const (
	portDelimiter = ":"
)

var (
	startTime time.Time
)

func main() {
	startTime = time.Now()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3200"
	}

	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	qr := qrhandlers.NewQrHandler()
	hc := healthcheckhandler.NewHealthCheckHandler(startTime)
	r.Mount("/qr", qr.Routes())
	r.Mount("/healthcheck", hc.Routes())

	http.ListenAndServe(portDelimiter+port, r)
}
