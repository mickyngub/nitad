package subcategory_test

import (
	"image"
	"net/http"
	"os"
	"testing"

	"github.com/birdglove2/nitad-backend/api"
	"github.com/birdglove2/nitad-backend/config"
	"github.com/birdglove2/nitad-backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestApp(t *testing.T) *fiber.App {
	config.Loadenv()
	database.ConnectDb(os.Getenv("MONGO_URI"))

	app := fiber.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gcpService := NewMockUploader(ctrl)

	api.CreateAPI(app, gcpService)

	return app
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

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

			server := newTestApp(t)
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
// 	image, err := os.Open("bear_test.jpg")
// 	if err != nil {
// 		zap.S().Fatal("error", err.Error())
// 	}
// 	// fmt.Println("image", image)
// 	testCases := []struct {
// 		name           string
// 		body           map[string]interface{}
// 		collectionName string
// 		buildMock      func(gcpService *MockUploader, collectionName string)
// 		checkResponse  func(*testing.T, *http.Response)
// 	}{
// 		{
// 			name: "OK",
// 			body: map[string]interface{}{
// 				"title": "test dummy subcategory title",
// 				"image": image,
// 			},
// 			collectionName: "subcategory",
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
// 			server := newTestServer(t, gcpService)
// 			url := "/api/v1/subcategory"

// 			tc.buildMock(gcpService, tc.collectionName)

// 			data, _ := json.Marshal(tc.body)
// 			fmt.Println("This is data ", tc.body)

// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			if err != nil {
// 				zap.S().Fatal("error", err.Error())
// 			}
// 			request.Header.Add("Content-Type", "multipart/form-data")

// 			// Buffer the body

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

// // it should add a subcategory successfully

// // it should not add a subcategory if required fields are not provided correctly

// // it should list all added subcategories

// // it should get the subcategory by valid id
