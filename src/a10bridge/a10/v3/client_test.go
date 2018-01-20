package v3_test

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v3"
	"a10bridge/config"
	"a10bridge/testing"
	"errors"
	tst "testing"

	"github.com/stretchr/testify/assert"
)

var helper v3.TestHelper = v3.TestHelper{}

func TestSessionResource(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	testConnect(testServer, assert)
	testConnect_ServerError(testServer, assert)
	testConnect_failedAuthentication(testServer, assert)

	testClose(testServer, assert, client)
	testClose_ServerError(testServer, assert, client)
	testClose_Failure(testServer, assert, client)
}

func TestServerResource(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	testGetServer(testServer, assert, client)
	testGetServer_ServerError(testServer, assert, client)
	testGetServer_Failure(testServer, assert, client)

	testCreateServer(testServer, assert, client)
	testCreateServer_ServerError(testServer, assert, client)
	testCreateServer_Failure(testServer, assert, client)

	testUpdateServer(testServer, assert, client)
	testUpdateServer_ServerError(testServer, assert, client)
	testUpdateServer_Failure(testServer, assert, client)
}

func TestHealthMonitorResource(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	testGetMonitor(testServer, assert, client)
	testGetMonitor_ServerError(testServer, assert, client)
	testGetMonitor_Failure(testServer, assert, client)

	testCreateMonitor(testServer, assert, client)
	testCreateMonitor_ServerError(testServer, assert, client)
	testCreateMonitor_Failure(testServer, assert, client)

	testUpdateMonitor(testServer, assert, client)
	testUpdateMonitor_ServerError(testServer, assert, client)
	testUpdateMonitor_Failure(testServer, assert, client)
}

func TestServiceGroupResource(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	testGetServiceGroup(testServer, assert, client)
	testGetServiceGroup_ServerError(testServer, assert, client)
	testGetServiceGroup_Failure(testServer, assert, client)

	testCreateServiceGroup(testServer, assert, client)
	testCreateServiceGroup_ServerError(testServer, assert, client)
	testCreateServiceGroup_Failure(testServer, assert, client)

	testUpdateServiceGroup(testServer, assert, client)
	testUpdateServiceGroup_ServerError(testServer, assert, client)
	testUpdateServiceGroup_Failure(testServer, assert, client)
}

func TestServiceGroupMemberResource(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	testCreateMember(testServer, assert, client)
	testCreateMember_ServerError(testServer, assert, client)
	testCreateMember_Failure(testServer, assert, client)

	testDeleteMember(testServer, assert, client)
	testDeleteMember_ServerError(testServer, assert, client)
	testDeleteMember_Failure(testServer, assert, client)
}

func TestErrorCodeDetection(t *tst.T) {
	sessionId := "test_session_id"
	assert := assert.New(t)
	testServer := testing.NewTestServer(t).Start()
	defer testServer.Stop()

	client, err := buildClient(testServer, sessionId)
	assert.Nil(err, "Failed to build client for testing")

	a10err := helper.BuildError(errors.New("test"))
	assert.False(client.IsServerNotFound(a10err))
	assert.False(client.IsHealthMonitorNotFound(a10err))
	assert.False(client.IsServiceGroupNotFound(a10err))

	a10err = helper.SetErrorCode(a10err, 1023460352)
	assert.True(client.IsServerNotFound(a10err))
	assert.True(client.IsHealthMonitorNotFound(a10err))
	assert.True(client.IsServiceGroupNotFound(a10err))

	assert.False(client.IsMemberAlreadyExists(a10err))
	a10err = helper.SetErrorCode(a10err, 1405)
	assert.True(client.IsMemberAlreadyExists(a10err))
}

func buildClient(testServer *testing.ServerConfig, sessionId string) (api.Client, error) {
	responseBody := `{"authresponse":{"signature":"` + sessionId + `","description":"the signature should be set in Authorization header for following request."}}`

	testServer.AddRequest().
		Response().
		Body(responseBody, "application/json")

	instance := config.A10Instance{
		APIVersion: 3,
		APIUrl:     testServer.GetURL(),
		UserName:   "usr",
		Password:   "pwd",
	}

	return v3.Connect(&instance)
}
