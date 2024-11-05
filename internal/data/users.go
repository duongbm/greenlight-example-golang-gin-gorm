package data

import (
	"errors"
	"github.com/duongbm/greenlight-gin/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(user *User) error {
	query := m.DB.Create(&user)
	if query.Error != nil {
		switch {
		case query.Error.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return query.Error
		}
	}
	return nil
}

func (m *UserModel) GetByEmail(email string) (user *User, err error) {
	query := m.DB.Where("email = ?", email).Find(&user)
	if query.Error != nil {
		switch {
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, query.Error
		}
	}
	return user, nil
}

func (m *UserModel) Update(user *User) error {
	query := m.DB.Clauses(clause.Returning{}).Where("id = ?", user.Id).Updates(&user)
	if query.Error != nil {
		switch {
		case query.Error.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(query.Error, gorm.ErrRecordNotFound):
			return ErrEditConflict
		default:
			return query.Error
		}
	}
	return nil
}

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must be not more 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must be at least 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
