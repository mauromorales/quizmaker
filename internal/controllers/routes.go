package controllers

import (
	"fmt"
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
			Handler: (&QuizController{}).Create,
		},
	}

	return routes
}

func GetFullURL(routeName string) (string, error) {
	u, err := url.Parse(Settings.Host)
	if err != nil {
		return "", fmt.Errorf("parsing host: %w", err)
	}

	u.Path = ""
	for _, r := range GetRoutes() {
		if r.Name == routeName {
			u.Path = r.Path
		}
	}

	if u.Path == "" {
		return "", fmt.Errorf("no route %s found", routeName)
	}

	return u.String(), nil
}
