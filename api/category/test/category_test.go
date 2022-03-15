package category

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type AddCategoryResult struct {
	Result category.Category
}

func TestAddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app, gcpService := setup.NewTestApp(t)

	subcate1 := setup.AddMockSubcategory(t)
	subcate2 := setup.AddMockSubcategory(t)

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
				"title":       "test add cate",
				"subcategory": []string{subcate1.ID.Hex(), subcate2.ID.Hex()},
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
			URL := "/api/v1/category"

			tc.buildMock(gcpService, tc.collectionName)

			request, err := setup.Upload("POST", URL, tc.body)
			require.Nil(t, err)
			request.Header.Add("Authorization", "bearer "+setup.Token)

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

func TestEditCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app, gcpService := setup.NewTestApp(t)

	subcate1 := setup.AddMockSubcategory(t)
	subcate2 := setup.AddMockSubcategory(t)

	cate := setup.AddMockCategory(t, subcate1)

	url := "/api/v1/category/" + cate.ID.Hex()

	testCases := []struct {
		name           string
		body           map[string]interface{}
		method         string
		collectionName string
		buildMock      func(gcpService *setup.MockUploader, collectionName string)
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "OK add subcate",
			body: map[string]interface{}{
				"title":       "test edit cate",
				"subcategory": []string{subcate1.ID.Hex(), subcate2.ID.Hex()},
			},
			collectionName: "category",
			method:         http.MethodPut,
			buildMock:      func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "OK change subcate",
			body: map[string]interface{}{
				"title":       "test edit cate",
				"subcategory": []string{subcate2.ID.Hex()},
			},
			collectionName: "category",
			method:         http.MethodPut,
			buildMock:      func(gcpService *setup.MockUploader, collectionName string) {},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.buildMock(gcpService, tc.collectionName)

			request, err := setup.Upload(tc.method, url, tc.body)
			require.Nil(t, err)
			request.Header.Add("Authorization", "bearer "+setup.Token)

			resp, err := app.Test(request)
			require.Nil(t, err)
			tc.checkResponse(t, resp)

			checkCate, err := setup.CateRepo.GetCategoryById(context.Background(), cate.ID)
			require.Nil(t, err)

			require.Equal(t, checkCate.Title, tc.body["title"].(string))
			for index := range checkCate.Subcategory {
				require.Equal(t, checkCate.Subcategory[index].ID.Hex(), tc.body["subcategory"].([]string)[index])
			}
		})
	}

	setup.DeleteMockCategory(t, cate)
	setup.DeleteMockSubcategory(t, subcate1)
	setup.DeleteMockSubcategory(t, subcate2)

}
