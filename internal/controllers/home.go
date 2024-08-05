package controllers

import (
	"encoding/base64"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	qrcode "github.com/skip2/go-qrcode"
)

type HomeController struct {
}

func (c *HomeController) Index(gctx *gin.Context) {
	url, err := GetFullURL("QuizNew")
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	var png []byte
	// TODO: Memoize this url. It won't change at least until the app is restarted
	png, err = qrcode.Encode(url, qrcode.Medium, 512)
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
