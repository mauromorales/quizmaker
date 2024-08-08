package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type (
	QuestionController struct{}
)

func (c *QuestionController) Answer(gctx *gin.Context) {
	err := gctx.Request.ParseForm()
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	qid := gctx.Param("id")
	selectedAnswer := gctx.Request.FormValue("answer")
	// TODO: When no answer is selected, flash an error and redirect to quiz again

	var question models.Question
	err = Settings.DB.First(&question, "ID = ?", qid).Error
	if handleError(gctx.Writer, err, http.StatusNotFound) {
		return
	}

	question.UserAnswer, err = strconv.Atoi(selectedAnswer)
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	err = Settings.DB.Save(&question).Error
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	redirectURL, err := GetFullURL(gctx.Request, "QuizShow", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	gctx.Redirect(http.StatusFound, redirectURL)
}
