package data

import (
	"github.com/duongbm/greenlight-gin/internal/validator"
	pq "github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type MovieModel struct {
	DB *gorm.DB
}

func (m *MovieModel) Insert(movie *Movie) error {
	return m.DB.Table("movies").Create(movie).Error
}

func (m *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var movie Movie
	query := m.DB.Table("movies").Find(&movie, id)
	if query.RowsAffected == 0 {
		return nil, ErrRecordNotFound
	}
	if query.Error != nil {
		return nil, query.Error
	}
	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = ?, year = ?, runtime = ?, genres = ? , version = version + 1
		WHERE id = ? AND version = ?
		RETURNING version`

	tx := m.DB.Raw(query, movie.Title, movie.Year, movie.Runtime, movie.Genres, movie.Id, movie.Version).Scan(&movie)
	if tx.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return tx.Error
}

func (m *MovieModel) Delete(id int64) error {
	query := m.DB.Table("movies").Where("id = ?", id).Delete(&Movie{})
	if query.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return query.Error
}

type Movie struct {
	Id        int64          `json:"id"`
	CreatedAt time.Time      `json:"-"`
	Title     string         `json:"title"`
	Year      int32          `json:"year,omitempty"`
	Runtime   Runtime        `json:"runtime,omitempty"`
	Genres    pq.StringArray `json:"genres,omitempty" gorm:"type:text[]"`
	Version   int32          `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime >= 0, "runtime", "must be positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate genres")
}
