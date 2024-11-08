package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) logError(c *gin.Context, err error) {
	app.logger.Error(err, map[string]string{"request_method": c.Request.Method, "request_url": c.Request.URL.String()})
}

func (app *application) errorResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(status, gin.H{"error": message})
}

func (app *application) serverErrorResponse(c *gin.Context, err error) {
	app.logError(c, err)
	message := "the server encountered a problem and could not process the request"
	app.errorResponse(c, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(c *gin.Context) {
	message := "the request method is not allowed"
	app.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(c *gin.Context, err map[string]string) {
	app.errorResponse(c, http.StatusUnprocessableEntity, err)
}

func (app *application) editConflictResponse(c *gin.Context) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(c, http.StatusUnprocessableEntity, message)
}
