package subcategory

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// validate requested string of subcategoryIds
// and return valid []objectId, otherwise error
func ValidateIds(sids []string) ([]primitive.ObjectID, errors.CustomError) {
	objectIds := make([]primitive.ObjectID, len(sids))

	for i, sid := range sids {
		objectId, err := ValidateId(sid)
		if err != nil {
			return objectIds, err
		}

		objectIds[i] = objectId
	}

	return objectIds, nil
}

// validate requested string of a single subcategoryId
// and return valid objectId, otherwise error
func ValidateId(sid string) (primitive.ObjectID, errors.CustomError) {
	objectId, err := functions.IsValidObjectId(sid)
	if err != nil {
		return objectId, err

	}
	// if err != nil >> id is not found
	if _, err = FindById(objectId); err != nil {
		return objectId, err
	}

	return objectId, nil
}

func HandleUpdateImage(c *fiber.Ctx, p *Subcategory, oid primitive.ObjectID) (*Subcategory, errors.CustomError) {
	oldSubcategory, err := FindById(oid)
	if err != nil {
		return p, err
	}

	s := BsonToSubcategory(oldSubcategory)

	files, err := functions.ExtractUpdatedFiles(c, "image")

	if err != nil {
		return p, err
	}
	if files == nil {
		// no file passed, use old image url
		p.Image = s.Image
	} else {
		// delete old files
		defer gcp.DeleteImages([]string{s.Image}, collectionName)

		// upload new files
		imageURLs, err := gcp.UploadImages(files, collectionName)
		if err != nil {
			// if upload error, delete uploaded file if it was uploaed
			defer gcp.DeleteImages(imageURLs, collectionName)
			return p, err
		}

		// if upload success, pass the url to the subcategory struct
		p.Image = imageURLs[0]
	}

	return p, nil
}

func HandleDeleteImage(oid primitive.ObjectID) errors.CustomError {
	oldSubcategory, err := FindById(oid)
	if err != nil {
		return err
	}

	s := BsonToSubcategory(oldSubcategory)

	err = gcp.DeleteImages([]string{s.Image}, collectionName)
	if err != nil {
		return err
	}
	return nil
}

func BsonToSubcategory(b bson.M) Subcategory {
	// convert bson to subcategory
	var s Subcategory
	bsonBytes, _ := bson.Marshal(b)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
