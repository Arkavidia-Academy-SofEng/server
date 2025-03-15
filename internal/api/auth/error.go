package auth

import (
	"fmt"
)

var (
	ErrorEmailAlreadyExists = fmt.Errorf("email already exists")
	ErrorInvalidCredentials = fmt.Errorf("invalid email or password")
	ErrorUserNotFound       = fmt.Errorf("user not found")
)
