package project

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type ProjectResponse struct {
	Result project.Project
}

func TestGetProjectById(t *testing.T) {
	app, _ := setup.NewTestApp(t)

	subcate := setup.AddMockSubcategory(t)

	cate := setup.AddMockCategory(t, subcate)

	proj := setup.AddMockProject(t, cate)

	url := "/api/v1/project/" + proj.ID.Hex()

	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.Nil(t, err)

	resp, err := app.Test(request)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	response := new(ProjectResponse)

	json.Unmarshal(bodyBytes, response)

	require.Equal(t, proj.ID, response.Result.ID, "The project Id should be equal")

	setup.DeleteMock(t, proj, cate, subcate)
}

func TestAddProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collectionName := "project"
	url := "/api/v1/project"

	subcate1 := setup.AddMockSubcategory(t)
	subcate2 := setup.AddMockSubcategory(t)
	cate1 := setup.AddMockCategory(t, subcate1)
	cate2 := setup.AddMockCategory(t, subcate2)

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
				"title":       "test add project",
				"description": "add project description",
				"authors":     "add project authors",
				"emails":      "add project emails",
				"inspiration": "add project inspiration",
				"abstract":    "add project abstract",
				"videos":      "add project videos",
				"keywords":    "add project keywords",
				"virtualLink": "add project virtualLink",
				"status":      "add project status",
				"images":      setup.OpenFileFromPath("dummy_image.jpg"),
				"report":      setup.OpenFileFromPath("dummy_report.pdf"),
				"category":    []string{cate1.ID.Hex(), cate2.ID.Hex()},
				"subcategory": []string{subcate1.ID.Hex(), subcate2.ID.Hex()},
			},
			buildMock: func(gcpService *setup.MockUploader, collectionName string) {
				gcpService.EXPECT().UploadFiles(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return([]string{"dummy image url"}, nil)
				gcpService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return("dummy report url", nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
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

			projResponse := new(ProjectResponse)
			bodyBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err)
			json.Unmarshal(bodyBytes, projResponse)

			proj, err := setup.ProjectRepo.GetProjectById(context.Background(), projResponse.Result.ID.Hex())
			require.Nil(t, err)

			require.Equal(t, tc.body["title"].(string), proj.Title)
			require.Equal(t, tc.body["description"].(string), proj.Description)
			require.Equal(t, []string{tc.body["authors"].(string)}, proj.Authors)
			require.Equal(t, []string{tc.body["emails"].(string)}, proj.Emails)
			require.Equal(t, tc.body["inspiration"].(string), proj.Inspiration)
			require.Equal(t, tc.body["abstract"].(string), proj.Abstract)
			require.Equal(t, []string{tc.body["videos"].(string)}, proj.Videos)
			require.Equal(t, []string{tc.body["keywords"].(string)}, proj.Keywords)
			require.Equal(t, tc.body["virtualLink"].(string), proj.VirtualLink)
			require.Equal(t, tc.body["status"].(string), proj.Status)
			require.Equal(t, proj.Images, []string{"dummy image url"})
			require.Equal(t, proj.Report, "dummy report url")
			require.Equal(t, 0, proj.Views)
			require.NotNil(t, proj.CreatedAt)
			require.NotNil(t, proj.UpdatedAt)

			require.Equal(t, cate1.ID.Hex(), proj.Category[0].ID.Hex())
			require.Equal(t, cate2.ID.Hex(), proj.Category[1].ID.Hex())

			require.Equal(t, subcate1.ID.Hex(), proj.Category[0].Subcategory[0].ID.Hex())
			require.Equal(t, subcate2.ID.Hex(), proj.Category[1].Subcategory[0].ID.Hex())

			setup.DeleteMockProject(t, proj)
		})
	}
	setup.DeleteMockCategory(t, cate1)
	setup.DeleteMockCategory(t, cate2)
	setup.DeleteMockSubcategory(t, subcate1)
	setup.DeleteMockSubcategory(t, subcate2)

}
