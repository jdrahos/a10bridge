package v3

import "fmt"

type a10Error struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"msg"`
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
