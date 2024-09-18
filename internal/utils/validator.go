package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func IsValidationError(err error) bool {
	var ve validator.ValidationErrors
	return errors.As(err, &ve)
}

func MakeValidationError(err error, model interface{}) []validationError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]validationError, len(ve))
		for i, fe := range ve {
			fieldName := getJSONTag(model, fe.StructField())
			out[i] = validationError{fieldName, msgForTag(fe.Tag())}
		}
		return out
	}
	return []validationError{}
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func getJSONTag(model interface{}, fieldName string) string {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	// Check if the model is a pointer and dereference it if needed
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Look for the field by name
	field, found := modelType.FieldByName(fieldName)
	if !found {
		return fieldName // fallback to field name if no json tag found
	}

	// Get the JSON tag
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return fieldName // fallback if there's no json tag or it's ignored
	}

	return jsonTag
}
