package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	// health check handler
	router.GET("/health", app.healthcheckHandler)

	// movies handler
	router.GET("/movies/:id", app.showMovieHandler)
	return router
}
