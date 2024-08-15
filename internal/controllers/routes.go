package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Route describes a route for httprouter
type Route struct {
	Name    string
	Method  string
	Path    string
	Format  string
	Handler gin.HandlerFunc
}

type Routes []Route

func GetRoutes() Routes {
	routes := Routes{
		Route{
			Name:    "SessionList",
			Method:  "GET",
			Path:    "/",
			Format:  "html",
			Handler: (&SessionController{}).List,
		},
		Route{
			Name:    "QuizNew",
			Method:  "GET",
			Path:    "/quizzes/new",
			Format:  "html",
			Handler: (&QuizController{}).New,
		},
		Route{
			Name:    "QuizCreate",
			Method:  "POST",
			Path:    "/quizzes",
			Format:  "html",
			Handler: (&QuizController{}).Create,
		},
		Route{
			Name:    "QuizShow",
			Method:  "GET",
			Path:    "/quiz",
			Format:  "html",
			Handler: (&QuizController{}).Show,
		},
		Route{
			Name:    "QuestionAnswer",
			Method:  "POST",
			Path:    "/questions/:id",
			Format:  "html",
			Handler: (&QuestionController{}).Answer,
		},
	}

	return routes
}

func RouteByName(name string) (Route, error) {
	for _, r := range GetRoutes() {
		if r.Name == name {
			return r, nil
		}
	}
	return Route{}, fmt.Errorf("route %s not found", name)
}

func GetFullURL(request *http.Request, routeName string, params map[string]string) (string, error) {
	var err error
	u := url.URL{Host: request.Host}

	if request.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if u.Path, err = GetRoutePath(routeName, params); err != nil {
		return "", err
	}

	return u.String(), nil
}

func GetRoutePath(routeName string, params map[string]string) (string, error) {
	var path string
	var err error
	var route Route

	if route, err = RouteByName(routeName); err != nil {
		return "", err
	}
	path = route.Path

	// Replace path parameters
	for key, value := range params {
		// Escape path parameters
		encodedValue := url.PathEscape(value)
		path = strings.ReplaceAll(path, ":"+key, encodedValue)
	}

	return path, nil
}
