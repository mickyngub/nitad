package project

import (
	"log"

	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
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
	pipe = database.AppendLookupStage(pipe, "subcategory")
	pipe = database.AppendUnsetStage(pipe, "category.subcategory")
	return pipe
}

func AppendSortStage(pipe mongo.Pipeline, pq *ProjectQuery) mongo.Pipeline {
	pq = SetDefaultSort(pq)
	return append(pipe, bson.D{{Key: "$sort", Value: bson.D{
		{Key: "views", Value: pq.ByViews},
		{Key: "title", Value: pq.ByName},
		{Key: "updatedAt", Value: pq.ByUpdatedAt},
		{Key: "createdAt", Value: pq.ByCreatedAt},
	}}})
}

func SetDefaultSort(pq *ProjectQuery) *ProjectQuery {
	if pq.ByViews == 0 {
		// sort for most view
		pq.ByViews = -1
	}
	if pq.ByName == 0 {
		// sort by alphabet
		pq.ByName = 1
	}
	if pq.ByUpdatedAt == 0 {
		pq.ByUpdatedAt = -1
	}

	if pq.ByCreatedAt == 0 {
		pq.ByCreatedAt = -1
	}

	return pq
}

func IncrementView(id primitive.ObjectID) {
	projectCollection, ctx := database.GetCollection(collectionName)

	_, err := projectCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}},
		},
	)

	// NOTE: logging ??
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateAndRemoveDuplicateIds(sids []string, cids []string) ([]primitive.ObjectID, []primitive.ObjectID, errors.CustomError) {
	var subcategoryIds []primitive.ObjectID
	var categoryIds []primitive.ObjectID

	soids, err := subcategory.ValidateIds(sids)
	if err != nil {
		return subcategoryIds, categoryIds, err
	}

	coids, err := category.ValidateIds(cids)
	if err != nil {
		return subcategoryIds, categoryIds, err
	}

	subcategoryIds = functions.RemoveDuplicateObjectIds(soids)
	categoryIds = functions.RemoveDuplicateObjectIds(coids)

	return subcategoryIds, categoryIds, nil
}

func HandleDeleteImages(oid primitive.ObjectID) errors.CustomError {
	project, err := FindById(oid)
	if err != nil {
		return err
	}
	p := BsonToProject(project)

	err = gcp.DeleteImages(p.Images, collectionName)
	if err != nil {
		return err
	}
	return nil
}

func HandleUpdateImages(c *fiber.Ctx, upr *UpdateProjectRequest, oid primitive.ObjectID) (*UpdateProjectRequest, errors.CustomError) {
	oldProject, err := database.FindById(oid, collectionName)
	if err != nil {
		return upr, err
	}
	pr := BsonToProject(oldProject)

	if len(upr.DeleteImages) > 0 {
		pr.Images = RemoveSliceFromSlice(pr.Images, upr.DeleteImages)
		defer gcp.DeleteImages(upr.DeleteImages, collectionName)
	}

	files, err := functions.ExtractUpdatedFiles(c, "images")
	if err != nil {
		return upr, err
	}

	if files == nil {
		// no file passed, use old image url
		upr.Images = pr.Images
		// log.Println("file == nil", upr.Images, pr.Images)
		return upr, nil
	} else {
		// if file pass
		// log.Println("file passed", upr.Images, pr.Images)
		// upload file
		imageURLs, err := gcp.UploadImages(files, collectionName)

		// log.Println("imageUrls", imageURLs)

		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			defer gcp.DeleteImages(imageURLs, collectionName)
			// log.Println("file bugs", upr.Images, pr.Images)
			return upr, err
		}

		// concat uploaded file to the existing ones
		imageURLs = append(pr.Images, imageURLs...)
		// log.Println("check ", imageURLs, pr.Images)

		// if upload success, pass the url to the subcategory struct
		upr.Images = imageURLs
		// log.Println("file latest", upr.Images)
	}

	return upr, nil
}

func BsonToProject(b bson.M) ProjectRequest {
	// convert bson to subcategory
	var pr ProjectRequest
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &pr)
	return pr
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
