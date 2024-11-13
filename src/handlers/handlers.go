package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"

	"highload-sn-backend/db/postgres"
	"highload-sn-backend/types"

	"github.com/gorilla/mux"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		encodeError(w, err)
	}
	defer r.Body.Close()

	var req types.LoginRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		encodeError(w, err)
		return
	}

	err = validateLoginReq(req)
	if err != nil {
		encodeError(w, err)
		return
	}

	token, err := postgres.GetSession(req.Id)
	if err != nil {
		if !errors.Is(err, types.ErrNoUser) {
			encodeError(w, err)
			return
		}

		token = uuid.NewString()
		err = postgres.AddSession(req.Id, token)
		if err != nil {
			encodeError(w, err)
			return
		}
	}

	_ = json.NewEncoder(w).Encode(types.TokenResponse{Token: token})
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		encodeError(w, err)
	}
	defer r.Body.Close()

	var req types.RegisterUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		encodeError(w, err)
		return
	}

	err = validateRegisterUserReq(req)
	if err != nil {
		encodeError(w, err)
		return
	}

	user, err := toDBUser(req)
	if err != nil {
		encodeError(w, err)
		return
	}

	err = postgres.AddUsers([]types.User{*user})
	if err != nil {
		encodeError(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(types.UserIdResponse{UserId: user.Id})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		encodeError(w, fmt.Errorf("%w: id is required", types.ErrInvalidData))
		return
	}

	users, err := postgres.GetUsers([]string{id})
	if err != nil {
		encodeError(w, err)
		return
	}

	if len(users) == 0 {
		encodeError(w, fmt.Errorf("%w: user with id=%s not found", types.ErrNoUser, id))
		return
	}

	_ = json.NewEncoder(w).Encode(fromDBUser(users[0]))
}
