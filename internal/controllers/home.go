package controllers

import (
	"encoding/base64"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type HomeController struct {
}

var QuizNewQRImageMemoization map[string][]byte

func (c *HomeController) Index(gctx *gin.Context) {
	NewQuizURL, err := GetFullURL(gctx.Request, "QuizNew")
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	png, err := getQRCodePNG(NewQuizURL)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	result := struct {
		QRCodePNG  string
		NewQuizURL string
	}{
		QRCodePNG:  base64.StdEncoding.EncodeToString(png),
		NewQuizURL: NewQuizURL,
	}

	Render([]string{"main_layout", path.Join("home", "index")}, gctx.Writer, result)
}

func getQRCodePNG(url string) ([]byte, error) {
	var png []byte
	var cached bool
	var err error

	png, cached = QuizNewQRImageMemoization[url]
	if cached && len(QuizNewQRImageMemoization[url]) > 0 {
		Settings.InfoLogger.Println("Using memoized QR code")
		return png, nil
	}

	if png, err = qrcode.Encode(url, qrcode.Medium, 512); err != nil {
		return png, err
	}

	if QuizNewQRImageMemoization == nil {
		Settings.InfoLogger.Println("non initialized map")
		QuizNewQRImageMemoization = map[string][]byte{
			url: png,
		}
	} else {
		QuizNewQRImageMemoization[url] = png
	}

	return png, nil
}
