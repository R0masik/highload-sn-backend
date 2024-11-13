package types

import (
	"fmt"
	"time"
)

const (
	Male   = "male"
	Female = "female"
)

var (
	SexSet = map[string]struct{}{
		Male:   {},
		Female: {},
	}

	ErrInvalidData = fmt.Errorf("invalid data")
	ErrNoUser      = fmt.Errorf("no user")
)

// requests

type LoginRequest struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Biography string `json:"biography"`
	City      string `json:"city"`
	Password  string `json:"password"`
}

// responses

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserIdResponse struct {
	UserId string `json:"user_id"`
}

type UserResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Biography string `json:"biography"`
	City      string `json:"city"`
}

// db

type QueryItem struct {
	Query  string
	Params []any
}

type User struct {
	Id           string    `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	BirthDate    time.Time `db:"birthdate"`
	Sex          string    `db:"sex"`
	Biography    string    `db:"biography"`
	City         string    `db:"city"`
	PasswordHash string    `db:"password_hash"`
}

type Session struct {
	UserId string `db:"user_id"`
	Token  string `db:"token"`
}
