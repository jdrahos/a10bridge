package v2_test

import (
	"a10bridge/a10/v2"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildError(t *testing.T) {
	assert := assert.New(t)
	expectedMessage := "test"
	expectedCode := 0
	err := errors.New(expectedMessage)
	a10err := v2.TestHelper{}.BuildError(err)
	assert.Equal(expectedCode, a10err.Code(), "Error code needs to be 0, this should not be used for building a10 errors")
	assert.Equal(expectedMessage, a10err.Message())
	assert.Equal(fmt.Sprintf("%d - %s", expectedCode, expectedMessage), a10err.Error())
}

func TestBuildError_null(t *testing.T) {
	var err error = nil
	a10err := v2.TestHelper{}.BuildError(err)
	assert.Nil(t, a10err)
}
