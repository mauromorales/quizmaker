package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type (
	QuestionController struct{}
)

func (c *QuestionController) Answer(gctx *gin.Context) {
	session, err := currentSession(gctx)
	if handleError(gctx.Writer, err, http.StatusUnauthorized) {
		return
	}

	err = gctx.Request.ParseForm()
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	qid := gctx.Param("id")
	selectedAnswer := gctx.Request.FormValue("answer")
	// TODO: if answer is empty for whatever reason, flash an error and redirect to quiz show
	// TODO: If the question doesn't belong to the current session, return an error
	// TODO: When no answer is selected, flash an error and redirect to quiz again
	fmt.Printf("session.ID = %+v\n", session.ID)

	var question models.Question
	err = Settings.DB.First(&question, "ID = ?", qid).Error
	if handleError(gctx.Writer, err, http.StatusNotFound) {
		return
	}

	redirectURL, err := GetFullURL(gctx.Request, "QuizShow", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// Don't allow answering expired or already answered questions
	if question.Expired() || question.UserAnswer != 0 {
		// TODO: Flash error
	} else {
		question.UserAnswer, err = strconv.Atoi(selectedAnswer)
		if handleError(gctx.Writer, err, http.StatusBadRequest) {
			return
		}

		err = Settings.DB.Save(&question).Error
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}
	}

	gctx.Redirect(http.StatusFound, redirectURL)
}
