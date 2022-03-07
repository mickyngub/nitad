package subcategory

import (
	"net/http"
	"testing"

	"github.com/birdglove2/nitad-backend/api/setup"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetSubcategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		checkResponse func(*testing.T, *http.Response)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ctrl := gomock.NewController(t)
			// defer ctrl.Finish()

			// config.Loadenv()
			// database.ConnectDb(os.Getenv("MONGO_URI"))

			// app := fiber.New()

			// gcpService := NewMockUploader(ctrl)

			// api.CreateAPI(app, gcpService)

			server := setup.NewTestApp(t)
			url := "/api/v1/subcategory"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				zap.S().Fatal("error", err.Error())
			}

			resp, err := server.Test(request)
			if err != nil {
				zap.S().Fatal("error", err.Error())
			}

			// bodyBytes, err := io.ReadAll(resp.Body)
			// if err != nil {
			// zap.S().Fatal("error", err.Error())
			// }
			// bodyString := string(bodyBytes)
			// fmt.Println("This is RESPONSE SUBCATE  ", bodyString)
			require.Equal(t, err, nil)

			tc.checkResponse(t, resp)

		})
	}
}

// func TestAddSubcategory(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// bs, err := ioutil.ReadFile("bear_test.jpg")
// 	// if err != nil {
// 	// 	fmt.Println("Error file:", err)
// 	// 	os.Exit(1)
// 	// }

// 	// images := randomImages(1)
// 	// fmt.Println("images", images)
// 	// image, err := getImageFromFilePath("bear_test")
// 	// if err != nil {
// 	// 	fmt.Println("get image error", err)
// 	// }
// 	// image, err := os.Open("bear_test.jpg")
// 	// if err != nil {
// 	// 	zap.S().Fatal("error", err.Error())
// 	// }
// 	// fmt.Println("image", image)
// 	testCases := []struct {
// 		name           string
// 		body           fiber.Map
// 		collectionName string
// 		buildMock      func(gcpService *MockUploader, collectionName string)
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "OK",
// 			body: fiber.Map{
// 				"title":       "user.Username",
// 				"subcategory": []string{"62251e85d5f281183909394a"},
// 			},
// 			collectionName: "category",
// 			buildMock: func(gcpService *MockUploader, collectionName string) {
// 				// gcpService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return("dummy image url", nil)
// 			},
// 			checkResponse: func(t *testing.T, resp *http.Response) {
// 				require.Equal(t, http.StatusOK, resp.StatusCode)
// 			},
// 		},
// 		// {
// 		// 	name: "FAIL",
// 		// 	buildMock: func(gcpService *MockUploader, collectionName string) {
// 		// 		gcpService.EXPECT().UploadFiles(gomock.Any(), gomock.Any(), gomock.Eq(collectionName)).Times(1).Return(nil, nil)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, resp *http.Response) {
// 		// 		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
// 		// 	},
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			gcpService := NewMockUploader(ctrl)
// 			server := newTestApp(t)
// 			URL := "/api/v1/category"

// 			tc.buildMock(gcpService, tc.collectionName)

// 			data, err := json.Marshal(tc.body)
// 			if err != nil {
// 				fmt.Println("This is data  err", err.Error())
// 			}
// 			request, err := http.NewRequest("POST", URL, bytes.NewReader(data))
// 			resp, err := server.Test(request)

// 			// form := url.Values{}
// 			// form.Add("title", "this is tile")
// 			// form.Add("subcategory", "62251e85d5f281183909394a")
// 			// form.Add("image", asString)
// 			// request, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
// 			// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 			if err != nil {
// 				zap.S().Fatal("error", err.Error())
// 			}
// 			// var b bytes.Buffer
// 			// w := multipart.NewWriter(&b)
// 			// request.Header.Add("Content-Type", w.FormDataContentType())
// 			// request.Header.Add("Content-Length", strconv.Itoa(len(data)))

// 			// Buffer the body

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

// // it should add a subcategory successfully

// // it should not add a subcategory if required fields are not provided correctly

// // it should list all added subcategories

// // it should get the subcategory by valid id
