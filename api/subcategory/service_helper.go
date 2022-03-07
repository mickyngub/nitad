package subcategory

import (
	"github.com/birdglove2/nitad-backend/errors"

	"github.com/birdglove2/nitad-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// receive array of subcategoryIds, then
// find and return non-duplicated subcategories, and their ids
func FindByIds(sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var subcategories []Subcategory

	sids = utils.RemoveDuplicateIds(sids)

	// for _, sid := range sids {
	// oid, err := utils.IsValidObjectId(sid)
	// if err != nil {
	// 	return subcategories, objectIds, err
	// }

	// subcate, err := s.repository.GetById(ctx, oid)
	// if err != nil {
	// 	return subcategories, objectIds, err
	// }
	// objectIds = append(objectIds, oid)
	// subcategories = append(subcategories, *subcate)
	// }

	return subcategories, objectIds, nil
}
