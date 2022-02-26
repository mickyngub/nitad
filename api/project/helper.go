package project

import (
	"context"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/redis"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

func AppendCountStage(pipe mongo.Pipeline) mongo.Pipeline {
	return append(pipe, bson.D{{Key: "$count", Value: "id"}})
}

func AppendQueryStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	pq = SetDefaultQuery(pq)

	pipe = AppendSortStage(pipe, pq)
	pipe = AppendPaginationStage(pipe, pq)
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
	// log.Println("9")
	countInt := redis.GetCacheInt(key)
	// log.Println("10")
	if countInt != 0 {
		redis.SetCacheInt(key, countInt+1)
		zap.S().Info("incrementing view of ", key, " = ", countInt+1)
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
		zap.S().Warn("Incrementing view error: ", err.Error())
	}
}

func HandleDeleteImages(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	project, err := GetById(oid)
	if err != nil {
		return err
	}

	err = gcp.DeleteImages(ctx, project.Images, collectionName)
	if err != nil {
		zap.S().Warn("Delete images failed", project.Images)
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
		err = gcp.DeleteImages(c.Context(), up.DeleteImages, collectionName)
		if err != nil {
			zap.S().Warn("Delete images failed", up.DeleteImages)
		}
	}

	files, err := utils.ExtractUpdatedFiles(c, "images")
	if err != nil {
		return up, err
	}

	if len(files) > 0 {
		// if file pass, upload file
		imageURLs, err := gcp.UploadImages(c.Context(), files, collectionName)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			gcp.DeleteImages(c.Context(), imageURLs, collectionName)
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
