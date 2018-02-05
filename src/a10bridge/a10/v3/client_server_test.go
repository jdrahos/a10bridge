package v3_test

import (
	"a10bridge/a10/api"
	"a10bridge/model"
	"a10bridge/testing"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func testUpdateServer(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPut).
		Path("/axapi/v3/slb/server/"+node.A10Server).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "server": {
    "name": "`+node.A10Server+`",
    "host": "`+node.IPAddress+`",
    "action": "enable",
    "weight": `+node.Weight+`,
    "conn-limit": 8000000
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.UpdateServer(&node)
	assert.Nil(err, "Unexpected error when updating server")
}

func testUpdateServer_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.UpdateServer(&node)
	assert.NotNil(err, "Expected error when update server call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testUpdateServer_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.UpdateServer(&node)
	assert.NotNil(err, "Expected error when update server call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testCreateServer(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/server/").
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "server": {
    "name": "`+node.A10Server+`",
    "host": "`+node.IPAddress+`",
    "action": "enable",
    "weight": `+node.Weight+`,
    "conn-limit": 8000000
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.CreateServer(&node)
	assert.Nil(err, "Unexpected error when creating server")
}

func testCreateServer_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.CreateServer(&node)
	assert.NotNil(err, "Expected error when create server call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testCreateServer_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	node := model.Node{
		A10Server: "server",
		IPAddress: "10.10.10.11",
		Weight:    "1",
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.CreateServer(&node)
	assert.NotNil(err, "Expected error when create server call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testGetServer(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	serverName := "server"
	weight := "1"
	ipAddress := "10.201.14.203"

	testServer.Reset().
		AddRequest().
		Method(http.MethodGet).
		Path("/axapi/v3/slb/server/"+serverName).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Response().
		Body(`{"server":{"name":"`+serverName+`","host":"`+ipAddress+`","gslb_external_address":"0.0.0.0","weight":`+weight+`,"health_monitor":"(default)","status":1,"conn_limit":8000000,"conn_limit_log":1,"conn_resume":0,"stats_data":1,"extended_stats":0,"slow_start":0,"spoofing_cache":0,"template":"default","port_list":[{"port_num":81,"protocol":2,"status":1,"weight":1,"no_ssl":0,"conn_limit":8000000,"conn_limit_log":0,"conn_resume":0,"template":"default","stats_data":1,"health_monitor":"(default)","extended_stats":0},{"port_num":90,"protocol":2,"status":1,"weight":1,"no_ssl":0,"conn_limit":8000000,"conn_limit_log":1,"conn_resume":0,"template":"default","stats_data":1,"health_monitor":"(default)","extended_stats":0}]}}`,
			"application/json")

	node, err := client.GetServer(serverName)

	assert.Nil(err, "Unexpected error when closing client session")
	assert.NotNil(node, "Expected node instance")
	assert.Equal("", node.Name)
	assert.Equal(serverName, node.A10Server)
	assert.Equal(weight, node.Weight)
	assert.Equal(ipAddress, node.IPAddress)
	assert.Equal(0, len(node.Labels))
}

func testGetServer_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	_, err := client.GetServer("doesn't matter")
	assert.NotNil(err, "Expected error when get server call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testGetServer_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	_, err := client.GetServer("doesn't matter")
	assert.NotNil(err, "Expected error when get server call fails in a10")
	assert.Equal(errorCode, err.Code())
}
