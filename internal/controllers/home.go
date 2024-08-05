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
	var png []byte
	png, err := qrcode.Encode(Settings.Host, qrcode.Medium, 256)
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
