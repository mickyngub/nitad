package project

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/stretchr/testify/require"
)

func TestAddProjectService(t *testing.T) {
	setup.NewTestApp(t)

	subcate := setup.AddMockSubcategory(t)
	cate := setup.AddMockCategory(t, subcate)

	setup.AddMockProject(t, cate)
	//TODO: add from http request

}

type Response struct {
	Result project.Project
}

func TestGetProjectById(t *testing.T) {
	app := setup.NewTestApp(t)

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

	response := new(Response)
	json.Unmarshal(bodyBytes, response)

	require.Equal(t, proj.ID, response.Result.ID, "The project Id should be equal")
	setup.DeleteMock(t, proj, cate, subcate)
}
