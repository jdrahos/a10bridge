package v3_test

import (
	"a10bridge/a10/api"
	"a10bridge/model"
	"a10bridge/testing"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func testUpdateMonitor(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	monitor := model.HealthCheck{
		Name:                      "test_name",
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
		Method(http.MethodPut).
		Path("/axapi/v3/health/monitor/"+monitor.Name).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "monitor": {
    "name": "`+monitor.Name+`",
    "retry": `+strconv.Itoa(monitor.RetryCount)+`,
    "up-retry": `+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
    "interval": `+strconv.Itoa(monitor.Interval)+`,
    "timeout": `+strconv.Itoa(monitor.Timeout)+`,
    "override-port": `+strconv.Itoa(monitor.Port)+`,
    "passive":0,
    "strict-retry-on-server-err-resp":1,
    "disable-after-down":0,
    "method":{
      "http": {
        "http":1,
        "http-port": `+strconv.Itoa(monitor.Port)+`,
        "http-url":1,
        "http-expect":1,
        "http-response-code": "`+monitor.ExpectCode+`",
        "url-type":"GET",
        "url-path":"`+monitor.Endpoint+`",
        "http-kerberos-auth":0
      }
    }
  }
}`).
		Response().
		Body(`{
	"monitor": {
		"name":"`+monitor.Name+`",
		"dsr-l2-strict":0,
		"retry":`+strconv.Itoa(monitor.RetryCount)+`,
		"up-retry":`+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
		"override-port":`+strconv.Itoa(monitor.Port)+`,
		"passive":0,
		"strict-retry-on-server-err-resp":1,
		"disable-after-down":0,
		"interval":`+strconv.Itoa(monitor.Interval)+`,
		"timeout":`+strconv.Itoa(monitor.Timeout)+`,
		"ssl-ciphers":"DEFAULT",
		"uuid":"6f577314-fb09-11e7-bdaf-97f82d417abc",
		"method": {
		"http": {
			"http":1,
			"http-port":`+strconv.Itoa(monitor.Port)+`,
			"http-expect":1,
			"http-response-code":"`+monitor.ExpectCode+`",
			"http-url":1,
			"url-type":"GET",
			"url-path":"`+monitor.Endpoint+`",
			"http-kerberos-auth":0,
			"uuid":"6f57cd5a-fb09-11e7-bdaf-97f82d417abc",
			"a10-url":"/axapi/v3/health/monitor/`+monitor.Name+`/method/http"
		},
		"a10-url":"/axapi/v3/health/monitor/`+monitor.Name+`/method"
		}
	}
}`, "application/json")

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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

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
		Path("/axapi/v3/health/monitor/").
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "monitor": {
    "name": "`+monitor.Name+`",
    "retry": `+strconv.Itoa(monitor.RetryCount)+`,
    "up-retry": `+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
    "interval": `+strconv.Itoa(monitor.Interval)+`,
    "timeout": `+strconv.Itoa(monitor.Timeout)+`,
    "override-port": `+strconv.Itoa(monitor.Port)+`,
    "passive":0,
    "strict-retry-on-server-err-resp":1,
    "disable-after-down":0,
    "method":{
      "http": {
        "http":1,
        "http-port": `+strconv.Itoa(monitor.Port)+`,
        "http-url":1,
        "http-expect":1,
        "http-response-code": "`+monitor.ExpectCode+`",
        "url-type":"GET",
        "url-path":"`+monitor.Endpoint+`",
        "http-kerberos-auth":0
      }
    }
  }
}`).
		Response().
		Body(`{
	"monitor": {
		"name":"`+monitor.Name+`",
		"dsr-l2-strict":0,
		"retry":`+strconv.Itoa(monitor.RetryCount)+`,
		"up-retry":`+strconv.Itoa(monitor.RequiredConsecutivePasses)+`,
		"override-port":`+strconv.Itoa(monitor.Port)+`,
		"passive":0,
		"strict-retry-on-server-err-resp":1,
		"disable-after-down":0,
		"interval":`+strconv.Itoa(monitor.Interval)+`,
		"timeout":`+strconv.Itoa(monitor.Timeout)+`,
		"ssl-ciphers":"DEFAULT",
		"uuid":"6f577314-fb09-11e7-bdaf-97f82d417abc",
		"method": {
		"http": {
			"http":1,
			"http-port":`+strconv.Itoa(monitor.Port)+`,
			"http-expect":1,
			"http-response-code":"`+monitor.ExpectCode+`",
			"http-url":1,
			"url-type":"GET",
			"url-path":"`+monitor.Endpoint+`",
			"http-kerberos-auth":0,
			"uuid":"6f57cd5a-fb09-11e7-bdaf-97f82d417abc",
			"a10-url":"/axapi/v3/health/monitor/`+monitor.Name+`/method/http"
		},
		"a10-url":"/axapi/v3/health/monitor/`+monitor.Name+`/method"
		}
	}
}`, "application/json")

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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	err := client.CreateHealthMonitor(&monitor)
	assert.NotNil(err, "Expected error when create monitor call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testGetMonitor(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	expected := model.HealthCheck{
		Name:                      "test_name",
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
		Method(http.MethodGet).
		Path("/axapi/v3/health/monitor/"+expected.Name).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Response().
		Body(`{
	"monitor": {
		"name":"`+expected.Name+`",
		"dsr-l2-strict":0,
		"retry":`+strconv.Itoa(expected.RetryCount)+`,
		"up-retry":`+strconv.Itoa(expected.RequiredConsecutivePasses)+`,
		"override-port":`+strconv.Itoa(expected.Port)+`,
		"passive":0,
		"strict-retry-on-server-err-resp":1,
		"disable-after-down":0,
		"interval":`+strconv.Itoa(expected.Interval)+`,
		"timeout":`+strconv.Itoa(expected.Timeout)+`,
		"ssl-ciphers":"DEFAULT",
		"uuid":"6f577314-fb09-11e7-bdaf-97f82d417abc",
		"method": {
		"http": {
			"http":1,
			"http-port":`+strconv.Itoa(expected.Port)+`,
			"http-expect":1,
			"http-response-code":"`+expected.ExpectCode+`",
			"http-url":1,
			"url-type":"GET",
			"url-path":"`+expected.Endpoint+`",
			"http-kerberos-auth":0,
			"uuid":"6f57cd5a-fb09-11e7-bdaf-97f82d417abc",
			"a10-url":"/axapi/v3/health/monitor/`+expected.Name+`/method/http"
		},
		"a10-url":"/axapi/v3/health/monitor/`+expected.Name+`/method"
		}
	}
}`,
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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	_, err := client.GetHealthMonitor("doesn't matter")
	assert.NotNil(err, "Expected error when get monitor call fails in a10")
	assert.Equal(errorCode, err.Code())
}
