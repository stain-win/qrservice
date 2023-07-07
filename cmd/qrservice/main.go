package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	qr := qrhandlers.NewQrHandler()
	hc := healthcheckhandler.NewHealthCheckHandler(startTime)
	r.Mount("/qr", qr.Routes())
	r.Mount("/healthcheck", hc.Routes())

	http.ListenAndServe(portDelimiter+port, r)
}
