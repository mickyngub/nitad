package category_test

import (
	"fmt"
	"testing"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/stretchr/testify/require"
)

func createDummySubcategory(t *testing.T) *subcategory.Subcategory {
	fmt.Println("hhh 0")

	dummySubcate := subcategory.Subcategory{
		Title: "dummy subcate title",
		Image: "dummy subcate image url",
	}

	fmt.Println("hhh 1", &dummySubcate)

	adddedSubcategory, err := subcategory.Add(&dummySubcate)
	fmt.Println("hhh 2222", adddedSubcategory)
	if err != nil {
		fmt.Println("hhh err", err)
	}

	require.Equal(t, err, nil)
	fmt.Println("hhh err2")

	require.Equal(t, dummySubcate.Title, adddedSubcategory.Title)
	fmt.Println("hhh err3")

	require.Equal(t, dummySubcate.Image, adddedSubcategory.Image)
	require.NotEqual(t, nil, adddedSubcategory.ID)
	fmt.Println("hhh 2")

	return adddedSubcategory
}

// func newTestServer(t *testing.T, gcpService gcp.Uploader) *fiber.App {
// 	app := fiber.New()
// 	api.CreateAPI(app, gcpService)
// 	return app
// }

// func TestAddCategory(t *testing.T) {
// 	fmt.Println("Test Add Category 1")

// 	// fmt.Println("ALL env", ())
// 	config.Loadenv()
// 	mongo_URI := "mongodb+srv://nitad:od4POFf8eBwNpDiT@nitad.mvsvb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

// 	database.ConnectDb(mongo_URI)
// 	fmt.Println("Test Add Category 2")

// 	fmt.Println("Test Add Category 3")

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	gcpService := NewMockUploader(ctrl)

// 	server := newTestServer(t, gcpService)
// 	fmt.Println("Test Add Category 4")

// 	addedSubcategory := createDummySubcategory(t)
// 	fmt.Println("Test Add Category 5")

// 	bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDY0NzI3NDksImlkIjoiNjIxYTAwM2JiYWNiZDEyYTQxNDM3YjI5Iiwic3ViIjoiNzFiY2VjMmUtZTI3My00M2IwLTliMTMtNjNiMjYwNGI4NWQxIiwidXNlcm5hbWUiOiJuaXRhZGFkbWluMSJ9.wrkMDcW88oaOLWf33bLQ4JDdw8yuyLJKidDcdL-JvDg"
// 	testCases := []struct {
// 		name          string
// 		body          map[string]interface{}
// 		checkResponse func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "OK",
// 			body: map[string]interface{}{
// 				"title":       "test dummy category title",
// 				"subcategory": []string{addedSubcategory.ID.Hex()},
// 			},
// 			checkResponse: func(t *testing.T, resp *http.Response) {
// 				require.Equal(t, http.StatusOK, resp.StatusCode)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			URL := "/api/v1/category"

// 			form := url.Values{}
// 			form.Add("title", "tc.body.title")
// 			form.Add("subcategory", addedSubcategory.ID.Hex())

// 			data, _ := json.Marshal(tc.body)
// 			fmt.Println("This is data ", data)

// 			// request, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(data))
// 			request, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(form.Encode()))
// 			request.Header.Add("Authorization", bearerToken)
// 			if err != nil {
// 				zap.S().Fatal("error", err.Error())
// 			}

// 			fmt.Println("This is Request ", request)

// 			resp, err := server.Test(request)
// 			if err != nil {
// 				zap.S().Fatal("error", err.Error())
// 			}

// 			bodyBytes, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				zap.S().Fatal("error", err.Error())
// 			}
// 			bodyString := string(bodyBytes)
// 			fmt.Println("This is RESPONSE ", bodyString)
// 			require.Equal(t, err, nil)

// 			tc.checkResponse(t, resp)

// 		})
// 	}

// }
