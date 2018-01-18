package v2_test

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v2"
	"a10bridge/config"
	"a10bridge/testing"
	tst "testing"

	"github.com/stretchr/testify/assert"
)

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

func buildClient(testServer *testing.ServerConfig, sessionId string) (api.Client, error) {
	testServer.AddRequest().
		Response().
		Body(`{"session_id":"`+sessionId+`"}`, "application/json")

	instance := config.A10Instance{
		APIVersion: 2,
		APIUrl:     testServer.GetUrl(),
		UserName:   "usr",
		Password:   "pwd",
	}

	return v2.Connect(&instance)
}
