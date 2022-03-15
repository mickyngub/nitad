package subcategory

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type AddSubcategoryResult struct {
	Result subcategory.Subcategory
}

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
				"title": "test add subcate",
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
			request.Header.Add("Authorization", "bearer "+setup.Token)
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
