package subcategory

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
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

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
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
		body          map[string]io.Reader
		buildMock     func(gcpService *setup.MockUploader, collectionName string)
		checkResponse func(*testing.T, *http.Response)
	}{
		{
			name:   "OK",
			method: http.MethodPost,
			body: map[string]io.Reader{
				"title": strings.NewReader("test add subcate"),
				"image": mustOpen("dummy_image.jpg"),
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
			body: map[string]io.Reader{
				"image": mustOpen("dummy_image.jpg"),
			},
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:   "No require image",
			method: http.MethodPost,
			body: map[string]io.Reader{
				"title": strings.NewReader("test add subcate"),
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
			request, err := CreateMultipartFormDataRequest(tc.method, url, tc.body)
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

func CreateMultipartFormDataRequest(method string, url string, values map[string]io.Reader) (req *http.Request, err error) {

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {

		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		// Add an image file
		if x, ok := r.(*os.File); ok {

			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {

				return nil, err
			}

		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}

		}

		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err = http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}
func TestAddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name           string
		body           fiber.Map
		collectionName string
		buildMock      func(gcpService *setup.MockUploader, collectionName string)
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"title":       "test add cate",
				"subcategory": []string{"622f61bed0570c74dfbef2a5"},
			},
			collectionName: "category",
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {
				// gcpService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return("dummy image url", nil)
			},
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

			form := url.Values{}

			token := login(t, app)
			form.Add("title", "test add cate 1")
			form.Add("subcategory", "622f61bed0570c74dfbef2a5")
			form.Add("subcategory", "622f1c453f5016e3e450ccbd")

			// form.Add("image", asString)
			request, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("Authorization", "bearer "+token)

			require.Nil(t, err)

			resp, err := app.Test(request)
			require.Nil(t, err)

			cate := new(AddCategoryResult)
			bodyBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err)
			json.Unmarshal(bodyBytes, cate)

			tc.checkResponse(t, resp)
			setup.DeleteMockCategory(t, &cate.Result)
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
