package controllers

import (
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

	viewData := struct {
		Completed  []models.Session
		InProgress []models.Session
	}{
		Completed:  complete,
		InProgress: inProgress,
	}

	Render([]string{"main_layout", path.Join("sessions", "list")}, gctx.Writer, viewData)
}
