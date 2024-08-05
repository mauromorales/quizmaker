package controllers

import (
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
	}

	return routes
}
