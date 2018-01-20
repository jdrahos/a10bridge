package v2_test

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v2"
	"a10bridge/config"
	"a10bridge/testing"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func testConnect(testServer *testing.ServerConfig, assert *assert.Assertions) {
	expectedUser := "test-user"
	expectedPassword := "test-user"
	sessionId := "31a9decc4370910de86156fd518888"
	responseBody := `{"session_id":"` + sessionId + `"}`

	testServer.Reset().
		AddRequest().
		Path("/services/rest/V2.1/").
		Method("GET").
		Query("format", "json").
		Query("method", "authenticate").
		Query("username", expectedUser).
		Query("password", expectedPassword).
		Response().
		Body(responseBody, "application/json")

	instance := config.A10Instance{
		APIVersion: 2,
		APIUrl:     testServer.GetURL(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v2.Connect(&instance)

	assert.Nil(err, "Unexpected error during authentication")
	assert.NotNil(client, "Expected not nil client after authentication")

	actualSessionId := v2.TestHelper{}.GetSessionID(client)
	assert.Equal(sessionId, actualSessionId, "Client using incorrect session id")
}

func testConnect_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions) {
	expectedUser := "test-user"
	expectedPassword := "test-user"

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	instance := config.A10Instance{
		APIVersion: 2,
		APIUrl:     testServer.GetURL(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v2.Connect(&instance)

	assert.NotNil(err, "Expected error when building client and the server responds with 500 during authentication, %s")
	assert.Nil(client, "Expected nil client when authentication fails")
}

func testConnect_failedAuthentication(testServer *testing.ServerConfig, assert *assert.Assertions) {
	expectedUser := "test-user"
	expectedPassword := "test-user"

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": 520486915, "msg": " Admin password error"}}}`, "application/json")

	instance := config.A10Instance{
		APIVersion: 2,
		APIUrl:     testServer.GetURL(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v2.Connect(&instance)

	assert.NotNil(err, "Expected error when building client and the server responds with 500 during authentication, %s")
	assert.Nil(client, "Expected nil client when authentication fails")
}

func testClose(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Method(http.MethodGet).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "session.close").
		Query("session_id", v2.TestHelper{}.GetSessionID(client)).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.Close()

	assert.Nil(err, "Unexpected error when closing client session")
}

func testClose_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.Close()

	assert.NotNil(err, "Expected error when session close call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testClose_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	err := client.Close()
	assert.NotNil(err, "Expected error when session close call fails in a10")
	assert.Equal(errorCode, err.Code())
}
