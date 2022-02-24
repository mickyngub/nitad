package subcategory

import (
	"context"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// receive array of subcategoryIds, then
// find and return non-duplicated subcategories, and their ids
func FindByIds(sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var subcategories []Subcategory

	sids = functions.RemoveDuplicateIds(sids)

	for _, sid := range sids {
		oid, err := functions.IsValidObjectId(sid)
		if err != nil {
			return subcategories, objectIds, err
		}

		s, err := GetById(oid)
		if err != nil {
			return subcategories, objectIds, err
		}
		objectIds = append(objectIds, oid)
		subcategories = append(subcategories, s)
	}

	return subcategories, objectIds, nil
}

// validate requested string of a single subcategoryId
// and return valid objectId, otherwise error
func ValidateId(sid string) (Subcategory, errors.CustomError) {
	var s Subcategory
	objectId, err := functions.IsValidObjectId(sid)
	if err != nil {
		return s, err
	}

	if s, err = GetById(objectId); err != nil {
		return s, err
	}

	return s, nil
}

func HandleUpdateImage(c *fiber.Ctx, s *Subcategory, oid primitive.ObjectID) (*Subcategory, errors.CustomError) {
	oldSubcategory, err := GetById(oid)
	if err != nil {
		return s, err
	}

	files, err := functions.ExtractUpdatedFiles(c, "image")

	if err != nil {
		return s, err
	}
	if files == nil {
		// no file passed, use old image url
		s.Image = oldSubcategory.Image
	} else {
		// delete old files
		gcp.DeleteImages(c.Context(), []string{s.Image}, collectionName)

		// upload new files
		imageURLs, err := gcp.UploadImages(c.Context(), files, collectionName)
		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			gcp.DeleteImages(c.Context(), imageURLs, collectionName)
			return s, err
		}

		// if upload success, pass the url to the subcategory struct
		s.Image = imageURLs[0]
	}

	s.CreatedAt = oldSubcategory.CreatedAt
	return s, nil
}

func HandleDeleteImage(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	oldSubcategory, err := GetById(oid)
	if err != nil {
		return err
	}

	err = gcp.DeleteImages(ctx, []string{oldSubcategory.Image}, collectionName)
	if err != nil {
		return err
	}
	return nil
}

func BsonToSubcategory(b bson.M) Subcategory {
	// convert bson to Subcategory struct
	var s Subcategory
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
