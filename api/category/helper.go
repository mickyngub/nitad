package category

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// use for checking ids from ProjectRequest
// receive array of categoryIds, then
// find and return non-duplicated categories, and their ids
// return []Category
func (c *categoryService) FindByIds2(ctx context.Context, cids []string) ([]Category, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var categories []Category

	cids = utils.RemoveDuplicateIds(cids)

	for _, cid := range cids {
		category, err := c.GetCategoryById(ctx, cid)
		if err != nil {
			return categories, objectIds, err
		}

		objectIds = append(objectIds, category.ID)
		categories = append(categories, *category)

	}

	return categories, objectIds, nil

}

// TODO: this function is written in O(n^3), should find a better way to handle this later.
// merge multiple categories with multiple sids
// such that the finalCate will result in
// multiple categories that contain only relevant subcategories
// need to do this because the GetById of category return all subcategories
// that are in the category, so we need to filter some out
func FilterCatesWithSids(categories []Category, sids []primitive.ObjectID) ([]Category, errors.CustomError) {
	finalCate := []Category{}
	for _, cate := range categories {
		subcateThatIsInCate := []subcategory.Subcategory{}
		for _, subcate := range cate.Subcategory {
			for index, id := range sids {
				if subcate.ID == id {
					subcateThatIsInCate = append(subcateThatIsInCate, subcate)
					sids = remove(sids, index)
					index--
				}
			}
		}
		cate.Subcategory = subcateThatIsInCate
		finalCate = append(finalCate, cate)
	}

	if len(sids) > 0 {
		return finalCate, errors.NewBadRequestError("conflict: some subcategories are not in any categories")
	}
	return finalCate, nil
}

// remove the value at index in slice unordered
func remove(slice []primitive.ObjectID, i int) []primitive.ObjectID {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
