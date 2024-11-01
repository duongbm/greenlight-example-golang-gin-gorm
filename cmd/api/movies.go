package main

import (
	"github.com/duongbm/greenlight-gin/internal/data"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (app *application) createMovieHandler(c *gin.Context) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, "create new movie")
}

func (app *application) showMovieHandler(c *gin.Context) {
	id := c.Param("id")

	_id, _ := strconv.ParseInt(id, 10, 64)
	movie := data.Movie{
		Id:        _id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	c.JSON(http.StatusOK, movie)
}
