package qrhandlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/skip2/go-qrcode"
	"github.com/stain-win/qrservice/errorhandler"
	"image/color"
	"net/http"
)
var qrContentKey = "content"
type qrHandler struct {}

type QrJsonResponse struct {
	Content string `json:"content"`
	Matrix [][]int `json:"matrix"`
}

func (qrj QrJsonResponse) Render (w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (qr qrHandler) Routes () chi.Router {
	r := chi.NewRouter()
	r.Use(qr.QRContext)
	r.Get("/", qr.ByteQr)
	r.Get("/img", qr.ImgQr)
	r.Get("/json", qr.JsonQr)
	return r
}

func (qr qrHandler) QRContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctn := r.URL.Query().Get("ctn")
		if ctn == "" {
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
	fmt.Println(len(c.Bitmap()))
	render.Status(r, http.StatusCreated)
	//w.Write([]byte(c.ToSmallString(true)))
	w.Write([]byte(boolToByte(c)))

}

func (qr qrHandler) JsonQr(w http.ResponseWriter, r *http.Request) {
	content := r.Context().Value(qrContentKey).(string)
	c, _ := qrcode.New(content, qrcode.Medium)

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewQrJsonResponse(c))
}

func (qr qrHandler) ImgQr(w http.ResponseWriter, r *http.Request){
	content := r.Context().Value(qrContentKey).(string)
	var png []byte
	c, _ := qrcode.New(content, qrcode.Medium)
	c.BackgroundColor = color.RGBA{0x00, 0x00, 0xff, 0xff}
	c.ForegroundColor = color.RGBA{0x44, 0xFF, 0x00, 0xff}
	png, _ = c.PNG(512)
	//png, _ = qrcode.Encode(content, qrcode.Medium, 256)
	//err = qrcode.WriteColorFile(content, qrcode.Medium, 256, color.Gray, "0000FF", "qr.png")
	w.Write(png)
}

func NewQrHandler () *qrHandler {
	return &qrHandler{}
}

func NewQrJsonResponse (qrCode *qrcode.QRCode) *QrJsonResponse {
	resp := &QrJsonResponse{}
	resp.Content = qrCode.Content
	resp.Matrix = boolToInt(qrCode)

	return resp
}

func boolToByte (q *qrcode.QRCode) string {
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

func boolToInt (q *qrcode.QRCode) [][]int {
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
