package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stain-win/qrservice/healthcheckhandler"
	"github.com/stain-win/qrservice/qrhandlers"
	"net/http"
	"time"
)

const (
	portDelimiter = ":"
)

var (
	port      = flag.String("port", "3200", "http service port")
	startTime time.Time
)

func main() {
	startTime = time.Now()
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

	http.ListenAndServe(portDelimiter+*port, r)
}
