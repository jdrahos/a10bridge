package a10

import (
	"a10bridge/a10/api"
	"fmt"
)

type a10ErrorImpl struct {
	ErrorCode    int
	ErrorMessage string
}

func (err a10ErrorImpl) Code() int {
	return err.ErrorCode
}

func (err a10ErrorImpl) Message() string {
	return err.ErrorMessage
}

func (err a10ErrorImpl) Error() string {
	return fmt.Sprintf("%d - %s", err.ErrorCode, err.ErrorMessage)
}

//BuildError builds a10 error with code set to 0 to indicate the error is not returned by a10 itself
func buildError(message string) api.A10Error {
	return a10ErrorImpl{
		ErrorCode:    0,
		ErrorMessage: message,
	}
}
