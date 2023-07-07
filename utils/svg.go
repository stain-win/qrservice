package utils

import (
	"encoding/xml"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
)

type QRBlock struct {
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Fill   string `xml:"fill,attr"`
}

var (
	blocksize   = 16
	borderwidth = 4
)

type QRSVG struct {
	XMLName  xml.Name  `xml:"svg"`
	NS       string    `xml:"xmlns,attr"`
	Width    int       `xml:"width,attr"`
	Height   int       `xml:"height,attr"`
	Style    string    `xml:"style,attr"`
	QRBlocks []QRBlock `xml:"rect"`
}

func NewQRSVG(qrc *qrcode.QRCode) string {
	var i image.Image = qrc.Image(20)
	var w int = i.Bounds().Max.X

	var svg QRSVG

	svg.NS = "http://www.w3.org/2000/svg"
	svg.Width = (w + 2*borderwidth) * blocksize
	svg.Height = svg.Width

	svg.QRBlocks = make([]QRBlock, 1+w*w)
	svg.QRBlocks[0] = QRBlock{0, 0, svg.Width, svg.Height, hex(qrc.BackgroundColor)}

	for x := 0; x < w; x++ {
		for y := 0; y < w; y++ {
			svg.QRBlocks[1+x*w+y] = QRBlock{
				X:      (x + borderwidth) * blocksize,
				Y:      (y + borderwidth) * blocksize,
				Width:  blocksize,
				Height: blocksize,
				Fill:   hex(i.At(x, y)),
			}
		}
	}
	x, _ := xml.Marshal(svg)
	return string(x)
}

func hex(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
