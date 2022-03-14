package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/stretchr/testify/require"
)

// func TestAddProjectService(t *testing.T) {
// 	setup.NewTestApp(t)

// 	subcate := setup.AddMockSubcategory(t)
// 	cate := setup.AddMockCategory(t, subcate)

// 	setup.AddMockProject(t, cate)

// 	var client *http.Client
// 	URL := "/api/v1/project"
// 	err := Upload(client, URL, values)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
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
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

type Response struct {
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

	response := new(Response)

	json.Unmarshal(bodyBytes, response)

	require.Equal(t, proj.ID, response.Result.ID, "The project Id should be equal")

	setup.DeleteMock(t, proj, cate, subcate)
}
