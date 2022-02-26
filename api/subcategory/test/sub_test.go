package subcategory_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mongo-go/testdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// const API_PREFIX = "/api/v1"

// // var app fiber.App = app.
// var app = fiber.New()

// // var v1 fiber.Router = app.Group(API_PREFIX)
// var subcategoryRouter fiber.Router = app.Use("/api/v1/subcategory")
// var subcategoryController *subcategory.Controller = subcategory.NewController(subcategoryRouter)

// func TestX(t *testing.T) {
// 	main.main()
// }
// func TestC(t *testing.T) {
// 	app := fiber.New()

// 	// Create route with GET method for test
// 	// app.Get("/", subcategoryController.ListSubcategory)
// 	api.CreateAPI(app)
// 	req := httptest.NewRequest("GET", "/api/v1/subcategory", nil)

// 	// Perform the request plain with the app,
// 	// the second argument is a request latency
// 	// (set to -1 for no latency)
// 	resp, _ := app.Test(req, 1)
// 	assert.Equal(t, resp, 20)
// 	// assert.Equal(t, nil, err, "they should be equal")

// }

// func Test_List_Subcategory(t *testing.T) {

// 	app := config.InitApp()

// 	fmt.Println("Hi 1")

// 	resp, err := app.Test(httptest.NewRequest("GET", "/api/v1/subcategory", nil))
// 	// zap.S().Info(resp)

// 	fmt.Println("Hi 2", resp)
// 	// resp, _ := http.PostForm(URL, form)
// 	bodyByte, _ := ioutil.ReadAll(resp.Body)

// 	var jsonMap map[string]interface{}
// 	json.Unmarshal(bodyByte, &jsonMap)

// 	// assert.Equal(t, true, jsonMap["success"])
// 	assert.Equal(t, "dummy subcategory title", jsonMap["result"].(map[string]interface{})["title"])
// 	// utAssertEqual(t, nil, err)
// 	assert.Equal(t, nil, err, "they should be equal")

// }

// func MockAddSubcategory() subcategory.Subcategory {
// 	s := subcategory.Subcategory{
// 		Title: "Mock Title Subcategory",
// 		Image: "Mock Image",
// 	}
// 	subcate, _ := subcategory.Add(&s)
// 	return *subcate
// }

// func MockAddCategory(s subcategory.Subcategory) category.Category {
// 	c := category.Category{
// 		Title: "Mock Title Category",
// 	}
// 	cate, _ := category.Add(&c, []primitive.ObjectID{s.ID})
// 	return *cate
// }

var testDb *testdb.TestDB

func setup(t *testing.T) *mongo.Collection {
	if testDb == nil {
		testDb = testdb.NewTestDB("mongodb://localhost", "your_db", time.Duration(2)*time.Second)

		err := testDb.Connect()
		if err != nil {
			t.Fatal(err)
		}
	}

	coll, err := testDb.CreateRandomCollection(testdb.NoIndexes)
	if err != nil {
		t.Fatal(err)
	}

	return coll // random *mongo.Collection in "your_db"
}

func Test1(t *testing.T) {
	coll := setup(t)
	defer coll.Drop(context.Background())

	// Test queries using coll
}

func Test_Create_Project(t *testing.T) {
	_, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		zap.S().Fatal(err.Error())
	}
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
	// if err != nil {
	// 	zap.S().Fatal(err.Error())
	// }

	// // List databases
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	zap.S().Fatal(err.Error())
	// }
	// log.Println(databases)
	// config.Loadenv()
	// database.ConnectDb()
	// defer database.DisconnectDb()

	// s := subcategory.Subcategory{
	// 	Title: "Mock Title Subcategory",
	// 	Image: "Mock Image",
	// }
	// fmt.Println("check s ", s)
	// pointer := &s
	// fmt.Println("check s 2", pointer)
	// subcategory.Add(pointer)

	// assert.Equal(t, s.Title, subcate.Title)
	// assert.Equal(t, s.Image, subcate.Image)
	// assert.Equal(t, s.Image, subcate.Image)

	// mockAddCategory(subcate)

	// dummySlice := []string{"dummy slice 1", "dummy slice 2"}

	// p := project.Project{
	// 	Title:       "pr.Title",
	// 	Description: "pr.Description",
	// 	Authors:     dummySlice,
	// 	Emails:      dummySlice,
	// 	Inspiration: "pr.Inspiration",
	// 	Abstract:    "pr.Abstract",
	// 	Videos:      dummySlice,
	// 	Keywords:    dummySlice,
	// 	Status:      "pr.Status",
	// 	Category:    []category.Category{cate},
	// }

	// result, err := project.Add(&p)
	// assert.Equal(t, nil, err, "err should be nil")
	// fmt.Println("result", result)

}
