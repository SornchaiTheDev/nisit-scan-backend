package libs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validatorImpl struct {
	validate *validator.Validate
}

type ValidateError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

var Validator = NewValidator()

func NewValidator() *validatorImpl {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &validatorImpl{
		validate: validate,
	}
}

func (v *validatorImpl) Validate(data interface{}) *ValidateError {
	errs := v.validate.Struct(data)

	fields := []string{}

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			fields = append(fields, fmt.Sprintf("field %s is %s", err.Field(), err.Tag()))
		}
		if len(errs.(validator.ValidationErrors)) > 0 {
			return &ValidateError{Message: "These fields are required",
				Errors: fields,
			}
		}
	}
	return nil
}
