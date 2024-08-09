package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	templatepkg "text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/jimmykarily/quizmaker/internal/models"
	settingspkg "github.com/jimmykarily/quizmaker/internal/settings"
)

var Settings settingspkg.Settings

// Render renders the given templates using the provided data and writes the result
// to the provided ResponseWriter.
func Render(templates []string, w http.ResponseWriter, data interface{}) {
	var (
		err         error
		tmplFile    *os.File
		tmplContent []byte
	)

	tmpl := templatepkg.New("page_template")
	tmpl = tmpl.Delims("[[", "]]")
	for _, template := range templates {
		tmplFile, err = os.Open("views/" + template + ".html")
		if err != nil {
			break
		}
		tmplContent, err = io.ReadAll(tmplFile)
		if err != nil {
			break
		}

		tmpl = templatepkg.Must(tmpl.Funcs(templatepkg.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"sub": func(a, b int) int {
				return a - b
			},
		}).Parse(string(tmplContent)))
	}

	if handleError(w, err, 500) {
		return
	}

	err = tmpl.ExecuteTemplate(w, templates[0], data)
	if handleError(w, err, 500) {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SetupRoutes(e *gin.Engine, routes Routes) {
	e.Static("/assets", "./assets")
	for _, r := range routes {
		e.Handle(r.Method, r.Path, r.Handler)
	}
}

// Write the error to the response writer and return  true if there was an error
func handleError(w http.ResponseWriter, err error, code int) bool {
	if err != nil {
		if Settings.ErrorLogger != nil { // we don't set it in tests
			Settings.ErrorLogger.Println(err.Error())
		}
		http.Error(w, err.Error(), code)
		return true
	}
	return false
}

func validCookieValue(ctx *gin.Context) (CookieValue, error) {
	var result CookieValue

	sc := securecookie.New([]byte(Settings.CookieSecret), nil)

	cookie, err := ctx.Request.Cookie(COOKIE_NAME)
	if err != nil { // no cookie found
		return result, fmt.Errorf("finding the %s cookie: %w", COOKIE_NAME, err)
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

func CreateCookie(email, userAgent string) (*http.Cookie, error) {
	currentTimestamp := time.Now().Format(COOKIE_TIMESTAMP_FORMAT)
	value := CookieValue{
		Email:     email,
		Timestamp: currentTimestamp,
		UserAgent: userAgent,
	}

	sc := securecookie.New([]byte(Settings.CookieSecret), nil)
	encoded, err := sc.Encode(COOKIE_NAME, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode cookie: %w", err)
	}

	return &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(COOKIE_LIFETIME_SEC * time.Second),
		HttpOnly: true,
	}, nil
}
