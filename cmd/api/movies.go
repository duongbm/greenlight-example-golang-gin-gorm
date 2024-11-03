package main

import (
	"errors"
	"github.com/duongbm/greenlight-gin/internal/data"
	"github.com/duongbm/greenlight-gin/internal/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (app *application) createMovieHandler(c *gin.Context) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (app *application) showMovieHandler(c *gin.Context) {
	id := c.Param("id")

	_id, _ := strconv.ParseInt(id, 10, 64)

	movie, err := app.models.Movies.Get(_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, movie)
}
