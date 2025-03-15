package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/bcrypt"
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func NewUlidFromTimestamp(time time.Time) (string, error) {
	ms := ulid.Timestamp(time)
	entropy := ulid.Monotonic(rand.Reader, 0)

	id, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func GetUserDifferenceData(DbUser entity.User, NewUser auth.UpdateUserRequest) (entity.User, error) {
	var result entity.User
	result.ID = DbUser.ID

	if NewUser.Username != "" && NewUser.Username != DbUser.Username {
		result.Username = NewUser.Username
	}

	if NewUser.Password != "" {
		hashedPass, err := bcrypt.HashPassword(NewUser.Password)
		if err != nil {
			return entity.User{}, err
		}
		result.Password = hashedPass
	}

	return result, nil
}

func MakeUserData(user entity.User) map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	}
}
