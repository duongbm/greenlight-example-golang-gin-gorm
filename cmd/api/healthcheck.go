package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) healthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "Available",
		"environment": app.config.env,
		"version":     version,
	})
}
