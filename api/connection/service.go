package connection

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSubcategoryThatAreNotInAnyCategory() ([]subcategory.Subcategory, errors.CustomError) {
	returnSubcategory := []subcategory.Subcategory{}
	subcategories, err := subcategory.FindAll()
	if err != nil {
		return returnSubcategory, err
	}

	categories, err := category.FindAll()
	if err != nil {
		return returnSubcategory, err
	}

	sidsThatAreAlreadyIncate := make(map[primitive.ObjectID]bool)
	for _, category := range categories {
		for _, subcate := range category.Subcategory {
			sidsThatAreAlreadyIncate[subcate.ID] = true
		}
	}

	for _, subcategory := range subcategories {
		_, isIn := sidsThatAreAlreadyIncate[subcategory.ID]
		if !isIn {
			returnSubcategory = append(returnSubcategory, subcategory)
		}
	}
	return returnSubcategory, nil
}
