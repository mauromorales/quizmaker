package controllers

import (
	"encoding/base64"
	"net/http"
	"path"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type (
	SessionController struct{}
)

func (c *SessionController) List(gctx *gin.Context) {
	sessions := []models.Session{}
	err := Settings.DB.Find(&sessions).Error
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	var complete, inProgress []models.Session
	for _, s := range sessions {
		if s.Complete {
			complete = append(complete, s)
		} else {
			inProgress = append(inProgress, s)
		}
	}

	sort.Slice(complete, func(i, j int) bool {
		return complete[i].Score > complete[j].Score
	})

	NewQuizURL, err := GetFullURL(gctx.Request, "QuizNew", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	png, err := getQRCodePNG(NewQuizURL)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	qp, err := models.NewQuestionPoolFromFile(Settings.QuestionPoolFile)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	viewData := struct {
		QRCodePNG  string
		NewQuizURL string
		Completed  []models.Session
		InProgress []models.Session
		Prices     models.PricesList
	}{
		QRCodePNG:  base64.StdEncoding.EncodeToString(png),
		NewQuizURL: NewQuizURL,
		Completed:  complete,
		InProgress: inProgress,
		Prices:     qp.Prices,
	}

	Render([]string{"main_layout", path.Join("sessions", "list")}, gctx.Writer, viewData)
}
