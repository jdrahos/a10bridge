package v2

import (
	"a10bridge/a10/api"
	"a10bridge/args"
	"a10bridge/model"
	"a10bridge/util"
)

type v2Client struct {
	baseRequest baseRequest
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

func Connect(arguments args.Args) (api.Client, api.A10Error) {
	var client api.Client
	urltpl := "{{.A10URL}}/services/rest/V2.1/?format=json&method=authenticate&username={{.A10User}}&password={{.A10Pwd}}"

	request := loginRequest{
		A10URL:  arguments.A10URL,
		A10User: arguments.A10User,
		A10Pwd:  arguments.A10Pwd,
	}

	url, err := util.ApplyTemplate(request, "login", urltpl)
	if err != nil {
		return client, buildA10Error(err)
	}

	response := loginResponse{}
	err = util.HttpGet(url, &response)
	if err != nil {
		return client, buildA10Error(err)
	}

	client = v2Client{
		baseRequest: baseRequest{
			A10URL:    arguments.A10URL,
			SessionID: response.SessionID,
		},
	}

	return client, buildA10Error(err)
}

func (client v2Client) Close() api.A10Error {
	urltpl := "{{.A10URL}}/services/rest/V2.1/?format=json&method=session.close&session_id={{.SessionID}}"
	request := client.baseRequest
	url, err := util.ApplyTemplate(request, "logout", urltpl)
	if err != nil {
		return buildA10Error(err)
	}

	response := logoutResponse{}
	err = util.HttpGet(url, &response)
	if err != nil {
		return buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return response.Result.Error
	}

	return nil
}

func (client v2Client) GetServer(serverName string) (*model.Server, api.A10Error) {
	var server *model.Server
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.search"

	request := getServerRequest{
		Base: client.baseRequest,
		Name: serverName,
	}

	url, err := util.ApplyTemplate(request, "get server", urltpl)
	if err != nil {
		return server, buildA10Error(err)
	}

	response := getServerResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/name.request", request, &response)
	if err != nil {
		return server, buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return server, response.Result.Error
	}

	server = &model.Server{
		Name:   response.Server.Name,
		IP:     response.Server.IP,
		Weight: response.Server.Weight,
	}

	return server, nil
}

func (client v2Client) CreateServer(server *model.Server) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.create"

	request := createServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}

	url, err := util.ApplyTemplate(request, "create server", urltpl)
	if err != nil {
		return buildA10Error(err)
	}

	response := createServerResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/server.request", request, &response)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}

func (client v2Client) UpdateServer(server *model.Server) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.server.update"

	request := updateServerRequest{
		Base:   client.baseRequest,
		Server: server,
	}

	url, err := util.ApplyTemplate(request, "update server", urltpl)
	if err != nil {
		return buildA10Error(err)
	}

	response := updateServerResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/server.request", request, &response)
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

	url, err := util.ApplyTemplate(request, "get monitor", urltpl)
	if err != nil {
		return monitor, buildA10Error(err)
	}

	response := getMonitorResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/name.request", request, &response)
	if err != nil {
		return monitor, buildA10Error(err)
	}

	if response.Result.Status == "fail" {
		return monitor, response.Result.Error
	}

	monitor = &model.HealthCheck{
		Name: response.Monitor.Name,
	}

	return monitor, nil
}

func (client v2Client) CreateHealthMonitor(monitor *model.HealthCheck) api.A10Error {
	urltpl := "{{.Base.A10URL}}/services/rest/V2.1/?session_id={{.Base.SessionID}}&format=json&method=slb.hm.create"

	request := createMonitorRequest{
		Base:    client.baseRequest,
		Monitor: monitor,
	}

	url, err := util.ApplyTemplate(request, "monitor request", urltpl)
	if err != nil {
		return buildA10Error(err)
	}

	response := createMonitorResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/health.monitor.request", request, &response)
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

	url, err := util.ApplyTemplate(request, "monitor request", urltpl)
	if err != nil {
		return buildA10Error(err)
	}

	response := updateMonitorResponse{}
	err = util.HttpPost(url, "a10/v2/tpl/health.monitor.request", request, &response)
	if err != nil {
		return buildA10Error(err)
	}

	return nil
}
