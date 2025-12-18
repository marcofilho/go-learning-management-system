package dto_test

import (
	"errors"
	"testing"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	err := errors.New("test error")
	message := "An error occurred"
	errorResponse := dto.NewErrorResponse(err, message)

	assert.Equal(t, err.Error(), errorResponse.Error)
	assert.Equal(t, message, errorResponse.Message)
}
