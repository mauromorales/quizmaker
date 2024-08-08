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
	"gorm.io/gorm/clause"
)

const (
	COOKIE_NAME             = "quizmaker-cookie"
	COOKIE_LIFETIME_SEC     = 3600
	COOKIE_TIMESTAMP_FORMAT = "2006-01-02 15:04:05"
)

type (
	QuizController struct{}
	CookieValue    struct {
		Email     string
		Timestamp string
		UserAgent string
	}
)

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

func (c *QuizController) Show(gctx *gin.Context) {
	currentSession, err := currentSession(gctx)
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	err = Settings.DB.Preload(clause.Associations).Find(&currentSession).Error
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	currentQuestion, err := currentSession.CurrentQuestion()
	// TODO: Return a flash error (when flashes are implemented)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}
	if currentQuestion.ID == 0 {
		// TODO: Quiz is finished, show an flash alert
		return
	}

	viewData := struct{ Question models.Question }{
		Question: currentQuestion,
	}

	Render([]string{"main_layout", path.Join("quizzes", "show")}, gctx.Writer, viewData)
}

func (c *QuizController) Create(gctx *gin.Context) {
	submittedEmail := gctx.Request.FormValue("email")
	if !models.ValidEmail(submittedEmail) {
		handleError(gctx.Writer, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	session, err := ensureQuizSession(gctx)
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	qp, err := models.NewQuestionPoolFromFile(Settings.QuestionPoolFile)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// TODO: Don't hardcode
	q, err := models.NewQuizWithOpts(models.QuizOptions{
		TotalQuestions:     15,
		MinDifficulty:      1,
		MaxDifficulty:      10,
		QuestionTimeoutSec: 10,
		AvailableQuestions: qp.Questions,
	})

	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	err = q.PersistForSessionEmail(Settings.DB, session.Email)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	redirectURL, err := GetFullURL(gctx.Request, "QuizShow")
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	gctx.Redirect(http.StatusFound, redirectURL)
}

func validCookieValue(ctx *gin.Context) (CookieValue, error) {
	var result CookieValue

	sc := securecookie.New([]byte(Settings.CookieSecret), nil)

	cookie, err := ctx.Request.Cookie(COOKIE_NAME)
	if err != nil { // no cookie found
		return result, err
	}

	if err = sc.Decode(COOKIE_NAME, cookie.Value, &result); err != nil {
		return result, fmt.Errorf("invalid cookie format: %w", err)
	}

	if err := validTimestamp(result.Timestamp); err != nil {
		return result, fmt.Errorf("invalid timestamp: %w", err)
	}

	return result, nil
}

func currentSession(ctx *gin.Context) (models.Session, error) {
	cookieValue, err := validCookieValue(ctx)
	if err != nil {
		return models.Session{}, err
	}

	session, err := models.SessionForEmail(Settings.DB, cookieValue.Email)
	if err != nil {
		return session, fmt.Errorf("finding user session: %w", err)
	}

	return session, nil
}

func ensureQuizSession(ctx *gin.Context) (models.Session, error) {
	var session models.Session

	submittedEmail := ctx.Request.FormValue("email")

	cookieValue, err := validCookieValue(ctx)
	if errors.Is(err, http.ErrNoCookie) { // no cookie found
		session, err = models.SessionForEmail(Settings.DB, submittedEmail)
		if err == nil {
			return session, errors.New("email has already been used previously")
		}

		return newSessionForEmail(ctx, submittedEmail) // fresh email
	}
	if err != nil { // other errors (expired or invalid cookie)
		return session, err
	}

	// valid cookie with email. Let's lookup the session.
	if cookieValue.Email != submittedEmail {
		return session, fmt.Errorf("already started with email: %s", cookieValue.Email)
	}

	// valid cookie with email. Let's lookup the session.
	session, err = models.SessionForEmail(Settings.DB, cookieValue.Email)
	// User has a valid cookie but we can't find a session.
	// Create a new one (we probably deleted the session from db).
	if err != nil {
		return newSessionForEmail(ctx, cookieValue.Email)
	}

	return session, nil
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
	value := CookieValue{
		Email:     email,
		Timestamp: currentTimestamp,
		UserAgent: userAgent,
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
