package v3_test

import (
	"a10bridge/a10/api"
	"a10bridge/model"
	"a10bridge/testing"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func testCreateMember(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	member := model.Member{
		ServiceGroupName: "sg_name",
		ServerName:       "srv_name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/service-group/"+member.ServiceGroupName+"/member/").
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "member" : {
    "name" : "`+member.ServerName+`",
    "port" : `+strconv.Itoa(member.Port)+`,
    "member-state": "enable",
    "member-stats-data-disable": 0,
    "member-priority": 1
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.CreateMember(&member)
	assert.Nil(err, "Unexpected error when creating member")
}

func testCreateMember_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	member := model.Member{
		ServiceGroupName: "sg name",
		ServerName:       "srv name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.CreateMember(&member)
	assert.NotNil(err, "Expected error when create member call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testCreateMember_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	member := model.Member{
		ServiceGroupName: "sg name",
		ServerName:       "srv name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.CreateMember(&member)
	assert.NotNil(err, "Expected error when create member call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testDeleteMember(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	member := model.Member{
		ServiceGroupName: "sg_name",
		ServerName:       "srv_name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodDelete).
		Path("/axapi/v3/slb/service-group/"+member.ServiceGroupName+"/member/"+member.ServerName+"+"+strconv.Itoa(member.Port)).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.DeleteMember(&member)
	assert.Nil(err, "Unexpected error when creating member")
}

func testDeleteMember_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	member := model.Member{
		ServiceGroupName: "sg name",
		ServerName:       "srv name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.DeleteMember(&member)
	assert.NotNil(err, "Expected error when create member call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testDeleteMember_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	member := model.Member{
		ServiceGroupName: "sg name",
		ServerName:       "srv name",
		Port:             8080,
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.DeleteMember(&member)
	assert.NotNil(err, "Expected error when create member call fails in a10")
	assert.Equal(errorCode, err.Code())
}
