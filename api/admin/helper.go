package admin

import (
	"github.com/birdglove2/nitad-backend/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, errors.CustomError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalServerError("Failed Hashing password")
	}
	return string(hashedPassword), nil

}
