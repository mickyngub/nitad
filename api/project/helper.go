package project

import (
	"log"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLookupStage() mongo.Pipeline {
	pipe := mongo.Pipeline{}
	// pipe = database.AppendLookupStage(pipe, "category")
	// pipe = database.AppendLookupStage(pipe, "category.subcategory")
	// pipe = database.AppendUnwindStage(pipe, "category")

	return pipe
}

var SORTING = map[string]string{
	"views":     "views",
	"name":      "title",
	"updatedAt": "updatedAt",
	"createdAt": "createdAt",
}

// this will not sort updatedAt and createdAt
func AppendQueryStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	pq = SetDefaultQuery(pq)

	pipe = AppendSortStage(pipe, pq)
	pipe = AppendPaginationStage(pipe, pq)
	// return append(pipe, bson.D{{Key: "$sort", Value: bson.D{
	// 	{Key: SORTING[pq.Sort], Value: pq.By},
	// 	{Key: "title", Value: pq.ByName},
	// 	{Key: "updatedAt", Value: pq.ByUpdatedAt},
	// 	{Key: "createdAt", Value: pq.ByCreatedAt},
	// 	{Key: pq}
	// }}})
	return pipe
}

func AppendSortStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe, bson.D{
		{Key: "$sort", Value: bson.D{{
			Key: SORTING[pq.Sort], Value: pq.By,
		}}}})
}

func AppendPaginationStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	return append(pipe,
		bson.D{{Key: "$skip", Value: (pq.Page - 1) * pq.Limit}},
		bson.D{{Key: "$limit", Value: pq.Limit}})
}

func SetDefaultQuery(pq *ProjectQuery) *ProjectQuery {
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

	if pq.Page == 0 {
		pq.Page = 1
	}

	if pq.Limit == 0 {
		pq.Limit = 15
	}
	return pq
}

func IncrementViewCache(id string, views int) {
	key := "views" + id

	countInt := redis.GetCacheInt(key)
	if countInt != 0 {
		redis.SetCacheInt(key, countInt+1)
		return
	}

	redis.SetCacheInt(key, 1)
}

func IncrementView(id primitive.ObjectID, val int) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: val}}},
		},
	)

	if err != nil {
		log.Println("=====Incrementing view error: ", err.Error())
	}
}

func HandleDeleteImages(oid primitive.ObjectID) errors.CustomError {
	project, err := GetById(oid)
	if err != nil {
		return err
	}

	err = gcp.DeleteImages(project.Images, collectionName)
	if err != nil {
		log.Println("=====Delete images failed!!", project.Images, "======")
	}
	return nil
}

func HandleUpdateImages(c *fiber.Ctx, up *UpdateProject, oid primitive.ObjectID) (*UpdateProject, errors.CustomError) {
	oldProject, err := GetById(oid)
	if err != nil {
		return up, err
	}

	up.Images = oldProject.Images

	if len(up.DeleteImages) > 0 {
		// remove deleteImages from Images attrs
		up.Images = RemoveSliceFromSlice(up.Images, up.DeleteImages)
		err = gcp.DeleteImages(up.DeleteImages, collectionName)
		if err != nil {
			log.Println("=====Delete images failed!!", up.DeleteImages, "======")
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
