package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFailure(t *testing.T) {
	status := 200
	errorMessage := "error_msg"
	failure := CreateFailure(status, errorMessage)

	assert.Equal(t, failure.Status(), status)
	assert.Equal(t, failure.Message(), errorMessage)
}

func TestFailureParseJSON(t *testing.T) {
	status := 200
	errorMessage := "error_msg"
	failure := CreateFailure(status, errorMessage)

	failureJSON, err := failure.ToJSON()
	assert.Nil(t, err)

	failure2, err := CreateFailureFromJSON(failureJSON)
	assert.Nil(t, err)

	assert.Equal(t, failure.Status(), failure2.Status())
	assert.Equal(t, failure.Message(), failure2.Message())
}
