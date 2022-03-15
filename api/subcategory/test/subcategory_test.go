package subcategory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetSubcategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server, _ := setup.NewTestApp(t)

	testCases := []struct {
		name          string
		url           string
		checkResponse func(*testing.T, *http.Response)
	}{
		{
			name: "OK",
			url:  "/api/v1/subcategory",
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodGet, tc.url, nil)
			require.Nil(t, err)

			resp, err := server.Test(request)
			require.Nil(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func TestAddSubcategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectionName := "subcategory"
	url := "/api/v1/connection/subcategory"

	app, gcpService := setup.NewTestApp(t)
	token := login(t, app)

	testCases := []struct {
		name          string
		method        string
		body          map[string]interface{}
		buildMock     func(gcpService *setup.MockUploader, collectionName string)
		checkResponse func(*testing.T, *http.Response)
	}{
		{
			name:   "OK",
			method: http.MethodPost,
			body: map[string]interface{}{
				"title": "test add subcate 121313",
				"image": setup.OpenFileFromPath("dummy_image.jpg"),
			},
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {
				gcpService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return("dummy image url", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "No require title",
			method: http.MethodPost,
			body: map[string]interface{}{
				"image": setup.OpenFileFromPath("dummy_image.jpg"),
			},
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:   "No require image",
			method: http.MethodPost,
			body: map[string]interface{}{
				"title": "test add subcate",
			},
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildMock(gcpService, collectionName)
			request, err := setup.Upload(tc.method, url, tc.body)
			request.Header.Add("Authorization", "bearer "+token)
			require.Nil(t, err)

			resp, err := app.Test(request)
			require.Nil(t, err)
			tc.checkResponse(t, resp)

			subcate := new(AddSubcategoryResult)
			bodyBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err)
			json.Unmarshal(bodyBytes, subcate)

			setup.DeleteMockSubcategory(t, &subcate.Result)
		})
	}
}

type AddSubcategoryResult struct {
	Result subcategory.Subcategory
}

type AddCategoryResult struct {
	Result category.Category
}

func TestAddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subcate1 := setup.AddMockSubcategory(t)
	// subcate2 := setup.AddMockSubcategory(t)

	// fmt.Println("subcate1", subcate1.ID.Hex())
	buf := new(bytes.Buffer)
	sids := []string{"623082745043756f707e5674", "623082745043756f707e5675"}
	// json.NewEncoder(buf).Encode([]string{subcate1.ID.Hex(), subcate2.ID.Hex()})
	json.NewEncoder(buf).Encode(sids)

	x := buf.Bytes()
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, x); err != nil {
		fmt.Println(err)
	}

	testCases := []struct {
		name           string
		body           map[string]interface{}
		collectionName string
		buildMock      func(gcpService *setup.MockUploader, collectionName string)
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"title": "test add cate",
				// "subcategory": strings.NewReader(subcate1.ID.Hex()),
				"subcategory": sids,
			},
			collectionName: "category",
			buildMock:      func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			app, gcpService := setup.NewTestApp(t)
			URL := "/api/v1/category"

			tc.buildMock(gcpService, tc.collectionName)

			// form := url.Values{}

			token := login(t, app)

			// request, err := setup.CreateMultipartFormDataRequest("POST", URL, tc.body)
			request, err := setup.Upload("POST", URL, tc.body)
			require.Nil(t, err)
			request.Header.Add("Authorization", "bearer "+token)

			// form.Add("title", "test add cate 1")
			// form.Add("subcategory", "622f61bed0570c74dfbef2a5")
			// form.Add("subcategory", "622f1c453f5016e3e450ccbd")

			// request, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
			// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			// request.Header.Add("Authorization", "bearer "+token)

			require.Nil(t, err)

			resp, err := app.Test(request)
			require.Nil(t, err)

			cate := new(AddCategoryResult)
			bodyBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err)
			json.Unmarshal(bodyBytes, cate)

			tc.checkResponse(t, resp)
			setup.DeleteMockCategory(t, &cate.Result)
			setup.DeleteMockSubcategory(t, subcate1)
		})
	}
}

type AuthResult struct {
	Result AuthResponse
}

type AuthResponse struct {
	Username     string
	AccessToken  string
	RefreshToken string
}

func login(t *testing.T, app *fiber.App) string {
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
