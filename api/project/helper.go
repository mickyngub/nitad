package project

import (
	"log"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLookupStage() mongo.Pipeline {
	pipe := mongo.Pipeline{}
	pipe = database.AppendLookupStage(pipe, "category")
	pipe = database.AppendUnwindStage(pipe, "category")

	pipe = database.AppendLookupStage(pipe, "subcategory")

	return pipe
}

var SORTING = map[string]string{
	"views":     "views",
	"name":      "title",
	"updatedAt": "updatedAt",
	"createdAt": "createdAt",
}

// this will not sort updatedAt and createdAt
func AppendSortStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	// pq = SetDefaultSort(pq)

	if pq.Sort == "" {
		pq.Sort = "views"
	}
	if pq.By == 0 {
		if pq.Sort == "name" {
			pq.By = 1
		} else {
			pq.By = -1
		}
	}

	return append(pipe, bson.D{{Key: "$sort", Value: bson.D{
		{Key: SORTING[pq.Sort], Value: pq.By},
		// {Key: "title", Value: pq.ByName},
		// {Key: "updatedAt", Value: pq.ByUpdatedAt},
		// {Key: "createdAt", Value: pq.ByCreatedAt},
		// {Key: pq}
	}}})
}

// func SetDefaultSort(pq *ProjectQuery) *ProjectQuery {
// 	if pq.ByViews == 0 {
// 		// sort for most view
// 		pq.ByViews = -1
// 	}
// 	if pq.ByName == 0 {
// 		// sort by alphabet
// 		pq.ByName = 1
// 	}
// 	if pq.ByUpdatedAt == 0 {
// 		pq.ByUpdatedAt = -1
// 	}

// 	if pq.ByCreatedAt == 0 {
// 		pq.ByCreatedAt = -1
// 	}

// 	return pq
// }

func IncrementView(id primitive.ObjectID) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}},
		},
	)

	if err != nil {
		log.Fatal("=====Incrementing view error: ", err.Error())
	}
}

func HandleDeleteImages(oid primitive.ObjectID) errors.CustomError {
	project, err := GetById(oid)
	if err != nil {
		return err
	}

	err = gcp.DeleteImages(project.Images, collectionName)
	if err != nil {
		return err
	}
	return nil
}

func HandleUpdateImages(c *fiber.Ctx, up *UpdateProject, oid primitive.ObjectID) (*UpdateProject, errors.CustomError) {
	oldProject, err := FindById(oid)
	if err != nil {
		return up, err
	}

	up.Images = oldProject.Images

	if len(up.DeleteImages) > 0 {
		// remove deleteImages from Images attrs
		oldProject.Images = RemoveSliceFromSlice(oldProject.Images, up.DeleteImages)
		err = gcp.DeleteImages(up.DeleteImages, collectionName)
		if err != nil {
			return up, err
		}
	}

	files, err := functions.ExtractUpdatedFiles(c, "images")
	if err != nil {
		return up, err
	}

	if len(files) > 0 {
		// if file pass, upload file
		imageURLs, err := gcp.UploadImages(files, collectionName)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			defer gcp.DeleteImages(imageURLs, collectionName)
			return up, err
		}

		// concat uploaded file to the existing ones
		up.Images = append(up.Images, imageURLs...)
	}

	up.CreatedAt = oldProject.CreatedAt
	return up, nil
}

func FindById(oid primitive.ObjectID) (ProjectForDecode, errors.CustomError) {
	b, err := database.GetElementById(oid, collectionName)
	if err != nil {
		return ProjectForDecode{}, err
	}
	// log.Println(b)

	return BsonToProjectForDecode(b), nil
}

func BsonToProjectForDecode(b interface{}) ProjectForDecode {
	// convert bson to project
	var p ProjectForDecode
	bsonBytes, err := bson.Marshal(b)
	if err != nil {
		log.Println("ERROR", err.Error())
	}
	err = bson.Unmarshal(bsonBytes, &p)
	if err != nil {
		log.Println("ERROR 2", err.Error())
	}
	//NOTE: these errors make p.Images cannot found sometime
	// resulting in false image management
	// log.Println(p.Images)
	return p
}

// this function remove the slice remove from the slice base
// ex:  base := []string{"test", "abc", "def", "ghi"}
//      remove := []string{"abc", "test"}
// return [def ghi]
// Used to remove deleteImages from ImageURLs
func RemoveSliceFromSlice(base []string, remove []string) []string {
	for i := 0; i < len(base); i++ {
		url := base[i]
		for _, rem := range remove {
			if url == rem {
				base = append(base[:i], base[i+1:]...)
				i--
				break
			}
		}
	}

	return base
}
