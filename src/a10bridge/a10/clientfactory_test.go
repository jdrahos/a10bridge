package a10_test

import (
	"a10bridge/a10"
	"a10bridge/config"
	tst "a10bridge/testing"
	"testing"
)

func TestBuildClient_unsupportedApiVersion(t *testing.T) {
	notexistentApiVersion := 55555

	_, err := a10.BuildClient(&config.A10Instance{
		APIVersion: notexistentApiVersion,
		APIUrl:     "localhost:12345",
	})

	if err == nil {
		t.Errorf("Expected error for not supported api version")
	}
}

func TestBuildClient_v2(t *testing.T) {
	expectedUser := "test-user"
	expectedPassword := "test-user"

	testServer := tst.NewTestServer(t).Start()
	defer testServer.Stop()

	testServer.AddRequest().
		Path("/services/rest/V2.1/").
		Method("GET").
		Query("format", "json").
		Query("method", "authenticate").
		Query("username", expectedUser).
		Query("password", expectedPassword).
		Response().
		Body(`{"session_id":"31a9decc4370910de86156fd518888"}`, "application/json")

	instance := config.A10Instance{
		APIVersion: 2,
		APIUrl:     testServer.GetUrl(),
		UserName:   expectedUser,
		Password:   expectedPassword,
	}

	client, err := a10.BuildClient(&instance)

	if err != nil {
		t.Errorf("Failed to build v2 client, %s", err)
	} else if client == nil {
		t.Errorf("Got nil v2 client, %s", err)
	}
}

func TestBuildClient_v3(t *testing.T) {
	testServer := tst.NewTestServer(t).Start()
	defer testServer.Stop()

	testServer.AddRequest().
		Path("/axapi/v3/auth").
		Method("POST").
		Header("Content-Type", "application/json").
		Body(`{
	"credentials": {
		"username": "test-user",
		"password": "test-password"
	}
}`).
		Response().Body(`{
"authresponse" : {
	"signature":"61ed181a0a8a5d06e972b3b4a237c0",
	"description":"the signature should be set in Authorization header for following request."
	}
}`, "application/json")

	instance := config.A10Instance{
		APIVersion: 3,
		APIUrl:     testServer.GetUrl(),
		UserName:   "test-user",
		Password:   "test-password",
	}

	client, err := a10.BuildClient(&instance)

	if err != nil {
		t.Errorf("Failed to build v3 client, %s", err)
	} else if client == nil {
		t.Errorf("Got nil v3 client, %s", err)
	}
}
