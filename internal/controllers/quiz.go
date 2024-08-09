package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	submitURL, err := GetFullURL(gctx.Request, "QuizCreate", nil)
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

	score := int(math.Round(models.QuestionList(currentSession.Questions).Score()))

	// Quiz is finished, show the results page
	if currentQuestion.ID == 0 {
		viewData := struct {
			Session         models.Session
			ScorePercentage string
		}{
			Session:         currentSession,
			ScorePercentage: strconv.Itoa(score),
		}
		Render([]string{"main_layout", path.Join("quizzes", "result")}, gctx.Writer, viewData)
		return
	}

	questionID := strconv.Itoa(int(currentQuestion.ID))
	submitURL, err := GetFullURL(gctx.Request, "QuestionAnswer", map[string]string{"id": questionID})
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// If it's the first time we show the question, make it "started"
	if currentQuestion.StartedAt.IsZero() {
		currentQuestion.StartedAt = time.Now()
		err = Settings.DB.Save(&currentQuestion).Error
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}
	}

	endTime := currentQuestion.StartedAt.Add(
		time.Duration(currentQuestion.AllowedSeconds) * time.Second)
	timeLeft := int(time.Until(endTime).Seconds())
	viewData := struct {
		Question  models.Question
		SubmitURL string
		TimeLeft  int
	}{
		Question:  currentQuestion,
		SubmitURL: submitURL,
		TimeLeft:  int(timeLeft),
	}

	Render([]string{"main_layout", path.Join("quizzes", "show")}, gctx.Writer, viewData)
}

func (c *QuizController) Create(gctx *gin.Context) {
	err := gctx.Request.ParseForm()
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	submittedEmail := gctx.Request.FormValue("email")
	if !models.ValidEmail(submittedEmail) {
		handleError(gctx.Writer, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	session, err := ensureQuizSession(gctx)
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	redirectURL, err := GetFullURL(gctx.Request, "QuizShow", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	// Reload the session with Questions
	err = Settings.DB.Preload(clause.Associations).Find(&session).Error
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}
	if len(session.Questions) > 0 { // Pre-existing session, just redirect to quiz
		gctx.Redirect(http.StatusFound, redirectURL)
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
		QuestionTimeoutSec: 30,
		AvailableQuestions: qp.Questions,
	})

	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	err = q.PersistForSessionEmail(Settings.DB, session.Email)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	gctx.Redirect(http.StatusFound, redirectURL)
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

	result, err = models.NewSessionForEmail(Settings.DB, email)
	if err != nil {
		return result, fmt.Errorf("creating a new session: %w", err)
	}

	// create the cookie too
	cookie, err := CreateCookie(email, ctx.Request.UserAgent())
	if err != nil {
		return result, fmt.Errorf("creating the cookie: %w", err)
	}
	http.SetCookie(ctx.Writer, cookie)

	return result, nil
}
