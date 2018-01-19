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

func testUpdateServiceGroup(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.update").
		Query("session_id", v2.TestHelper{}.GetSessionId(client)).
		Body(`{
  "service_group": {
    "name": "`+svcGroup.Name+`",
    "protocol": 2,
    "health_monitor": "`+svcGroup.Health.Name+`",
    "member_list": [
	  {
        "server" : "`+svcGroup.Members[0].ServerName+`",
        "port" : `+strconv.Itoa(svcGroup.Members[0].Port)+`
      }
    ] 
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.UpdateServiceGroup(&svcGroup)
	assert.Nil(err, "Unexpected error when updating service group")
}

func testUpdateServiceGroup_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.UpdateServiceGroup(&svcGroup)
	assert.NotNil(err, "Expected error when update service group call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testUpdateServiceGroup_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	err := client.UpdateServiceGroup(&svcGroup)
	assert.NotNil(err, "Expected error when update service group call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testCreateServiceGroup(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.create").
		Query("session_id", v2.TestHelper{}.GetSessionId(client)).
		Body(`{
  "service_group": {
    "name": "`+svcGroup.Name+`",
    "protocol": 2,
    "health_monitor": "`+svcGroup.Health.Name+`",
    "member_list": [
	  {
        "server" : "`+svcGroup.Members[0].ServerName+`",
        "port" : `+strconv.Itoa(svcGroup.Members[0].Port)+`
      }
    ] 
  }
}`).
		Response().
		Body(`{"response": {"status": "OK"}}`, "application/json")

	err := client.CreateServiceGroup(&svcGroup)
	assert.Nil(err, "Unexpected error when creating service group")
}

func testCreateServiceGroup_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := client.CreateServiceGroup(&svcGroup)
	assert.NotNil(err, "Expected error when create service group call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testCreateServiceGroup_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	svcGroup := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	err := client.CreateServiceGroup(&svcGroup)
	assert.NotNil(err, "Expected error when create service group call fails in a10")
	assert.Equal(errorCode, err.Code())
}

func testGetServiceGroup(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	expected := model.ServiceGroup{
		Name: "test name",
		Health: &model.HealthCheck{
			Name: "monitor name",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "test name",
				ServerName:       "server name",
				Port:             8080,
			},
		},
	}

	testServer.Reset().
		AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.search").
		Query("session_id", v2.TestHelper{}.GetSessionId(client)).
		Body(`{
  "name": "`+expected.Name+`"
}`).
		Response().
		Body(`{"service_group":{"name":"`+expected.Name+
			`","protocol":2,"lb_method":0,"health_monitor":"`+expected.Health.Name+
			`","policy_template":"","port_template":"","server_template":"","priority_affinity":0,"sample_rsp_time":0,`+
			`"sample_rsp_time_rpt_ext_ser_top_fastest":0,"sample_rsp_time_rpt_ext_ser_top_slowest":0,"sample_rsp_time_rpt_ext_ser_report_delay":0,`+
			`"traffic_repl_mirr_da_repl":0,"traffic_repl_mirr_sa_repl":0,"traffic_repl_mirr_sa_da_repl":0,"traffic_repl_mirr_ip_repl":0,`+
			`"traffic_repl_mirr":0,"action_list":[{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},`+
			`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},`+
			`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},`+
			`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},`+
			`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},`+
			`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0}],"min_active_member":{"status":0,"number":0,"priority_set":0},`+
			`"backup_server_event_log_enable":0,"client_reset":0,"stats_data":1,"extended_stats":0,"member_list":[`+
			`{"server":"`+expected.Members[0].ServerName+
			`","port":`+strconv.Itoa(expected.Members[0].Port)+
			`,"template":"default","priority":1,"status":1,"stats_data":1}]}}`,
			"application/json")

	svcGroup, err := client.GetServiceGroup(expected.Name)

	assert.Nil(err, "Unexpected error when getting service group")
	assert.NotNil(svcGroup, "Expected service group instance")
	assert.Equal(expected.Name, svcGroup.Name)
	assert.Equal(expected.Health.Name, svcGroup.Health.Name)
	assert.Equal(len(expected.Members), len(svcGroup.Members))
	assert.Equal(1, len(svcGroup.Members))
	assert.Equal(expected.Members[0].ServerName, svcGroup.Members[0].ServerName)
	assert.Equal(expected.Members[0].Port, svcGroup.Members[0].Port)
	assert.Equal(expected.Members[0].ServiceGroupName, svcGroup.Members[0].ServiceGroupName)
}

func testGetServiceGroup_ServerError(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	_, err := client.GetServiceGroup("doesn't matter")
	assert.NotNil(err, "Expected error when get service group call fails because of server issues")
	assert.Equal(0, err.Code(), "Expected 0 failure code for errors not returned by a10")
}

func testGetServiceGroup_Failure(testServer *testing.ServerConfig, assert *assert.Assertions, client api.Client) {
	errorCode := 1009
	testServer.Reset().
		AddRequest().
		Response().
		Body(`{"response": {"status": "fail", "err": {"code": `+strconv.Itoa(errorCode)+`, "msg": "Invalid session ID"}}}`, "application/json")

	_, err := client.GetServiceGroup("doesn't matter")
	assert.NotNil(err, "Expected error when get service group call fails in a10")
	assert.Equal(errorCode, err.Code())
}
