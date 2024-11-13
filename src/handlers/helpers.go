package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"highload-sn-backend/db/postgres"
	"net/http"
	"time"

	"highload-sn-backend/internal/log"
	"highload-sn-backend/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func encodeError(w http.ResponseWriter, err error) {
	log.Logger().Error(err)

	code := statusCodeFromError(err)
	resp := &types.ErrorResponse{
		Message:   err.Error(),
		RequestId: "",
		Code:      code,
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}

func statusCodeFromError(err error) int {
	switch {
	case errors.Is(err, types.ErrInvalidData):
		return http.StatusBadRequest
	case errors.Is(err, types.ErrNoUser):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func validateLoginReq(req types.LoginRequest) error {
	if len(req.Id) == 0 {
		return fmt.Errorf("%w: id is required", types.ErrInvalidData)
	}
	if len(req.Password) == 0 {
		return fmt.Errorf("%w: password is required", types.ErrInvalidData)
	}

	users, err := postgres.GetUsers([]string{req.Id})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return fmt.Errorf("%w: user with id=%s not found", types.ErrNoUser, req.Id)
	}

	err = compareHashAndPassword(users[0].PasswordHash, req.Password)
	if err != nil {
		return fmt.Errorf("%w: incorrect password", types.ErrInvalidData)
	}

	return nil
}

func validateRegisterUserReq(req types.RegisterUserRequest) error {
	if len(req.FirstName) == 0 {
		return fmt.Errorf("%w: first name is required", types.ErrInvalidData)
	}
	if len(req.LastName) == 0 {
		return fmt.Errorf("%w: last name is required", types.ErrInvalidData)
	}
	if _, err := time.Parse(time.DateOnly, req.BirthDate); err != nil {
		return fmt.Errorf("%w: birth date is invalid, format is %v", types.ErrInvalidData, time.DateOnly)
	}
	if _, ok := types.SexSet[req.Sex]; !ok {
		return fmt.Errorf("%w: sex possible values: male, female", types.ErrInvalidData)
	}
	if len(req.Biography) == 0 {
		return fmt.Errorf("%w: biography is required", types.ErrInvalidData)
	}
	if len(req.City) == 0 {
		return fmt.Errorf("%w: city is required", types.ErrInvalidData)
	}
	if len(req.Password) == 0 {
		return fmt.Errorf("%w: password is required", types.ErrInvalidData)
	}

	return nil
}

func compareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func toDBUser(req types.RegisterUserRequest) (*types.User, error) {
	id := uuid.NewString()
	birthdate, _ := time.Parse(time.DateOnly, req.BirthDate)
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: can not get hash from password: %w", types.ErrInvalidData, err)
	}

	return &types.User{
		Id:           id,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    birthdate,
		Sex:          req.Sex,
		Biography:    req.Biography,
		City:         req.City,
		PasswordHash: string(hash),
	}, nil
}

func fromDBUser(user types.User) types.UserResponse {
	return types.UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate.Format(time.DateOnly),
		Sex:       user.Sex,
		Biography: user.Biography,
		City:      user.City,
	}
}
