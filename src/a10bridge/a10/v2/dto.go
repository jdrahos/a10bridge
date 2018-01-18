package v2

import (
	"a10bridge/model"
)

type baseRequest struct {
	A10URL, SessionID string
}

type result struct {
	Error  a10Error `json:"err"`
	Status string   `json:"status"`
}

type simpleResponse struct {
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
	Result    result `json:"response"`
	SessionID string `json:"session_id"`
}

type logoutResponse = simpleResponse
type logoutRequest = baseRequest

type serverRequest struct {
	Base   baseRequest
	Server *model.Node
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
type createServerResponse = simpleResponse

type updateServerRequest = serverRequest
type updateServerResponse = simpleResponse

type monitorRequest struct {
	Base    baseRequest
	Monitor *model.HealthCheck
}

type getMonitorRequest = nameRequest
type getMonitorResponse struct {
	Result  result `json:"response"`
	Monitor struct {
		Name                      string `json:"name"`
		RetryCount                int    `json:"retry"`
		RequiredConsecutivePasses int    `json:"consec_pass_reqd"`
		Interval                  int    `json:"interval"`
		Timeout                   int    `json:"timeout"`
		HTTP                      struct {
			Endpoint   string `json:"url"`
			Port       int    `json:"port"`
			ExpectCode string `json:"expect_code"`
		} `json:"http"`
	} `json:"health_monitor"`
}

type createMonitorRequest = monitorRequest
type createMonitorResponse = simpleResponse

type updateMonitorRequest = monitorRequest
type updateMonitorResponse = simpleResponse

type serviceGroupRequest struct {
	Base         baseRequest
	ServiceGroup *model.ServiceGroup
}

type getServiceGroupRequest = nameRequest
type getServiceGroupResponse struct {
	Result       result `json:"response"`
	ServiceGroup struct {
		Name              string `json:"name"`
		HealthMonitorName string `json:"health_monitor"`
		Members           []struct {
			ServerName string `json:"server"`
			Port       int    `json:"port"`
		} `json:"member_list"`
	} `json:"service_group"`
}

type createServiceGroupRequest = serviceGroupRequest
type createServiceGroupResponse = simpleResponse

type updateServiceGroupRequest = serviceGroupRequest
type updateServiceGroupResponse = simpleResponse

type serviceGroupMemberRequest struct {
	Base   baseRequest
	Member *model.Member
}

type createServiceGroupMemberRequest = serviceGroupMemberRequest
type createServiceGroupMemberResponse = simpleResponse

type deleteServiceGroupMemberRequest = serviceGroupMemberRequest
type deleteServiceGroupMemberResponse = simpleResponse
