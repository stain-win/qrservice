package healthcheckhandler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

type HealthCheckHandler struct {
	TimeStart time.Time
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type HealthCheckStatusResponse struct {
	Status  string  `json:"status"`
	Elapsed float64 `json:"elapsed"`
}

func (hcr *HealthCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (hcr *HealthCheckStatusResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (hc HealthCheckHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", hc.HealthCheck)
	r.Get("/status", hc.Status)
	return r
}

func (hc HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewHealthCheckResponse())
}

func (hc HealthCheckHandler) Status(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewHealthCheckStatusResponse(hc.TimeStart))
}

func NewHealthCheckHandler(startTime time.Time) *HealthCheckHandler {

	return &HealthCheckHandler{TimeStart: startTime}
}

func NewHealthCheckResponse() *HealthCheckResponse {
	resp := &HealthCheckResponse{}
	resp.Status = "OK"
	return resp
}

func NewHealthCheckStatusResponse(timeStart time.Time) *HealthCheckStatusResponse {
	resp := &HealthCheckStatusResponse{}
	resp.Status = "OK"
	resp.Elapsed = time.Since(timeStart).Seconds()
	return resp
}
