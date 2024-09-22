package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type TestModel struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestMakeValidationError(t *testing.T) {
	// Create an instance of the validator
	v := validator.New()

	// Create a test case with invalid data
	model := TestModel{
		Name:  "",              // This field is required
		Email: "invalid-email", // This is not a valid email
	}

	// Validate the model
	err := v.Struct(model)

	// Check that the error is a validation error
	require.True(t, IsValidationError(err))

	// Generate validation errors
	validationErrors := MakeValidationError(err, model)

	// Assert the length of the validation errors
	require.Len(t, validationErrors, 2)

	// Assert specific validation errors
	require.Equal(t, "name", validationErrors[0].Field)
	require.Equal(t, "This field is required", validationErrors[0].Message)

	require.Equal(t, "email", validationErrors[1].Field)
	require.Equal(t, "Invalid email", validationErrors[1].Message)
}
