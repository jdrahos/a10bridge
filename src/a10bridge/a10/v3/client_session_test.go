package v3_test

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v3"
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
	responseBody := `{"authresponse":{"signature":"` + sessionId + `","description":"the signature should be set in Authorization header for following request."}}`

	testServer.Reset().
		AddRequest().
		Path("/axapi/v3/auth").
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		Body(`{
    "credentials": {
        "username": "`+expectedUser+`",
        "password": "`+expectedPassword+`"
    }
}`).
		Response().
		Body(responseBody, "application/json")

	instance := config.A10Instance{
		APIVersion: 3,
		APIUrl:     testServer.GetUrl(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v3.Connect(&instance)

	assert.Nil(err, "Unexpected error during authentication")
	assert.NotNil(client, "Expected not nil client after authentication")

	actualSessionId := helper.GetSessionId(client)
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
		APIVersion: 3,
		APIUrl:     testServer.GetUrl(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v3.Connect(&instance)

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
		APIVersion: 3,
		APIUrl:     testServer.GetUrl(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := v3.Connect(&instance)

	assert.NotNil(err, "Expected error when building client and the server responds with 500 during authentication, %s")
	assert.Nil(client, "Expected nil client when authentication fails")
}

func testClose(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/logoff").
		Header("Authorization", "A10 "+helper.GetSessionId(client)).
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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.Close()
	assert.NotNil(err, "Expected error when session close call fails in a10")
	assert.Equal(errorCode, err.Code())
}
