package app

import (
	"net/http"

	"kathenovino/app/actions/absences"
	"kathenovino/app/middleware"
	"kathenovino/public"

	"github.com/gobuffalo/buffalo"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.RequestID)
	root.Use(middleware.Database)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)

	root.GET("/absences", absences.List)
	root.GET("/", absences.New)
	root.POST("/absences/create", absences.Create)
	root.DELETE("/absences/{id}/delete", absences.Delete)
	root.ServeFiles("/", http.FS(public.FS()))
}
