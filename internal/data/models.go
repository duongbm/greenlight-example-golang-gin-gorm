package data

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	User   UserModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Movies: MovieModel{db},
		User:   UserModel{db},
	}
}
