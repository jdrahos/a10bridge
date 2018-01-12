package v2

import (
	"a10bridge/model"
	"fmt"
)

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

type baseRequest struct {
	A10URL, SessionID string
}

type result struct {
	Error  a10Error `json:"err"`
	Status string   `json:"status"`
}

type baseResponse struct {
	Result result `json:"response"`
}

type nameRequest struct {
	Base baseRequest
	Name string `json:"name"`
}

type loginRequest struct {
	A10URL, A10User, A10Pwd string
}
type loginResponse struct {
	SessionID string `json:"session_id"`
}

type logoutResponse = baseResponse
type logoutRequest = baseRequest

type serverRequest struct {
	Base   baseRequest
	Server *model.Server
}

type getServerRequest = nameRequest
type getServerResponse struct {
	Result result `json:"response"`
	Server struct {
		Name   string `json:"name"`
		IP     string `json:"host"`
		Weight int    `json:"weight"`
	} `json:"server"`
}

type createServerRequest = serverRequest
type createServerResponse struct {
	Result result `json:"response"`
}

type updateServerRequest = serverRequest
type updateServerResponse struct {
	Result result `json:"response"`
}

type monitorRequest struct {
	Base    baseRequest
	Monitor *model.HealthCheck
}

type getMonitorRequest = nameRequest
type getMonitorResponse struct {
	Result  result `json:"response"`
	Monitor struct {
		Name string `json:"name"`
	} `json:"healthMonitor"`
}

type createMonitorRequest = monitorRequest
type createMonitorResponse struct {
	Result result `json:"response"`
}

type updateMonitorRequest = monitorRequest
type updateMonitorResponse struct {
	Result result `json:"response"`
}
