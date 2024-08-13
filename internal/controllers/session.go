package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type (
	SessionController struct{}
)

func (c *SessionController) List(gctx *gin.Context) {
	sessions := []models.Session{}
	Settings.DB.Model(&models.Session{}).Select("sessions.*, count(questions.ID) as questions").
		Joins("left join questions on questions.session_email = sessions.email").
		Group("sessions.email").
		Scan(&sessions)

	viewData := struct {
		Completed  []models.Session
		InProgress []models.Session
	}{
		Completed:  sessions,
		InProgress: []models.Session{},
	}

	Render([]string{"main_layout", path.Join("sessions", "list")}, gctx.Writer, viewData)
}
