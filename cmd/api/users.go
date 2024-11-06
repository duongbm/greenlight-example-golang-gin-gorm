package main

import (
	"errors"
	"github.com/duongbm/greenlight-gin/internal/data"
	"github.com/duongbm/greenlight-gin/internal/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) registerUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: true,
	}

	// set password user
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(c, v.Errors)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.Error(err, nil)
		}
	})

	c.JSON(http.StatusCreated, user)
}
