package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/oklog/ulid/v2"
	"math/big"
	"time"
)

func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %v", err)
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func NewUlidFromTimestamp(time time.Time) (string, error) {
	ms := ulid.Timestamp(time)
	entropy := ulid.Monotonic(rand.Reader, 0)

	id, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
