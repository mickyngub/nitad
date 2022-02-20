package admin

import (
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, errors.CustomError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalServerError("Failed Hashing password")
	}
	return string(hashedPassword), nil

}

func ComparePassword(hashed, password string) errors.CustomError {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return errors.NewBadRequestError("Invalid Credentials")
	}
	return nil
}

func CreateJWTToken(a Admin) (map[string]interface{}, errors.CustomError) {
	var result map[string]interface{}
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin_id"] = a.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return result, errors.NewInternalServerError("Failed creating JWT")
	}

	result = map[string]interface{}{
		"id":       a.ID,
		"username": a.Username,
		"token":    t,
		"exp":      exp,
	}
	return result, nil
}
