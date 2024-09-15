package libs

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

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

	validate.RegisterValidation("date", dateValidator)
	validate.RegisterValidation("timestamp", timestampValidator)

	return &validatorImpl{
		validate: validate,
	}
}

func dateValidator(fl validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^(?:\d{2})/(?:\d{2})/(?:\d{4})$`, fl.Field().String())
	if err != nil {
		return false
	}
	return match
}

func timestampValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
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
