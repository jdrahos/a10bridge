package v3

import (
	"a10bridge/a10/api"
	"a10bridge/config"
	"a10bridge/model"
	"a10bridge/util"
	"strconv"
)

type v3Client struct {
	baseRequest   baseRequest
	commonHeaders map[string]string
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

func Connect(a10Instance *config.A10Instance) (api.Client, api.A10Error) {
	var client api.Client
	urltpl := "{{.A10URL}}/axapi/v3/auth"
	request := loginRequest{
		A10URL:  a10Instance.APIUrl,
		A10User: a10Instance.UserName,
		A10Pwd:  a10Instance.Password,
	}
	response := loginResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/auth.request", &request, &response, map[string]string{})
	if err != nil {
		return client, buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return client, response.Result.Error
	}

	client = v3Client{
		baseRequest: baseRequest{
			A10URL: a10Instance.APIUrl,
		},
		commonHeaders: map[string]string{
			"Authorization": "A10 " + response.Authresponse.Signature,
		},
	}

	return client, buildA10Error(err)
}

func (client v3Client) Close() api.A10Error {
	urltpl := "{{.A10URL}}/axapi/v3/logoff"
	request := client.baseRequest
	response := logoutResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/logout.request", &request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) GetServer(serverName string) (*model.Node, api.A10Error) {
	var server *model.Node
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/server/{{.Name}}"
	request := getServerRequest{
		Base: client.baseRequest,
		Name: serverName,
	}
	response := getServerResponse{}
	err := util.HttpGet(urltpl, request, &response, client.commonHeaders)
	if err != nil {
		return server, buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return server, response.Result.Error
	}

	server = &model.Node{
		A10Server: response.Server.Name,
		IPAddress: response.Server.IP,
		Weight:    strconv.Itoa(response.Server.Weight),
	}

	return server, nil
}

func (client v3Client) CreateServer(server *model.Node) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/server/"
	request := createServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}
	response := createServerResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/server.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) UpdateServer(server *model.Node) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/server/{{.Server.A10Server}}"
	request := updateServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}
	response := updateServerResponse{}
	err := util.HttpPut(urltpl, "a10/v3/tpl/server.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) GetHealthMonitor(monitorName string) (*model.HealthCheck, api.A10Error) {
	var monitor *model.HealthCheck
	urltpl := "{{.Base.A10URL}}/axapi/v3/health/monitor/{{.Name}}"
	request := getMonitorRequest{
		Base: client.baseRequest,
		Name: monitorName,
	}
	response := getMonitorResponse{}
	err := util.HttpGet(urltpl, request, &response, client.commonHeaders)
	if err != nil {
		return monitor, buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return monitor, response.Result.Error
	}

	mon := response.Monitor
	monitor = &model.HealthCheck{
		Name:                      mon.Name,
		Endpoint:                  mon.Method.HTTP.Endpoint,
		ExpectCode:                mon.Method.HTTP.ExpectCode,
		Port:                      mon.Method.HTTP.Port,
		Interval:                  mon.Interval,
		RetryCount:                mon.RetryCount,
		Timeout:                   mon.Timeout,
		RequiredConsecutivePasses: mon.RequiredConsecutivePasses,
	}

	return monitor, nil
}

func (client v3Client) CreateHealthMonitor(monitor *model.HealthCheck) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/health/monitor/"
	request := createMonitorRequest{
		Base:    client.baseRequest,
		Monitor: monitor,
	}
	response := createMonitorResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/health.monitor.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) UpdateHealthMonitor(monitor *model.HealthCheck) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/health/monitor/{{.Monitor.Name}}"
	request := updateMonitorRequest{
		Base:    client.baseRequest,
		Monitor: monitor,
	}
	response := updateMonitorResponse{}
	err := util.HttpPut(urltpl, "a10/v3/tpl/health.monitor.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) GetServiceGroup(serviceGroupName string) (*model.ServiceGroup, api.A10Error) {
	var serviceGroup *model.ServiceGroup
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/service-group/{{.Name}}"
	request := getServiceGroupRequest{
		Base: client.baseRequest,
		Name: serviceGroupName,
	}
	response := getServiceGroupResponse{}
	err := util.HttpGet(urltpl, request, &response, client.commonHeaders)
	if err != nil {
		return serviceGroup, buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return serviceGroup, response.Result.Error
	}

	sg := response.ServiceGroup
	serviceGroup = &model.ServiceGroup{
		Name: sg.Name,
		Health: &model.HealthCheck{
			Name: sg.HealthMonitorName,
		},
		Members: make([]*model.Member, len(sg.Members)),
	}

	for idx, member := range sg.Members {
		serviceGroup.Members[idx] = &model.Member{
			Port:             member.Port,
			ServerName:       member.ServerName,
			ServiceGroupName: serviceGroup.Name,
		}
	}

	return serviceGroup, nil
}

func (client v3Client) CreateServiceGroup(serviceGroup *model.ServiceGroup) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/service-group/"
	request := createServiceGroupRequest{
		Base:         client.baseRequest,
		ServiceGroup: serviceGroup,
	}
	response := createServiceGroupResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/svcgrp.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) UpdateServiceGroup(serviceGroup *model.ServiceGroup) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/service-group/{{.ServiceGroup.Name}}"
	request := updateServiceGroupRequest{
		Base:         client.baseRequest,
		ServiceGroup: serviceGroup,
	}
	response := updateServiceGroupResponse{}
	err := util.HttpPut(urltpl, "a10/v3/tpl/svcgrp.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) CreateMember(member *model.Member) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/service-group/{{.Member.ServiceGroupName}}/member/"
	request := createServiceGroupMemberRequest{
		Base:   client.baseRequest,
		Member: member,
	}
	response := createServiceGroupMemberResponse{}
	err := util.HttpPost(urltpl, "a10/v3/tpl/svcgrp.member.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) DeleteMember(member *model.Member) api.A10Error {
	urltpl := "{{.Base.A10URL}}/axapi/v3/slb/service-group/{{.Member.ServiceGroupName}}/member/{{.Member.ServerName}}+{{.Member.Port}}"
	request := deleteServiceGroupMemberRequest{
		Base:   client.baseRequest,
		Member: member,
	}
	response := deleteServiceGroupMemberResponse{}
	err := util.HttpDelete(urltpl, request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}
	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v3Client) IsServerNotFound(err api.A10Error) bool {
	return err.Code() == 1023460352
}

func (client v3Client) IsHealthMonitorNotFound(err api.A10Error) bool {
	return err.Code() == 1023460352
}

func (client v3Client) IsServiceGroupNotFound(err api.A10Error) bool {
	return err.Code() == 1023460352
}

func (client v3Client) IsMemberAlreadyExists(err api.A10Error) bool {
	return err.Code() == 1405
}
