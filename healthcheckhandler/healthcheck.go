package healthcheckhandler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type HealthCheckHandler struct {
}

type HealthCheckResponse struct {
	Status string `json:"status"`
	Elapsed int64 `json:"elapsed"`
}

func (hcr *HealthCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	hcr.Elapsed = 4
	return nil
}

func (hc HealthCheckHandler) Routes () chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", hc.HealthCheck)
	return r
}

func (hc HealthCheckHandler) HealthCheck (w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewHealthCheckResponse())
}

func NewHealthCheckHandler () *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func NewHealthCheckResponse () *HealthCheckResponse {
	resp := &HealthCheckResponse{}
	resp.Status = "OK"
	return resp
}
