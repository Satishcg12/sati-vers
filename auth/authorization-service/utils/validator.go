package utils

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		errorMessages := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			jsonTag := cv.getJSONTag(i, err.StructField())
			if jsonTag != "" {
				errorMessages[jsonTag] = cv.getErrorMsg(err)
			} else {
				errorMessages[err.Field()] = cv.getErrorMsg(err)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessages)
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

func (cv *CustomValidator) getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "email":
		return "Invalid email format"
	case "eqfield":
		return fe.Field() + " must be equal to " + fe.Param()
	default:
		return "Invalid value"
	}
}

func (cv *CustomValidator) getJSONTag(model interface{}, fieldName string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	field, found := modelType.FieldByName(fieldName)
	if !found {
		return ""
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	jsonTag = strings.Split(jsonTag, ",")[0]
	return jsonTag
}
