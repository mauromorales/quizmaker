package controllers

import (
	"fmt"
	"net/http"
	"net/url"

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
			Name:    "Home",
			Method:  "GET",
			Path:    "/",
			Format:  "html",
			Handler: (&HomeController{}).Index,
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

func GetFullURL(request *http.Request, routeName string) (string, error) {
	u := url.URL{Host: request.Host}

	if request.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if r, err := RouteByName(routeName); err == nil {
		u.Path = r.Path
	}

	if u.Path == "" {
		return "", fmt.Errorf("no route %s found", routeName)
	}

	return u.String(), nil
}
