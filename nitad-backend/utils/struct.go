package utils

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/jinzhu/copier"
)

// CopyStruct use `copier` pkg to copy struct field
// it log error if occurred, and return error from `copier` pkg
func CopyStruct(from interface{}, to interface{}) errors.CustomError {
	if err := copier.Copy(to, from); err != nil {
		return errors.NewInternalServerError("Copy struct failed" + err.Error())
	}

	return nil
}
