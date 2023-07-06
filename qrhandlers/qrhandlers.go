package qrhandlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/skip2/go-qrcode"
	"image/color"
	"net/http"
	"qrservice/errorhandler"
	"qrservice/utils"
)

var qrContentKey = "content"

type qrHandler struct{}

type QrJsonResponse struct {
	Content string  `json:"content"`
	Matrix  [][]int `json:"matrix"`
}

type QrCode struct {
	Content string `json:"qr_content"`
	BgColor string `json:"bg_color,omitempty"`
	Color   string `json:"color,omitempty"`
	Level   int    `json:"err_level,omitempty"`
	Size    int    `json:"size,omitempty"`
	Output  string `json:"output,omitempty"`
	Border  bool   `json:"border,omitempty"`
}

type QrCodeRequest struct {
	*QrCode
}

func (qrc *QrCodeRequest) Bind(r *http.Request) error {
	if qrc.QrCode == nil {
		return errors.New("missing required params")
	}

	return nil
}

func (qrj QrJsonResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (qr qrHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(qr.QRContext)
	r.Get("/", qr.ByteQr)
	r.Post("/img", qr.ImgQr)
	r.Get("/json", qr.JsonQr)
	return r
}

func (qr qrHandler) QRContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctn := r.URL.Query().Get("ctn")
		if ctn == "" && r.Method == "GET" {
			render.Render(w, r, errorhandler.ErrorRenderer(fmt.Errorf("QR code content is required")))
			return
		}
		queryParams := ctn
		ctx := context.WithValue(r.Context(), qrContentKey, queryParams)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (qr qrHandler) ByteQr(w http.ResponseWriter, r *http.Request) {
	content := r.Context().Value(qrContentKey).(string)
	c, _ := qrcode.New(content, qrcode.Medium)
	c.BackgroundColor = color.RGBA{0x00, 0x00, 0x33, 0xff}
	render.Status(r, http.StatusCreated)
	w.Write([]byte(boolToByte(c)))

}

func (qr qrHandler) JsonQr(w http.ResponseWriter, r *http.Request) {
	content := r.Context().Value(qrContentKey).(string)
	c, _ := qrcode.New(content, qrcode.Medium)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, NewQrJsonResponse(c))
}

func (qr qrHandler) ImgQr(w http.ResponseWriter, r *http.Request) {
	var c *qrcode.QRCode
	var png []byte
	content := &QrCodeRequest{}

	if err := render.Bind(r, content); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	c, _ = qrcode.New(content.Content, qrcode.Medium)
	fmt.Println(content.BgColor)
	c.BackgroundColor, _ = utils.ParseHexColor(content.BgColor)
	c.ForegroundColor, _ = utils.ParseHexColor(content.Color)

	png, _ = c.PNG(content.Size)

	render.Data(w, r, png)
}

func NewQrHandler() *qrHandler {
	return &qrHandler{}
}

func NewQrJsonResponse(qrCode *qrcode.QRCode) *QrJsonResponse {
	resp := &QrJsonResponse{}
	resp.Content = qrCode.Content
	resp.Matrix = boolToInt(qrCode)

	return resp
}

func boolToByte(q *qrcode.QRCode) string {
	bits := q.Bitmap()
	var buf bytes.Buffer
	for y := range bits {
		for x := range bits[y] {
			if bits[y][x] {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func boolToInt(q *qrcode.QRCode) [][]int {
	bits := q.Bitmap()
	intArray := make([][]int, len(bits))

	for i := range intArray {
		intArray[i] = make([]int, 0, len(bits[i]))
	}

	for y := range bits {
		for x := range bits[y] {
			if bits[y][x] {
				intArray[x] = append(intArray[x], 1)
			} else {
				intArray[x] = append(intArray[x], 0)
			}
		}
	}
	return intArray
}

// TODO move this to separate module
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
