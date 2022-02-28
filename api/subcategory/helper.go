package subcategory

import (
	"github.com/birdglove2/nitad-backend/api/collections_helper"
	"github.com/birdglove2/nitad-backend/errors"

	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// receive array of subcategoryIds, then
// find and return non-duplicated subcategories, and their ids
func FindByIds(sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var subcategories []Subcategory

	sids = utils.RemoveDuplicateIds(sids)

	for _, sid := range sids {
		oid, err := utils.IsValidObjectId(sid)
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
	objectId, err := utils.IsValidObjectId(sid)
	if err != nil {
		return s, err
	}

	if s, err = GetById(objectId); err != nil {
		return s, err
	}

	return s, nil
}

func HandleUpdateImage(c *fiber.Ctx, s *Subcategory) (*Subcategory, errors.CustomError) {
	oldSubcategory, err := GetById(s.ID)
	if err != nil {
		return s, err
	}

	files, err := utils.ExtractUpdatedFiles(c, "image")
	if err != nil {
		return s, err
	}

	s.Image = oldSubcategory.Image
	// if there is file passed, delete the old one and upload a new one
	if len(files) > 0 {
		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(c.Context(), files[0], s.Image, collectionName)
		if err != nil {
			return s, err
		}
		// if upload success, pass the url to the subcategory struct
		s.Image = newUploadFilename
	}

	s.CreatedAt = oldSubcategory.CreatedAt
	return s, nil
}

func BsonToSubcategory(b bson.M) Subcategory {
	// convert bson to Subcategory struct
	var s Subcategory
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
