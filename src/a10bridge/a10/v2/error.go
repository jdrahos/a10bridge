package v2

import (
	"a10bridge/a10/api"
	"fmt"
)

type a10Error struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"msg"`
}

func buildA10Error(err error) api.A10Error {
	if err == nil {
		return nil
	}
	return a10Error{
		ErrorCode:    0,
		ErrorMessage: err.Error(),
	}
}

func (err a10Error) Code() int {
	return err.ErrorCode
}

func (err a10Error) Message() string {
	return err.ErrorMessage
}

func (err a10Error) Error() string {
	return fmt.Sprintf("%d - %s", err.ErrorCode, err.ErrorMessage)
}
