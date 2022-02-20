package admin

import (
	"os"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/utils"
	"github.com/golang-jwt/jwt/v4"
)

type MsgToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

var JWT_SECRET = os.Getenv("JWT_SECRET")

func CreateToken(admin *Admin) (MsgToken, errors.CustomError) {
	var msgToken MsgToken
	claims := jwt.MapClaims{
		"id":       admin.ID,
		"username": admin.Username,
		"sub":      utils.UUID(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return msgToken, errors.NewBadRequestError(err.Error())
	}

	msgToken.AccessToken = accessToken

	rtClaims := jwt.MapClaims{
		"sub": utils.UUID(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return msgToken, errors.NewBadRequestError(err.Error())
	}
	msgToken.RefreshToken = refreshToken

	return msgToken, nil
}
