package validators

import (
	"strings"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Request interface{}

func ValidateStruct(req Request) errors.CustomError {

	var errfields []string
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			failedFieldArr := strings.Split(err.StructNamespace(), ".")
			failedField := failedFieldArr[len(failedFieldArr)-1]
			errfields = append(errfields, failedField+" "+err.Tag())
		}
		return errors.NewInvalidInputError(errfields)
	}
	return nil

}
