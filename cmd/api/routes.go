package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	// attach middleware
	router.Use(app.recoverPanic())

	router.NoRoute(app.notFoundResponse)
	router.NoMethod(app.methodNotAllowedResponse)

	// health check handler
	router.GET("/health", app.healthcheckHandler)

	// movies handler
	router.GET("/movies", app.listMovieHandler)
	router.GET("/movies/:id", app.showMovieHandler)
	router.PUT("/movies/:id", app.updateMovieHandler)
	router.PATCH("/movies/:id", app.partialUpdateMovieHandler)
	router.DELETE("/movies/:id", app.deleteMovieHandler)
	router.POST("/movies", app.createMovieHandler)

	router.POST("/users", app.registerUserHandler)
	return router
}
