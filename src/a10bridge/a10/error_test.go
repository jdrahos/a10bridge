package a10_test

import (
	"fmt"
	"testing"

	"a10bridge/a10"

	"github.com/stretchr/testify/assert"
)

func TestBuildError(t *testing.T) {
	expectedMessage := "test"
	expectedCode := 0
	err := a10.TestHelper{}.BuildError(expectedMessage)
	assert.Equal(t, expectedCode, err.Code(), "Error code needs to be 0, this should not be used for building a10 errors")
	assert.Equal(t, expectedMessage, err.Message())
	assert.Equal(t, fmt.Sprintf("%d - %s", expectedCode, expectedMessage), err.Error())
}
