package v2

import (
	"a10bridge/a10/api"
	"a10bridge/config"
	"a10bridge/model"
	"a10bridge/util"
	"strconv"
	"strings"
)

type v2Client struct {
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
	urltpl := "{{.A10URL}}/services/rest/V2.1/?format=json&method=authenticate&username={{.A10User}}&password={{.A10Pwd}}"
	request := loginRequest{
		A10URL:  a10Instance.APIUrl,
		A10User: a10Instance.UserName,
		A10Pwd:  a10Instance.Password,
	}
	commonHeaders := map[string]string{}
	response := loginResponse{}
	err := util.HttpGet(urltpl, &request, &response, commonHeaders)
	if err != nil {
		return client, buildA10Error(err)
	}

	client = v2Client{
		baseRequest: baseRequest{
			A10URL:    a10Instance.APIUrl,
			SessionID: response.SessionID,
		},
		commonHeaders: commonHeaders,
	}

	return client, buildA10Error(err)
}

func (client v2Client) Close() api.A10Error {
	urltpl := "{{.A10URL}}/services/rest/V2.1/?format=json&method=session.close&session_id={{.SessionID}}"

	request := client.baseRequest
	response := logoutResponse{}
	err := util.HttpGet(urltpl, &request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v2Client) GetServer(serverName string) (*model.Node, api.A10Error) {
	var server *model.Node
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.search"
	request := getServerRequest{
		Base: client.baseRequest,
		Name: serverName,
	}
	response := getServerResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/name.request", request, &response, client.commonHeaders)
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

func (client v2Client) CreateServer(server *model.Node) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.create"
	request := createServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}
	response := createServerResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/server.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) UpdateServer(server *model.Node) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.update"
	request := updateServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}
	response := updateServerResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/server.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) GetHealthMonitor(monitorName string) (*model.HealthCheck, api.A10Error) {
	var monitor *model.HealthCheck
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.hm.search"
	request := getMonitorRequest{
		Base: client.baseRequest,
		Name: monitorName,
	}
	response := getMonitorResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/name.request", request, &response, client.commonHeaders)
	if err != nil {
		return monitor, buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return monitor, response.Result.Error
	}

	mon := response.Monitor
	monitor = &model.HealthCheck{
		Name:                      mon.Name,
		Endpoint:                  strings.TrimPrefix(mon.HTTP.Endpoint, "GET "),
		ExpectCode:                mon.HTTP.ExpectCode,
		Port:                      mon.HTTP.Port,
		Interval:                  mon.Interval,
		RetryCount:                mon.RetryCount,
		Timeout:                   mon.Timeout,
		RequiredConsecutivePasses: mon.RequiredConsecutivePasses,
	}

	return monitor, nil
}

func (client v2Client) CreateHealthMonitor(monitor *model.HealthCheck) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.hm.create"
	request := createMonitorRequest{
		Base:    client.baseRequest,
		Monitor: monitor,
	}
	response := createMonitorResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/health.monitor.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) UpdateHealthMonitor(monitor *model.HealthCheck) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.hm.update"
	request := updateMonitorRequest{
		Base:    client.baseRequest,
		Monitor: monitor,
	}
	response := updateMonitorResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/health.monitor.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) GetServiceGroup(serviceGroupName string) (*model.ServiceGroup, api.A10Error) {
	var serviceGroup *model.ServiceGroup
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.service_group.search"
	request := getServiceGroupRequest{
		Base: client.baseRequest,
		Name: serviceGroupName,
	}
	response := getServiceGroupResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/name.request", request, &response, client.commonHeaders)
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

func (client v2Client) CreateServiceGroup(serviceGroup *model.ServiceGroup) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.service_group.create"
	request := createServiceGroupRequest{
		Base:         client.baseRequest,
		ServiceGroup: serviceGroup,
	}
	response := createServiceGroupResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/svcgrp.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) UpdateServiceGroup(serviceGroup *model.ServiceGroup) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.service_group.update"
	request := updateServiceGroupRequest{
		Base:         client.baseRequest,
		ServiceGroup: serviceGroup,
	}
	response := updateServiceGroupResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/svcgrp.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) CreateMember(member *model.Member) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.service_group.member.create"
	request := createServiceGroupMemberRequest{
		Base:   client.baseRequest,
		Member: member,
	}
	response := createServiceGroupMemberResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/svcgrp.member.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) DeleteMember(member *model.Member) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.service_group.member.delete"
	request := deleteServiceGroupMemberRequest{
		Base:   client.baseRequest,
		Member: member,
	}
	response := deleteServiceGroupMemberResponse{}
	err := util.HttpPost(urltpl, "a10/v2/tpl/svcgrp.member.request", request, &response, client.commonHeaders)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) IsServerNotFound(err api.A10Error) bool {
	return err.Code() == 67174402
}

func (client v2Client) IsHealthMonitorNotFound(err api.A10Error) bool {
	return err.Code() == 33619968
}

func (client v2Client) IsServiceGroupNotFound(err api.A10Error) bool {
	return err.Code() == 67305473
}

func (client v2Client) IsMemberAlreadyExists(err api.A10Error) bool {
	return err.Code() == 1405
}
