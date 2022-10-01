package validation

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Path        string
	Description string
	Err         error
}

func NewValidationError(path string, description string, err error) *ValidationError {
	return &ValidationError{Path: path, Description: description, Err: err}
}

func (e *ValidationError) Error() string {
	errorMessage := []string{}
	errorMessage = append(errorMessage, fmt.Sprintf("File %s is invalid.", e.Path))
	errorMessage = append(errorMessage, e.Description)
	if e.Err != nil {
		errorMessage = append(errorMessage, e.Err.Error())
	}
	return strings.Join(errorMessage, " ")
}
