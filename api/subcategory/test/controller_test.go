package subcategory_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"github.com/birdglove2/nitad-backend/utils"
	"github.com/go-playground/assert/v2"
)

// type Subcategory struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty`
// 	Title     string             `bson:"title,omitempty`
// 	Image     string             `bson:"image,omitempty`
// 	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
// 	updatedAt time.Time          `bson:"updated_at" json:"updated_at"`
// }

// type SubcategoryResponse struct {
// 	success bool
// 	result  Subcategory
// }

// func TestDoStuffWithTestServer(t *testing.T) {

// 	config.Loadenv()
// 	database.ConnectDb()
// 	app := config.InitApp()
// 	api.CreateAPI(app)

// 	req := httptest.NewRequest("GET", "/subcategory", nil)
// 	req.Header.Set("X-Custom-Header", "hi")

// 	resp, _ := app.Test(req)
// 	bodyByte, _ := ioutil.ReadAll(resp.Body)

// 	var jsonMap map[string]interface{}
// 	json.Unmarshal(bodyByte, &jsonMap)
// 	fmt.Println("this is getting method", jsonMap)
// }

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

func TestListSubcategory(t *testing.T) {

	url := "http://localhost:3000/api/v1/subcategory"
	resp, err := http.Get(url)

	bodyByte, _ := ioutil.ReadAll(resp.Body)

	var jsonMap map[string]interface{}
	json.Unmarshal(bodyByte, &jsonMap)
	log.Println(jsonMap["result"])

	assert.Equal(t, err, nil)
	assert.Equal(t, true, jsonMap["success"])
	assert.Equal(t, "dummy subcategory title", jsonMap["result"].(map[string]interface{})["title"])
	assert.Equal(t, "dummy subcategory image", jsonMap["result"].(map[string]interface{})["image"])
}

func TestAddSubcategory(t *testing.T) {
	form := url.Values{}
	form.Add("title", "dummy subcategory title")
	form.Add("image", "randomImages(2)[0]")

	URL := "http://localhost:3000/api/v1/subcategory"
	resp, _ := http.PostForm(URL, form)
	bodyByte, _ := ioutil.ReadAll(resp.Body)

	var jsonMap map[string]interface{}
	json.Unmarshal(bodyByte, &jsonMap)
	log.Println(jsonMap["result"])

	assert.Equal(t, true, jsonMap["success"])
	assert.Equal(t, "dummy subcategory title", jsonMap["result"].(map[string]interface{})["title"])
	assert.Equal(t, "dummy subcategory image", jsonMap["result"].(map[string]interface{})["image"])

}

// func TestGetSubcategoryById(t *testing.T) {

// }

// it should add a subcategory successfully
// it should not add a subcategory if required fields are not provided correctly

// it should list all added subcategories

// it should get the subcategory by valid id
