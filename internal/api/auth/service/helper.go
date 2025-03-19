package authService

import (
	"ProjectGolang/internal/entity"
	"crypto/rand"
	"fmt"
	"github.com/oklog/ulid/v2"
	"math/big"
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

//
//func GetUserDifferenceData(DbUser entity.User, NewUser auth.UpdateUserRequest) (entity.User, error) {
//	var result entity.User
//	result.ID = DbUser.ID
//
//	if NewUser.Username != "" && NewUser.Username != DbUser.Username {
//		result.Username = NewUser.Username
//	}
//
//	if NewUser.Password != "" {
//		hashedPass, err := bcrypt.HashPassword(NewUser.Password)
//		if err != nil {
//			return entity.User{}, err
//		}
//		result.Password = hashedPass
//	}
//
//	return result, nil
//}
//

func makeUserData(user entity.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"role":       user.Role,
		"is_premium": user.IsPremium,
	}
}

func generateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("OTP length must be greater than 0")
	}

	var otp string

	for i := 0; i < length; i++ {
		// Generate a random number between 0-9
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}

		otp += fmt.Sprintf("%d", num)
	}

	return otp, nil
}
