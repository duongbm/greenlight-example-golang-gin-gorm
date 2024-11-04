package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				app.serverErrorResponse(c, fmt.Errorf("%s", err))
				return
			}
		}()
		c.Next()
	}
}
