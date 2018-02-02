package main

import (
	"a10bridge/apiserver"
	"a10bridge/model"
	bridgeTesting "a10bridge/testing"
	"a10bridge/util"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes/fake"
)

type MainTestSuite struct {
	suite.Suite
	testServer *bridgeTesting.ServerConfig
	configFile string
	resolver   *bridgeTesting.ConfigurableResolver
}

func (suite *MainTestSuite) SetupTest() {
	suite.resolver.Reset()
	suite.testServer.Reset()
}

func TestA10Bridge(t *testing.T) {
	tests := new(MainTestSuite)

	tests.testServer = bridgeTesting.NewTestServer(t).Start()
	defer tests.testServer.Stop()

	tests.resolver = new(bridgeTesting.ConfigurableResolver)
	originalResolver := util.InjectIPResolver(tests.resolver)
	defer util.InjectIPResolver(originalResolver)

	suite.Run(t, tests)
}

func (suite MainTestSuite) TestV2Protocol() {
	suite.writeConfigFile("2")
	defer os.Remove(suite.configFile)

	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config="+suite.configFile)
	os.Args = append(os.Args, "-interval=10")
	os.Args = append(os.Args, "-sort")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)

	nodeList := newNodeListBuilder(suite).
		addK8sNode("node1", "10.10.10.1", matchingNodeSelector()).
		addK8sNode("node2", "10.10.10.2", notMatchingNodeSelector()).
		addK8sNode("node3", "10.10.10.3", notMatchingNodeSelector()).
		addK8sNode("node4", "10.10.10.4", matchingNodeSelector()).
		buildList()
	configMapList := newConfigMapListBuilder(suite).
		addConfigMap("cluster-configs", "ingress", configMapData()).
		buildList()
	daemonSetList := newDaemonSetListBuilder(suite).
		addDaemonSet("nginx-ingress-controller", "{{.DataCenter}}-nginx-{{.Type}}").
		withContainer("app-http", 90, "/health").
		buildDaemonSet().
		addDaemonSet("traefik-ingress-controller-80", "{{.DataCenter}}-traefik-{{.Type}}").
		withContainer("app-http", 80, "/health").
		buildDaemonSet().
		addDaemonSet("traefik-ingress-controller-81", "{{.DataCenter}}-traefik-{{.Type}}").
		withContainer("app-http", 81, "/health").
		buildDaemonSet().
		buildList()

	clientSet := fake.NewSimpleClientset(&nodeList, &configMapList, &daemonSetList)
	apiserver.InjectFakeClient(clientSet)

	sessionId := "31a9decc4370910de86156fd518888"
	suite.testServer.AddRequest().
		Path("/services/rest/V2.1/").
		Method("GET").
		Query("format", "json").
		Query("method", "authenticate").
		Query("username", "dingo").
		Query("password", "dongo").
		Response().
		Body(`{"session_id":"`+sessionId+`"}`, "application/json")

	existingNode1 := model.Node{
		A10Server: "node1",
		IPAddress: "10.10.1.1",
		Weight:    "5",
	}
	expectedNode1 := model.Node{
		A10Server: "node1",
		IPAddress: "10.10.10.1",
		Weight:    "1",
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.server.search").
		Query("session_id", sessionId).
		Body(v2NameRequest("node1")).
		Response().
		Body(v2ServerResponse(existingNode1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.server.update").
		Query("session_id", sessionId).
		Body(v2ServerRequest(expectedNode1)).
		Response().
		Body(v2OkResponse(), "application/json")

	expectedNode4 := model.Node{
		A10Server: "node4",
		IPAddress: "10.10.10.4",
		Weight:    "1",
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.server.search").
		Query("session_id", sessionId).
		Body(v2NameRequest("node4")).
		Response().
		Body(v2ErrorResponse(67174402), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.server.create").
		Query("session_id", sessionId).
		Body(v2ServerRequest(expectedNode4)).
		Response().
		Body(v2OkResponse(), "application/json")

	existingMonitor1 := model.HealthCheck{
		Name:                      "dc-nginx-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/healtz",
		ExpectCode:                "200",
	}

	expectedMonitor1 := model.HealthCheck{
		Name:                      "dc-nginx-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 3,
		Interval:                  15,
		Timeout:                   10,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.search").
		Query("session_id", sessionId).
		Body(v2NameRequest(expectedMonitor1.Name)).
		Response().
		Body(v2HealthMonitorResponse(existingMonitor1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.update").
		Query("session_id", sessionId).
		Body(v2HealthMonitorRequest(expectedMonitor1)).
		Response().
		Body(v2OkResponse(), "application/json")

	existingSvcGroup1 := model.ServiceGroup{
		Name: "dc-nginx-kube",
		Health: &model.HealthCheck{
			Name: "dc-nginx-test",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node1",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node3",
				Port:             80,
			},
		},
	}
	expectedSvcGroup1 := model.ServiceGroup{
		Name: "dc-nginx-kube",
		Health: &model.HealthCheck{
			Name: "dc-nginx-kube",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node1",
				Port:             90,
			},
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node4",
				Port:             90,
			},
		},
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.search").
		Query("session_id", sessionId).
		Body(v2NameRequest(expectedSvcGroup1.Name)).
		Response().
		Body(v2ServiceGroupResponse(existingSvcGroup1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.update").
		Query("session_id", sessionId).
		Body(v2ServiceGroupRequest(expectedSvcGroup1)).
		Response().
		Body(v2OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.member.create").
		Query("session_id", sessionId).
		Body(v2MemberRequest(expectedSvcGroup1.Members[0])).
		Response().
		Body(v2OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.member.create").
		Query("session_id", sessionId).
		Body(v2MemberRequest(expectedSvcGroup1.Members[1])).
		Response().
		Body(v2OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.member.delete").
		Query("session_id", sessionId).
		Body(v2MemberRequest(existingSvcGroup1.Members[0])).
		Response().
		Body(v2OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.member.delete").
		Query("session_id", sessionId).
		Body(v2MemberRequest(existingSvcGroup1.Members[1])).
		Response().
		Body(v2OkResponse(), "application/json")

	expectedMonitor2 := model.HealthCheck{
		Name:                      "dc-traefik-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 3,
		Interval:                  15,
		Timeout:                   10,
		Port:                      8080,
		Endpoint:                  "/syntheticHealth",
		ExpectCode:                "404",
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.search").
		Query("session_id", sessionId).
		Body(v2NameRequest(expectedMonitor2.Name)).
		Response().
		Body(v2ErrorResponse(33619968), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.hm.create").
		Query("session_id", sessionId).
		Body(v2HealthMonitorRequest(expectedMonitor2)).
		Response().
		Body(v2OkResponse(), "application/json")

	expectedSvcGroup2 := model.ServiceGroup{
		Name: "dc-traefik-kube",
		Health: &model.HealthCheck{
			Name: "dc-traefik-kube",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node1",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node4",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node1",
				Port:             81,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node4",
				Port:             81,
			},
		},
	}

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.search").
		Query("session_id", sessionId).
		Body(v2NameRequest(expectedSvcGroup2.Name)).
		Response().
		Body(v2ErrorResponse(67305473), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "slb.service_group.create").
		Query("session_id", sessionId).
		Body(v2ServiceGroupRequest(expectedSvcGroup2)).
		Response().
		Body(v2OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/services/rest/V2.1/").
		Query("format", "json").
		Query("method", "session.close").
		Query("session_id", sessionId).
		Response().
		Body(v2OkResponse(), "application/json")

	main()
	suite.testServer.AssertNoPendingRequests()
}

func (suite MainTestSuite) TestV3Protocol() {
	suite.writeConfigFile("3")
	defer os.Remove(suite.configFile)

	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config="+suite.configFile)
	os.Args = append(os.Args, "-interval=10")
	os.Args = append(os.Args, "-sort")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)

	nodeList := newNodeListBuilder(suite).
		addK8sNode("node1", "10.10.10.1", matchingNodeSelector()).
		addK8sNode("node2", "10.10.10.2", notMatchingNodeSelector()).
		addK8sNode("node3", "10.10.10.3", notMatchingNodeSelector()).
		addK8sNode("node4", "10.10.10.4", matchingNodeSelector()).
		buildList()
	configMapList := newConfigMapListBuilder(suite).
		addConfigMap("cluster-configs", "ingress", configMapData()).
		buildList()
	daemonSetList := newDaemonSetListBuilder(suite).
		addDaemonSet("nginx-ingress-controller", "{{.DataCenter}}-nginx-{{.Type}}").
		withContainer("app-http", 90, "/health").
		buildDaemonSet().
		addDaemonSet("traefik-ingress-controller-80", "{{.DataCenter}}-traefik-{{.Type}}").
		withContainer("app-http", 80, "/health").
		buildDaemonSet().
		addDaemonSet("traefik-ingress-controller-81", "{{.DataCenter}}-traefik-{{.Type}}").
		withContainer("app-http", 81, "/health").
		buildDaemonSet().
		buildList()

	clientSet := fake.NewSimpleClientset(&nodeList, &configMapList, &daemonSetList)
	apiserver.InjectFakeClient(clientSet)

	sessionId := "31a9decc4370910de86156fd518888"
	suite.testServer.AddRequest().
		Path("/axapi/v3/auth").
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		Body(v3AuthenicationRequest("dingo", "dongo")).
		Response().
		Body(v3AuthenicationResponse(sessionId), "application/json")

	existingNode1 := model.Node{
		A10Server: "node1",
		IPAddress: "10.10.1.1",
		Weight:    "5",
	}
	expectedNode1 := model.Node{
		A10Server: "node1",
		IPAddress: "10.10.10.1",
		Weight:    "1",
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/axapi/v3/slb/server/"+expectedNode1.A10Server).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3ServerResponse(existingNode1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPut).
		Path("/axapi/v3/slb/server/"+expectedNode1.A10Server).
		Header("Authorization", "A10 "+sessionId).
		Body(v3ServerRequest(expectedNode1)).
		Response().
		Body(v3ServerResponse(expectedNode1), "application/json")

	expectedNode4 := model.Node{
		A10Server: "node4",
		IPAddress: "10.10.10.4",
		Weight:    "1",
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/axapi/v3/slb/server/"+expectedNode4.A10Server).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3ErrorResponse(1023460352), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/server/").
		Header("Authorization", "A10 "+sessionId).
		Body(v3ServerRequest(expectedNode4)).
		Response().
		Body(v3ServerResponse(expectedNode4), "application/json")

	existingMonitor1 := model.HealthCheck{
		Name:                      "dc-nginx-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 45,
		Interval:                  1,
		Timeout:                   4654,
		Port:                      8080,
		Endpoint:                  "/healtz",
		ExpectCode:                "200",
	}

	expectedMonitor1 := model.HealthCheck{
		Name:                      "dc-nginx-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 3,
		Interval:                  15,
		Timeout:                   10,
		Port:                      8080,
		Endpoint:                  "/health",
		ExpectCode:                "200",
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/axapi/v3/health/monitor/"+expectedMonitor1.Name).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3HealthMonitorResponse(existingMonitor1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPut).
		Path("/axapi/v3/health/monitor/"+expectedMonitor1.Name).
		Header("Authorization", "A10 "+sessionId).
		Body(v3HealthMonitorRequest(expectedMonitor1)).
		Response().
		Body(v3HealthMonitorResponse(expectedMonitor1), "application/json")

	existingSvcGroup1 := model.ServiceGroup{
		Name: "dc-nginx-kube",
		Health: &model.HealthCheck{
			Name: "dc-nginx-test",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node1",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node3",
				Port:             80,
			},
		},
	}
	expectedSvcGroup1 := model.ServiceGroup{
		Name: "dc-nginx-kube",
		Health: &model.HealthCheck{
			Name: "dc-nginx-kube",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node1",
				Port:             90,
			},
			&model.Member{
				ServiceGroupName: "dc-nginx-kube",
				ServerName:       "node4",
				Port:             90,
			},
		},
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/services/rest/V2.1/").
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3ServiceGroupResponse(existingSvcGroup1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPut).
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name).
		Header("Authorization", "A10 "+sessionId).
		Body(v3ServiceGroupRequest(expectedSvcGroup1)).
		Response().
		Body(v3ServiceGroupResponse(expectedSvcGroup1), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name+"/member/").
		Header("Authorization", "A10 "+sessionId).
		Body(v3MemberRequest(expectedSvcGroup1.Members[0])).
		Response().
		Body(v3OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name+"/member/").
		Header("Authorization", "A10 "+sessionId).
		Body(v3MemberRequest(expectedSvcGroup1.Members[1])).
		Response().
		Body(v3OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodDelete).
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name+"/member/"+existingSvcGroup1.Members[0].ServerName+"+"+strconv.Itoa(existingSvcGroup1.Members[0].Port)).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3OkResponse(), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodDelete).
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup1.Name+"/member/"+existingSvcGroup1.Members[1].ServerName+"+"+strconv.Itoa(existingSvcGroup1.Members[1].Port)).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3OkResponse(), "application/json")

	expectedMonitor2 := model.HealthCheck{
		Name:                      "dc-traefik-kube",
		RetryCount:                5,
		RequiredConsecutivePasses: 3,
		Interval:                  15,
		Timeout:                   10,
		Port:                      8080,
		Endpoint:                  "/syntheticHealth",
		ExpectCode:                "404",
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/axapi/v3/health/monitor/"+expectedMonitor2.Name).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3ErrorResponse(1023460352), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/health/monitor/").
		Header("Authorization", "A10 "+sessionId).
		Body(v3HealthMonitorRequest(expectedMonitor2)).
		Response().
		Body(v3HealthMonitorResponse(expectedMonitor2), "application/json")

	expectedSvcGroup2 := model.ServiceGroup{
		Name: "dc-traefik-kube",
		Health: &model.HealthCheck{
			Name: "dc-traefik-kube",
		},
		Members: []*model.Member{
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node1",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node4",
				Port:             80,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node1",
				Port:             81,
			},
			&model.Member{
				ServiceGroupName: "dc-traefik-kube",
				ServerName:       "node4",
				Port:             81,
			},
		},
	}

	suite.testServer.AddRequest().
		Method(http.MethodGet).
		Path("/services/rest/V2.1/").
		Path("/axapi/v3/slb/service-group/"+expectedSvcGroup2.Name).
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3ErrorResponse(1023460352), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/slb/service-group/").
		Header("Authorization", "A10 "+sessionId).
		Body(v3ServiceGroupRequest(expectedSvcGroup2)).
		Response().
		Body(v3ServiceGroupResponse(expectedSvcGroup2), "application/json")

	suite.testServer.AddRequest().
		Method(http.MethodPost).
		Path("/axapi/v3/logoff").
		Header("Authorization", "A10 "+sessionId).
		Response().
		Body(v3OkResponse(), "application/json")

	main()
	suite.testServer.AssertNoPendingRequests()
}

type nodeListBuilder struct {
	suite *MainTestSuite
	nodes []corev1.Node
}

func newNodeListBuilder(suite MainTestSuite) *nodeListBuilder {
	builder := new(nodeListBuilder)
	builder.suite = &suite
	return builder
}

func (builder *nodeListBuilder) addK8sNode(name, ip string, labels map[string]string) *nodeListBuilder {
	node := corev1.Node{}
	node.SetLabels(labels)
	node.SetName(name)
	builder.nodes = append(builder.nodes, node)
	builder.suite.resolver.AddRecord(name, ip)
	return builder
}

func (builder *nodeListBuilder) buildList() corev1.NodeList {
	return corev1.NodeList{
		Items: builder.nodes,
	}
}

type configMapListBuilder struct {
	suite      *MainTestSuite
	configMaps []corev1.ConfigMap
}

func newConfigMapListBuilder(suite MainTestSuite) *configMapListBuilder {
	builder := new(configMapListBuilder)
	builder.suite = &suite
	return builder
}

func (builder *configMapListBuilder) addConfigMap(name string, namespace string, data map[string]string) *configMapListBuilder {
	configMap := corev1.ConfigMap{}
	configMap.SetName(name)
	configMap.SetNamespace(namespace)
	configMap.Data = data
	builder.configMaps = append(builder.configMaps, configMap)
	return builder
}

func (builder *configMapListBuilder) buildList() corev1.ConfigMapList {
	return corev1.ConfigMapList{
		Items: builder.configMaps,
	}
}

type daemonSetListBuilder struct {
	suite      *MainTestSuite
	daemonSets []extensionsv1beta1.DaemonSet
}

func newDaemonSetListBuilder(suite MainTestSuite) *daemonSetListBuilder {
	builder := new(daemonSetListBuilder)
	builder.suite = &suite
	return builder
}

type daemonSetBuilder struct {
	suite       *MainTestSuite
	listBuilder *daemonSetListBuilder
	daemonSet   extensionsv1beta1.DaemonSet
}

func (builder *daemonSetListBuilder) addDaemonSet(name string, sgTemplate string) *daemonSetBuilder {
	daemonSet := extensionsv1beta1.DaemonSet{}
	daemonSet.SetName(name)
	daemonSet.SetNamespace("ingress")
	daemonSet.SetAnnotations(map[string]string{
		"a10.service_group": sgTemplate,
	})
	daemonSet.Spec.Template.Spec.NodeSelector = matchingNodeSelector()
	return &daemonSetBuilder{
		listBuilder: builder,
		daemonSet:   daemonSet,
	}
}

func (builder *daemonSetBuilder) withContainer(containerPortName string, containerPort int, probePath string) *daemonSetBuilder {
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: probePath,
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}
	builder.daemonSet.Spec.Template.Spec.Containers = append(builder.daemonSet.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          containerPortName,
				HostPort:      int32(containerPort),
				ContainerPort: int32(containerPort + 1000),
			},
		},
		LivenessProbe: &livenessProbe,
	})
	return builder
}

func (builder *daemonSetBuilder) buildDaemonSet() *daemonSetListBuilder {
	builder.listBuilder.daemonSets = append(builder.listBuilder.daemonSets, builder.daemonSet)
	return builder.listBuilder
}

func (builder *daemonSetListBuilder) buildList() extensionsv1beta1.DaemonSetList {
	return extensionsv1beta1.DaemonSetList{
		Items: builder.daemonSets,
	}
}

func matchingNodeSelector() map[string]string {
	return map[string]string{
		"ingress": "true",
	}
}

func notMatchingNodeSelector() map[string]string {
	return map[string]string{
		"ingress": "false",
	}
}

func configMapData() map[string]string {
	return map[string]string{
		"name": "dc-kube-test",
	}
}

func (suite *MainTestSuite) writeConfigFile(version string) {
	binary, err := ioutil.ReadFile("testdata/config.template.v" + version)
	if err != nil {
		suite.T().Errorf("Failed to rad config template, error: %s", err)
	}
	config, err := util.ApplyTemplate(struct{ Url string }{Url: suite.testServer.GetURL()}, string(binary))
	if err != nil {
		suite.T().Errorf("Failed to build config from template, error: %s", err)
	}
	configFile := "/tmp/a10bridge/test.config"
	err = os.MkdirAll("/tmp/a10bridge/", 0700)
	if err != nil {
		suite.T().Errorf("Failed to create folder for the config, error: %s", err)
	}
	err = ioutil.WriteFile(configFile, []byte(config), 0700)
	if err != nil {
		suite.T().Errorf("Failed to write config, error: %s", err)
	}
	suite.configFile = configFile
}

func v2OkResponse() string {
	return `{"response": {"status": "OK"}}`
}

func v2ErrorResponse(errorCode int) string {
	return `{"response": {"status": "fail", "err": {"code": ` + strconv.Itoa(errorCode) + `, "msg": "Error"}}}`
}

func v2NameRequest(name string) string {
	return `{ "name": "` + name + `" }`
}

func v2ServerRequest(node model.Node) string {
	return `{
		"server": {
		  "name": "` + node.A10Server + `",
		  "host": "` + node.IPAddress + `",
		  "weight": ` + node.Weight + `,
		  "conn_limit_log": 1
		}
	  }`
}

func v2ServerResponse(node model.Node) string {
	return `{"server":{"name":"` + node.A10Server + `","host":"` + node.IPAddress + `","gslb_external_address":"0.0.0.0","weight":` + node.Weight + `,"health_monitor":"(default)","status":1,"conn_limit":8000000,"conn_limit_log":1,"conn_resume":0,"stats_data":1,"extended_stats":0,"slow_start":0,"spoofing_cache":0,"template":"default","port_list":[{"port_num":81,"protocol":2,"status":1,"weight":1,"no_ssl":0,"conn_limit":8000000,"conn_limit_log":0,"conn_resume":0,"template":"default","stats_data":1,"health_monitor":"(default)","extended_stats":0},{"port_num":90,"protocol":2,"status":1,"weight":1,"no_ssl":0,"conn_limit":8000000,"conn_limit_log":1,"conn_resume":0,"template":"default","stats_data":1,"health_monitor":"(default)","extended_stats":0}]}}`
}

func v2HealthMonitorRequest(monitor model.HealthCheck) string {
	return `{
		"health_monitor": {
		  "name": "` + monitor.Name + `",
		  "retry": ` + strconv.Itoa(monitor.RetryCount) + `,
		  "consec_pass_reqd": ` + strconv.Itoa(monitor.RequiredConsecutivePasses) + `,
		  "interval": ` + strconv.Itoa(monitor.Interval) + `,
		  "timeout": ` + strconv.Itoa(monitor.Timeout) + `,
		  "override_port": ` + strconv.Itoa(monitor.Port) + `,
		  "type": 3,
		  "http": {
			"port": ` + strconv.Itoa(monitor.Port) + `,
			"url": "GET ` + monitor.Endpoint + `",
			"expect_code": "` + monitor.ExpectCode + `",
			"passive": {
			  "status": 0,
			  "status_code_2xx": 0,
			  "threshold": 75,
			  "sample_threshold": 50,
			  "interval": 10
			}
		  }
		}
	  }`
}

func v2HealthMonitorResponse(monitor model.HealthCheck) string {
	return `{"health_monitor":{"name":"` + monitor.Name +
		`","retry":` + strconv.Itoa(monitor.RetryCount) + `,"consec_pass_reqd":` + strconv.Itoa(monitor.RequiredConsecutivePasses) +
		`,"interval":` + strconv.Itoa(monitor.Interval) +
		`,"timeout":` + strconv.Itoa(monitor.Timeout) +
		`,"strictly_retry":0,"disable_after_down":0,"override_ipv4":"0.0.0.0","override_ipv6":"::","override_port":` + strconv.Itoa(monitor.Port) +
		`,"type":3,"http":{"port":` + strconv.Itoa(monitor.Port) +
		`,"host":"","url":"GET ` + monitor.Endpoint +
		`","user":"","password":"","expect_code":"` + monitor.ExpectCode +
		`","maintenance_code":"","passive":{"status":0,"status_code_2xx":0,"threshold":75,"sample_threshold":50,"interval":10}}}}`
}

func v2ServiceGroupRequest(serviceGroup model.ServiceGroup) string {
	requestBody := `{
		"service_group": {
		  "name": "` + serviceGroup.Name + `",
		  "protocol": 2,
		  "health_monitor": "` + serviceGroup.Health.Name + `",
		  "member_list": [`

	for idx, member := range serviceGroup.Members {
		if idx != 0 {
			requestBody += ","
		}
		requestBody += ` {
				"server" : "` + member.ServerName + `",
				"port" : ` + strconv.Itoa(member.Port) + `
			}`
	}

	return requestBody + `		  ] 
		}
	  }`
}

func v2ServiceGroupResponse(serviceGroup model.ServiceGroup) string {
	responseBody := `{"service_group":{"name":"` + serviceGroup.Name +
		`","protocol":2,"lb_method":0,"health_monitor":"` + serviceGroup.Health.Name +
		`","policy_template":"","port_template":"","server_template":"","priority_affinity":0,"sample_rsp_time":0,` +
		`"sample_rsp_time_rpt_ext_ser_top_fastest":0,"sample_rsp_time_rpt_ext_ser_top_slowest":0,"sample_rsp_time_rpt_ext_ser_report_delay":0,` +
		`"traffic_repl_mirr_da_repl":0,"traffic_repl_mirr_sa_repl":0,"traffic_repl_mirr_sa_da_repl":0,"traffic_repl_mirr_ip_repl":0,` +
		`"traffic_repl_mirr":0,"action_list":[{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},` +
		`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},` +
		`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},` +
		`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},` +
		`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0},` +
		`{"action":0,"exceed_limit_only":0},{"action":0,"exceed_limit_only":0}],"min_active_member":{"status":0,"number":0,"priority_set":0},` +
		`"backup_server_event_log_enable":0,"client_reset":0,"stats_data":1,"extended_stats":0,"member_list":[`

	for idx, member := range serviceGroup.Members {
		if idx != 0 {
			responseBody += ","
		}
		responseBody += ` {"server":"` + member.ServerName +
			`","port":` + strconv.Itoa(member.Port) +
			`,"template":"default","priority":1,"status":1,"stats_data":1}`
	}

	return responseBody + `]}}`
}

func v2MemberRequest(member *model.Member) string {
	return `{
		"member" : {
		  "server" : "` + member.ServerName + `",
		  "port" : ` + strconv.Itoa(member.Port) + `
		},
		"name" : "` + member.ServiceGroupName + `"
	  }`
}

func v3OkResponse() string {
	return `{"response": {"status": "OK"}}`
}

func v3ErrorResponse(errorCode int) string {
	return `{"response":{"status":"fail","err":{"code":` + strconv.Itoa(errorCode) + `,"from":"HTTP","msg":"Failure"}}}`
}

func v3AuthenicationRequest(username, password string) string {
	return `{
		"credentials": {
			"username": "` + username + `",
			"password": "` + password + `"
		}
	}`
}

func v3AuthenicationResponse(sessionId string) string {
	return `{"authresponse":{"signature":"` + sessionId + `","description":"the signature should be set in Authorization header for following request."}}`
}

func v3ServerRequest(node model.Node) string {
	return `{
		"server": {
		  "name": "` + node.A10Server + `",
		  "host": "` + node.IPAddress + `",
		  "action": "enable",
		  "weight": ` + node.Weight + `,
		  "conn-limit": 1000
		}
	  }`
}

func v3ServerResponse(node model.Node) string {
	return `{
		"server": {
		  "name": "` + node.A10Server + `",
		  "host": "` + node.IPAddress + `",
		  "action": "enable",
		  "weight": ` + node.Weight + `,
		  "conn-limit": 1000
		}
	  }`
}

func v3HealthMonitorRequest(monitor model.HealthCheck) string {
	return `{
		"monitor": {
		  "name": "` + monitor.Name + `",
		  "retry": ` + strconv.Itoa(monitor.RetryCount) + `,
		  "up-retry": ` + strconv.Itoa(monitor.RequiredConsecutivePasses) + `,
		  "interval": ` + strconv.Itoa(monitor.Interval) + `,
		  "timeout": ` + strconv.Itoa(monitor.Timeout) + `,
		  "override-port": ` + strconv.Itoa(monitor.Port) + `,
		  "passive":0,
		  "strict-retry-on-server-err-resp":1,
		  "disable-after-down":0,
		  "method":{
			"http": {
			  "http":1,
			  "http-port": ` + strconv.Itoa(monitor.Port) + `,
			  "http-url":1,
			  "http-expect":1,
			  "http-response-code": "` + monitor.ExpectCode + `",
			  "url-type":"GET",
			  "url-path":"` + monitor.Endpoint + `",
			  "http-kerberos-auth":0
			}
		  }
		}
	  }`
}

func v3HealthMonitorResponse(monitor model.HealthCheck) string {
	return `{
		"monitor": {
			"name":"` + monitor.Name + `",
			"dsr-l2-strict":0,
			"retry":` + strconv.Itoa(monitor.RetryCount) + `,
			"up-retry":` + strconv.Itoa(monitor.RequiredConsecutivePasses) + `,
			"override-port":` + strconv.Itoa(monitor.Port) + `,
			"passive":0,
			"strict-retry-on-server-err-resp":1,
			"disable-after-down":0,
			"interval":` + strconv.Itoa(monitor.Interval) + `,
			"timeout":` + strconv.Itoa(monitor.Timeout) + `,
			"ssl-ciphers":"DEFAULT",
			"uuid":"6f577314-fb09-11e7-bdaf-97f82d417abc",
			"method": {
			"http": {
				"http":1,
				"http-port":` + strconv.Itoa(monitor.Port) + `,
				"http-expect":1,
				"http-response-code":"` + monitor.ExpectCode + `",
				"http-url":1,
				"url-type":"GET",
				"url-path":"` + monitor.Endpoint + `",
				"http-kerberos-auth":0,
				"uuid":"6f57cd5a-fb09-11e7-bdaf-97f82d417abc",
				"a10-url":"/axapi/v3/health/monitor/` + monitor.Name + `/method/http"
			},
			"a10-url":"/axapi/v3/health/monitor/` + monitor.Name + `/method"
			}
		}
	}`
}

func v3ServiceGroupRequest(serviceGroup model.ServiceGroup) string {
	requestBody := `{
		"service-group": {
		  "name": "` + serviceGroup.Name + `",
		  "protocol": "tcp",
		  "health-check": "` + serviceGroup.Health.Name + `",
		  "member-list": [`

	for idx, member := range serviceGroup.Members {
		if idx != 0 {
			requestBody += ", "
		}
		requestBody += ` {
			  "name" : "` + member.ServerName + `",
			  "port" : ` + strconv.Itoa(member.Port) + `
			}`
	}

	return requestBody + ` ] } }`
}

func v3ServiceGroupResponse(serviceGroup model.ServiceGroup) string {
	responseBody := `{
		"service-group": {
		  "name":"` + serviceGroup.Name + `",
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
		  "health-check":"` + serviceGroup.Health.Name + `",
		  "sample-rsp-time":0,
		  "uuid":"fafe860c-fb11-11e7-bdaf-97f82d417abc",
		  "member-list": [`

	for idx, member := range serviceGroup.Members {
		if idx != 0 {
			responseBody += ","
		}
		responseBody += `
			{
			  "name":"` + member.ServerName + `",
			  "port":` + strconv.Itoa(member.Port) + `,
			  "member-state":"enable",
			  "member-stats-data-disable":0,
			  "member-priority":1,
			  "uuid":"fb00a914-fb11-11e7-bdaf-97f82d417abc",
			  "a10-url":"/axapi/v3/slb/service-group/lga-kube-traefik-test/member/lga-kubnode07+81"
			}
		  `
	}
	return responseBody + `]
		}
	  }`
}

func v3MemberRequest(member *model.Member) string {
	return `{
		"member" : {
		  "name" : "` + member.ServerName + `",
		  "port" : ` + strconv.Itoa(member.Port) + `,
		  "member-state": "enable",
		  "member-stats-data-disable": 0,
		  "member-priority": 1
		}
	  }`
}
