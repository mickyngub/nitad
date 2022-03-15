package setup

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type AuthResult struct {
	Result AuthResponse
}

type AuthResponse struct {
	Username     string
	AccessToken  string
	RefreshToken string
}

func Login(t *testing.T, app *fiber.App) string {
	form := url.Values{}

	URL := "/api/v1/admin/login"
	form.Add("username", os.Getenv("ADMIN_USERNAME"))
	form.Add("password", os.Getenv("ADMIN_PASSWORD"))

	request, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	require.Nil(t, err)

	resp, err := app.Test(request)
	require.Nil(t, err)
	bodyBytes, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	result := AuthResult{}
	err = json.Unmarshal(bodyBytes, &result)
	require.Nil(t, err)

	return result.Result.AccessToken
}
