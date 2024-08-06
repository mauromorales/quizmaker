package controllers

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type QuizController struct {
}

func (c *QuizController) New(gctx *gin.Context) {
	submitURL, err := GetFullURL(gctx.Request, "QuizCreate")
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	viewData := struct {
		SubmitURL string
	}{
		SubmitURL: submitURL,
	}

	Render([]string{"main_layout", path.Join("quizzes", "new")}, gctx.Writer, viewData)
}

func (c *QuizController) Create(gctx *gin.Context) {
	qp, err := models.NewQuestionPoolFromFile(Settings.QuestionPoolFile)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// TODO: Don't hardcode
	q, err := qp.GenerateQuiz(models.QuizOptions{
		TotalQuestions:     10,
		MinDifficulty:      1,
		MaxDifficulty:      10,
		QuestionTimeoutSec: 10,
	})

	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	viewData := struct {
		Quiz models.Quiz
	}{
		Quiz: q,
	}

	Render([]string{"main_layout", path.Join("quizzes", "show")}, gctx.Writer, viewData)
}
