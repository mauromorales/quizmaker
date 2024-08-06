package controllers

import (
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
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

	tmpl := template.New("page_template")
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

		tmpl, err = tmpl.Parse(string(tmplContent))
		if err != nil {
			break
		}
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
