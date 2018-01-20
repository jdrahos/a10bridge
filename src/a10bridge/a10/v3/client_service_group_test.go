package v3_test

import (
	"a10bridge/a10/api"
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
		Method(http.MethodPut).
		Path("/axapi/v3/slb/service-group/"+svcGroup.Name).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "service-group": {
    "name": "`+svcGroup.Name+`",
    "protocol": "tcp",
    "health-check": "`+svcGroup.Health.Name+`",
    "member-list": [
      {
        "name" : "`+svcGroup.Members[0].ServerName+`",
        "port" : `+strconv.Itoa(svcGroup.Members[0].Port)+`
      }
    ] 
  }
}`).
		Response().
		Body(`{
  "service-group": {
    "name":"`+svcGroup.Name+`",
    "protocol":"tcp",
    "lb-method":"round-robin",
    "stateless-auto-switch":0,
    "reset-on-server-selection-fail":0,
    "priority-affinity":0,
    "backup-server-event-log":0,
    "strict-select":0,
    "stats-data-action":"stats-data-enable",
    "extended-stats":0,
    "traffic-replication-mirror":0,
    "traffic-replication-mirror-da-repl":0,
    "traffic-replication-mirror-ip-repl":0,
    "traffic-replication-mirror-sa-da-repl":0,
    "traffic-replication-mirror-sa-repl":0,
    "health-check":"`+svcGroup.Health.Name+`",
    "sample-rsp-time":0,
    "uuid":"fafe860c-fb11-11e7-bdaf-97f82d417abc",
    "member-list": [
      {
        "name":"`+svcGroup.Members[0].ServerName+`",
        "port":`+strconv.Itoa(svcGroup.Members[0].Port)+`,
        "member-state":"enable",
        "member-stats-data-disable":0,
        "member-priority":1,
        "uuid":"fb00a914-fb11-11e7-bdaf-97f82d417abc",
        "a10-url":"/axapi/v3/slb/service-group/lga-kube-traefik-test/member/lga-kubnode07+81"
      }
    ]
  }
}
`, "application/json")

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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

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
		Path("/axapi/v3/slb/service-group/").
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Body(`{
  "service-group": {
    "name": "`+svcGroup.Name+`",
    "protocol": "tcp",
    "health-check": "`+svcGroup.Health.Name+`",
    "member-list": [
      {
        "name" : "`+svcGroup.Members[0].ServerName+`",
        "port" : `+strconv.Itoa(svcGroup.Members[0].Port)+`
      }
    ] 
  }
}`).
		Response().
		Body(`{
  "service-group": {
    "name":"`+svcGroup.Name+`",
    "protocol":"tcp",
    "lb-method":"round-robin",
    "stateless-auto-switch":0,
    "reset-on-server-selection-fail":0,
    "priority-affinity":0,
    "backup-server-event-log":0,
    "strict-select":0,
    "stats-data-action":"stats-data-enable",
    "extended-stats":0,
    "traffic-replication-mirror":0,
    "traffic-replication-mirror-da-repl":0,
    "traffic-replication-mirror-ip-repl":0,
    "traffic-replication-mirror-sa-da-repl":0,
    "traffic-replication-mirror-sa-repl":0,
    "health-check":"`+svcGroup.Health.Name+`",
    "sample-rsp-time":0,
    "uuid":"fafe860c-fb11-11e7-bdaf-97f82d417abc",
    "member-list": [
      {
        "name":"`+svcGroup.Members[0].ServerName+`",
        "port":`+strconv.Itoa(svcGroup.Members[0].Port)+`,
        "member-state":"enable",
        "member-stats-data-disable":0,
        "member-priority":1,
        "uuid":"fb00a914-fb11-11e7-bdaf-97f82d417abc",
        "a10-url":"/axapi/v3/slb/service-group/lga-kube-traefik-test/member/lga-kubnode07+81"
      }
    ]
  }
}
`, "application/json")

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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

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
		Method(http.MethodGet).
		Path("/services/rest/V2.1/").
		Path("/axapi/v3/slb/service-group/"+expected.Name).
		Header("Authorization", "A10 "+helper.GetSessionID(client)).
		Response().
		Body(`{
  "service-group": {
    "name":"`+expected.Name+`",
    "protocol":"tcp",
    "lb-method":"round-robin",
    "stateless-auto-switch":0,
    "reset-on-server-selection-fail":0,
    "priority-affinity":0,
    "backup-server-event-log":0,
    "strict-select":0,
    "stats-data-action":"stats-data-enable",
    "extended-stats":0,
    "traffic-replication-mirror":0,
    "traffic-replication-mirror-da-repl":0,
    "traffic-replication-mirror-ip-repl":0,
    "traffic-replication-mirror-sa-da-repl":0,
    "traffic-replication-mirror-sa-repl":0,
    "health-check":"`+expected.Health.Name+`",
    "sample-rsp-time":0,
    "uuid":"fafe860c-fb11-11e7-bdaf-97f82d417abc",
    "member-list": [
      {
        "name":"`+expected.Members[0].ServerName+`",
        "port":`+strconv.Itoa(expected.Members[0].Port)+`,
        "member-state":"enable",
        "member-stats-data-disable":0,
        "member-priority":1,
        "uuid":"fb00a914-fb11-11e7-bdaf-97f82d417abc",
        "a10-url":"/axapi/v3/slb/service-group/lga-kube-traefik-test/member/lga-kubnode07+81"
      }
    ]
  }
}
`, "application/json")

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
		Body(`{"response":{"status":"fail","err":{"code":`+strconv.Itoa(errorCode)+`,"from":"HTTP","msg":"Unauthorized"}}}`, "application/json")

	_, err := client.GetServiceGroup("doesn't matter")
	assert.NotNil(err, "Expected error when get service group call fails in a10")
	assert.Equal(errorCode, err.Code())
}
