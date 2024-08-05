package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
)

type QuizController struct {
}

func (c *QuizController) Create(gctx *gin.Context) {
	// TODO: Generate a new quiz and redirect to the view page
	Render([]string{"main_layout", path.Join("quizzes", "show")}, gctx.Writer, nil)
}
