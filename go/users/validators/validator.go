package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validate = *validator.New()

func ValidateStruct(data interface{}) []ErrorResponse {
	errs := validate.Struct(data)
	validatorErrors := []ErrorResponse{}

	fmt.Println(validatorErrors)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validatorErrors = append(validatorErrors, elem)
		}
		return validatorErrors
	}
	return nil
}
