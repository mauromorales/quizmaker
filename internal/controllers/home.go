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

var QuizNewQRImageMemoization []byte

func (c *HomeController) Index(gctx *gin.Context) {
	png, err := getQRCodePNG()
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	result := struct {
		QRCodePNG string
	}{
		QRCodePNG: base64.StdEncoding.EncodeToString(png),
	}

	Render([]string{"main_layout", path.Join("home", "index")}, gctx.Writer, result)
}

func getQRCodePNG() ([]byte, error) {
	if len(QuizNewQRImageMemoization) > 0 {
		Settings.InfoLogger.Println("Using memoized QR code")
		return QuizNewQRImageMemoization, nil
	}

	var png []byte

	url, err := GetFullURL("QuizNew")
	if err != nil {
		return png, err
	}

	if png, err = qrcode.Encode(url, qrcode.Medium, 512); err != nil {
		return png, err
	}

	QuizNewQRImageMemoization = png

	return png, nil
}
