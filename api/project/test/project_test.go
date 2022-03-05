package project_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/stretchr/testify/require"
)

// func TestAddProjectService(t *testing.T) {
// 	newTestApp(t)

// 	subcate := addMockSubcategory(t)
// 	cate := addMockCategory(t, subcate)

// 	addMockProject(t, cate)
// 	//TODO: add from http request

// }

type Response struct {
	Result project.Project
}

func TestGetProjectById(t *testing.T) {
	app := newTestApp(t)

	subcate := addMockSubcategory(t)

	cate := addMockCategory(t, subcate)
	proj := addMockProject(t, cate)

	url := "/api/v1/project/" + proj.ID.Hex()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)

	resp, err := app.Test(request)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	response := new(Response)
	json.Unmarshal(bodyBytes, response)

	require.Equal(t, proj.ID, response.Result.ID, "The project Id should be equal")
	deleteMock(t, proj, cate, subcate)
}
