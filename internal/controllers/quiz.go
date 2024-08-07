package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/jimmykarily/quizmaker/internal/models"
)

const (
	COOKIE_NAME             = "quizmaker-cookie"
	COOKIE_LIFETIME_SEC     = 3600
	COOKIE_TIMESTAMP_FORMAT = "2006-01-02 15:04:05"
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
	err := ensureQuizSession(gctx)
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	qp, err := models.NewQuestionPoolFromFile(Settings.QuestionPoolFile)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// TODO: Don't hardcode
	q, err := qp.GenerateQuiz(models.QuizOptions{
		TotalQuestions:     15,
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

func ensureQuizSession(ctx *gin.Context) error {
	sc := securecookie.New([]byte(Settings.CookieSecret), nil)

	submittedEmail := ctx.Request.FormValue("email")
	if !models.ValidEmail(submittedEmail) {
		return errors.New("invalid email")
	}

	cookie, err := ctx.Request.Cookie(COOKIE_NAME)
	if err != nil { // no cookie found
		_, err := models.SessionForEmail(Settings.DB, submittedEmail)
		if err == nil {
			return errors.New("email has already been used previously")
		}

		_, err = newSessionForEmail(ctx, submittedEmail) // fresh email
		if err != nil {
			return fmt.Errorf("creating a new session: %w", err)
		}
		return nil
	}

	cookieValue := make(map[string]string)
	if err = sc.Decode(COOKIE_NAME, cookie.Value, &cookieValue); err != nil {
		return fmt.Errorf("invalid cookie format: %w", err)
	}

	if err := validTimestamp(cookieValue["timestamp"]); err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}

	// valid cookie with email. Let's lookup the session.
	_, err = models.SessionForEmail(Settings.DB, cookieValue["email"])
	// User has a valid cookie but we can't find a session.
	// Create a new one (we probably deleted the session from db).
	if err != nil {
		_, err = newSessionForEmail(ctx, cookieValue["email"])
		if err != nil {
			return err
		}
	}

	return nil
}

func validTimestamp(timestampStr string) error {
	timestamp, err := time.Parse(COOKIE_TIMESTAMP_FORMAT, timestampStr)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	if time.Since(timestamp).Seconds() > COOKIE_LIFETIME_SEC {
		return errors.New("cookie has expired")
	}

	return nil
}

func newSessionForEmail(ctx *gin.Context, email string) (models.Session, error) {
	var err error
	var result models.Session

	sc := securecookie.New([]byte(Settings.CookieSecret), nil)

	result, err = models.NewSessionForEmail(Settings.DB, email)
	if err != nil {
		return result, fmt.Errorf("creating a new session: %w", err)
	}

	// create the cookie too
	userAgent := ctx.Request.UserAgent()
	currentTimestamp := time.Now().Format(COOKIE_TIMESTAMP_FORMAT)
	value := map[string]string{
		"email":     email,
		"timestamp": currentTimestamp,
		"userAgent": userAgent,
	}

	encoded, err := sc.Encode(COOKIE_NAME, value)
	if err != nil {
		return result, fmt.Errorf("failed to encode cookie: %w", err)
	}

	cookie := &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(COOKIE_LIFETIME_SEC * time.Second),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, cookie)

	return result, nil
}
