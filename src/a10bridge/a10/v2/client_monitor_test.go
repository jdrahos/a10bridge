package v2_test

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v2"
	"a10bridge/model"
	"a10bridge/testing"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func testUpdateMonitor(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.update").
		Query("session_id", v2.TestHelper{}.GetSessionID(client)).
		Body(`{
  "health_monitor": {
    "name": "`+monitor.Name+`",
    "retry": `+strconv.Itoa(monitor.RetryCount)+`,
    "consec_pass_reqd": `+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
    "interval": `+strconv.Itoa(monitor.Interval)+`,
    "timeout": `+strconv.Itoa(monitor.Timeout)+`,
    "override_port": `+strconv.Itoa(monitor.Port)+`,
    "type": 3,
    "http": {
      "port": `+strconv.Itoa(monitor.Port)+`,
      "url": "GET `+monitor.Endpoint+`",
      "expect_code": "`+monitor.ExpectCode+`",
      "passive": {
        "status": 0,
        "status_code_2xx": 0,
        "threshold": 75,
        "sample_threshold": 50,
        "interval": 10
      }
    }
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.UpdateHealthMonitor(&monitor)
	assert.Nil(err, "Unexpected error when updating monitor")
}

func testUpdateMonitor_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.UpdateHealthMonitor(&monitor)
	assert.NotNil(err, "Expected error when update monitor call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testUpdateMonitor_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	err := client.UpdateHealthMonitor(&monitor)
	assert.NotNil(err, "Expected error when update monitor call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testCreateMonitor(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.create").
		Query("session_id", v2.TestHelper{}.GetSessionID(client)).
		Body(`{
  "health_monitor": {
    "name": "`+monitor.Name+`",
    "retry": `+strconv.Itoa(monitor.RetryCount)+`,
    "consec_pass_reqd": `+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
    "interval": `+strconv.Itoa(monitor.Interval)+`,
    "timeout": `+strconv.Itoa(monitor.Timeout)+`,
    "override_port": `+strconv.Itoa(monitor.Port)+`,
    "type": 3,
    "http": {
      "port": `+strconv.Itoa(monitor.Port)+`,
      "url": "GET `+monitor.Endpoint+`",
      "expect_code": "`+monitor.ExpectCode+`",
      "passive": {
        "status": 0,
        "status_code_2xx": 0,
        "threshold": 75,
        "sample_threshold": 50,
        "interval": 10
      }
    }
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.CreateHealthMonitor(&monitor)
	assert.Nil(err, "Unexpected error when creating monitor")
}

func testCreateMonitor_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.CreateHealthMonitor(&monitor)
	assert.NotNil(err, "Expected error when create monitor call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testCreateMonitor_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	monitor := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	err := client.CreateHealthMonitor(&monitor)
	assert.NotNil(err, "Expected error when create monitor call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testGetMonitor(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	expected := model.HealthCheck{
		Name:                      "test name",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.search").
		Query("session_id", v2.TestHelper{}.GetSessionID(client)).
		Body(`{
  "name": "`+expected.Name+`"
}`).
		Response().
		Body(`{"health_monitor":{"name":"`+expected.Name+
			`","retry":`+strconv.Itoa(expected.RetryCount)+`,"consec_pass_reqd":`+strconv.Itoa(expected.RequiredConsecutivePasses)+
			`,"interval":`+strconv.Itoa(expected.Interval)+
			`,"timeout":`+strconv.Itoa(expected.Timeout)+
			`,"strictly_retry":0,"disable_after_down":0,"override_ipv4":"0.0.0.0","override_ipv6":"::","override_port":`+strconv.Itoa(expected.Port)+
			`,"type":3,"http":{"port":`+strconv.Itoa(expected.Port)+
			`,"host":"","url":"GET `+expected.Endpoint+
			`","user":"","password":"","expect_code":"`+expected.ExpectCode+
			`","maintenance_code":"","passive":{"status":0,"status_code_2xx":0,"threshold":75,"sample_threshold":50,"interval":10}}}}`,
			"application/json")

	monitor, err := client.GetHealthMonitor(expected.Name)

	assert.Nil(err, "Unexpected error when getting monitor")
	assert.NotNil(monitor, "Expected health check instance")
	assert.Equal(expected.Name, monitor.Name)
	assert.Equal(expected.RetryCount, monitor.RetryCount)
	assert.Equal(expected.RequiredConsecutivePasses, monitor.RequiredConsecutivePasses)
	assert.Equal(expected.Interval, monitor.Interval)
	assert.Equal(expected.Timeout, monitor.Timeout)
	assert.Equal(expected.Port, monitor.Port)
	assert.Equal(expected.Endpoint, monitor.Endpoint)
	assert.Equal(expected.ExpectCode, monitor.ExpectCode)
}

func testGetMonitor_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	_, err := client.GetHealthMonitor("doesn't matter")
	assert.NotNil(err, "Expected error when get monitor call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testGetMonitor_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	_, err := client.GetHealthMonitor("doesn't matter")
	assert.NotNil(err, "Expected error when get monitor call fails in a10")
	assert.Equal(errorCode, err.Code())
}
