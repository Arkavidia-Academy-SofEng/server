package auth

import (
	"ProjectGolang/pkg/response"
	"net/http"
)

var (
	ErrEmailAlreadyExists = response.New(http.StatusConflict, "email already exists")
	ErrUserNotFound       = response.New(http.StatusNotFound, "user not found")
)
