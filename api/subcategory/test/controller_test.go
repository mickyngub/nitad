package subcategory_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//
func randomImages(n int) []*multipart.FileHeader {
	results := make([]*multipart.FileHeader, n)
	for i := 0; i < n; i++ {
		results[i] = &multipart.FileHeader{
			Filename: utils.RandomString(5) + ".jpeg",
		}
	}
	return results
}

// func TestGetSubcategoryById(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		checkResponse func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "OK",
// 			checkResponse: func(t *testing.T, resp *http.Response) {
// 				require.Equal(t, http.StatusOK, resp.StatusCode)
// 			},
// 		},
// 		{
// 			name: "Unauthorized",
// 			checkResponse: func(t *testing.T, resp *http.Response) {
// 				require.Equal(t, http.StatusOK, resp.StatusCode)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			fmt.Println(3)

// 			server := newTestServer(t)

// 			url := "/api/v1/subcategory"
// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			fmt.Println(1)
// 			require.Equal(t, err, nil)

// 			fmt.Println(2)

// 			resp, err := server.Test(request)

// 			fmt.Println(resp)

// 			require.Equal(t, err, nil)

// 			tc.checkResponse(t, resp)

// 		})
// 	}
// }

func TestAddSubcategoryById(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// gcpService := NewMockClientUploader(ctrl)

	// gcpService.EXPECT

	testCases := []struct {
		name           string
		body           map[string]interface{}
		collectionName string
		buildMock      func(gcpService *MockClientUploader, collectionName string)
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"title": "test title dummy subcategory",
				"image": randomImages(2),
			},
			collectionName: "asdfasdfasdf",
			buildMock: func(gcpService *MockClientUploader, collectionName string) {
				gcpService.EXPECT().UploadImages(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return([]string{"dummy url"}, nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "FAIL",
			buildMock: func(gcpService *MockClientUploader, collectionName string) {
				gcpService.EXPECT().UploadImages(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gcpService := NewMockClientUploader(ctrl)

			tc.buildMock(gcpService, tc.collectionName)

			server := newTestServer(t, gcpService)

			url := "/api/v1/subcategory"

			data, _ := json.Marshal(tc.body)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.Equal(t, err, nil)

			resp, err := server.Test(request)

			fmt.Println(resp)

			require.Equal(t, err, nil)

			tc.checkResponse(t, resp)

		})
	}
}

// it should add a subcategory successfully
// it should not add a subcategory if required fields are not provided correctly

// it should list all added subcategories

// it should get the subcategory by valid id
