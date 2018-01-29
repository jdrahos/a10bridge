package processor_test

import (
	"a10bridge/mocks"
	"a10bridge/model"
	"a10bridge/processor"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type K8sProcessorTestSuite struct {
	suite.Suite
	helper *processor.TestHelper
	client *mocks.K8sClient
}

func (suite *K8sProcessorTestSuite) SetupTest() {
	suite.client = new(mocks.K8sClient)
}
func TestK8sProcessor(t *testing.T) {
	tests := new(K8sProcessorTestSuite)
	tests.helper = new(processor.TestHelper)
	suite.Run(t, tests)
}

func (suite *K8sProcessorTestSuite) TestBuildEnvironment() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	configMap := &model.ConfigMap{
		Name:      "cluster-configs",
		Namespace: "ingress",
		Data: map[string]string{
			"name": "dc-type",
		},
	}

	client.On("GetConfigMap", "ingress", "cluster-configs").Once().Return(configMap, nil)
	env, err := processor.BuildEnvironment()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(env)
	suite.Assert().Equal("dc", env.DataCenter)
	suite.Assert().Equal("type", env.Type)
	suite.Assert().Equal("dc-type", env.Cluster)
}

func (suite *K8sProcessorTestSuite) TestBuildEnvironment_getConfigMapFails() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)

	client.On("GetConfigMap", "ingress", "cluster-configs").Once().Return(nil, errors.New("crap"))
	env, err := processor.BuildEnvironment()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(env)
}

func (suite *K8sProcessorTestSuite) TestBuildEnvironment_configWithoutClusterName() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	configMap := &model.ConfigMap{
		Name:      "cluster-configs",
		Namespace: "ingress",
		Data: map[string]string{
			"not_name": "dc-type",
		},
	}

	client.On("GetConfigMap", "ingress", "cluster-configs").Once().Return(configMap, nil)
	env, err := processor.BuildEnvironment()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(env)
}

func (suite *K8sProcessorTestSuite) TestFindNodes() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	nodeSelectors := map[string]string{
		"ingressNode": "true",
	}
	node1 := model.Node{
		Labels: map[string]string{
			"ingressNode": "true",
		},
		Name:      "node1",
		Weight:    "1",
		IPAddress: "10.10.10.1",
	}
	node2 := model.Node{
		Labels: map[string]string{
			"ingressNode": "false",
		},
		Name:      "node2",
		Weight:    "1",
		IPAddress: "10.10.10.2",
	}
	allNodes := []*model.Node{&node1, &node2}
	client.On("GetNodes").Return(allNodes, nil)

	nodes, err := processor.FindNodes(nodeSelectors)
	suite.Assert().Nil(err)
	suite.Assert().NotNil(nodes)
	suite.Assert().Equal(1, len(nodes))

	passingNode := nodes[0]
	suite.Assert().Equal(node1.Name, passingNode.Name)
}

func (suite *K8sProcessorTestSuite) TestFindNodes_getNodesFails() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	nodeSelectors := map[string]string{
		"ingressNode": "true",
	}
	client.On("GetNodes").Return(nil, errors.New("failed to get nodes"))

	nodes, err := processor.FindNodes(nodeSelectors)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(nodes)
}

func (suite *K8sProcessorTestSuite) TestFindIngressControllers() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	controller1 := model.IngressController{
		Name: "ingress1",
		Health: &model.HealthCheck{
			Name:       "health1",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   10,
			Port:       8080,
			RequiredConsecutivePasses: 5,
			RetryCount:                5,
			Timeout:                   10,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress1-{{.ClusterName}}",
	}
	controller2 := model.IngressController{
		Name: "ingress2",
		Health: &model.HealthCheck{
			Name:       "health2",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   10,
			Port:       8080,
			RequiredConsecutivePasses: 5,
			RetryCount:                5,
			Timeout:                   10,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress2-{{.ClusterName}}",
	}
	expectedControllers := []*model.IngressController{&controller1, &controller2}
	client.On("GetIngressControllers").Return(expectedControllers, nil)

	controllers, err := processor.FindIngressControllers()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(controllers)
	suite.Assert().Equal(expectedControllers, controllers)
}

func (suite *K8sProcessorTestSuite) TestFindIngressControllers_getIngressControllersFails() {
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	client.On("GetIngressControllers").Return(nil, errors.New("fail"))

	controllers, err := processor.FindIngressControllers()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(controllers)
}

func (suite *K8sProcessorTestSuite) TestBuildServiceGroups() {
	expectedServiceGroupName := "service-group-name"
	suite.helper.SetUtilApplyTemplate(func(data interface{}, tpl string) (string, error) {
		return expectedServiceGroupName, nil
	})
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	controller := model.IngressController{
		Name: "ingress1",
		Health: &model.HealthCheck{
			Name:       "health1",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   10,
			Port:       8080,
			RequiredConsecutivePasses: 5,
			RetryCount:                5,
			Timeout:                   10,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress1-{{.ClusterName}}",
	}
	controllers := []*model.IngressController{&controller}
	environment := model.Environment{
		Cluster:    "dc-type",
		DataCenter: "dc",
		Type:       "type",
	}

	serviceGroups := processor.BuildServiceGroups(controllers, &environment)
	suite.Assert().NotNil(serviceGroups)
	suite.Assert().Equal(1, len(serviceGroups))
	actualServiceGroup, found := serviceGroups[expectedServiceGroupName]
	suite.Assert().True(found)
	suite.Assert().Equal(expectedServiceGroupName, actualServiceGroup.Name)
	suite.Assert().NotNil(actualServiceGroup.Health)
	suite.Assert().Equal(expectedServiceGroupName, actualServiceGroup.Health.Name)
	suite.Assert().Equal(controller.Health.Endpoint, actualServiceGroup.Health.Endpoint)
	suite.Assert().Equal(controller.Health.ExpectCode, actualServiceGroup.Health.ExpectCode)
	suite.Assert().Equal(controller.Health.Interval, actualServiceGroup.Health.Interval)
	suite.Assert().Equal(controller.Health.Port, actualServiceGroup.Health.Port)
	suite.Assert().Equal(controller.Health.RequiredConsecutivePasses, actualServiceGroup.Health.RequiredConsecutivePasses)
	suite.Assert().Equal(controller.Health.RetryCount, actualServiceGroup.Health.RetryCount)
	suite.Assert().Equal(controller.Health.Timeout, actualServiceGroup.Health.Timeout)
}

func (suite *K8sProcessorTestSuite) TestBuildServiceGroups_collapseControllersWithTheSameTemplate() {
	expectedServiceGroupName := "service-group-name"
	suite.helper.SetUtilApplyTemplate(func(data interface{}, tpl string) (string, error) {
		return expectedServiceGroupName, nil
	})
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	controller1 := model.IngressController{
		Name: "ingress1",
		Health: &model.HealthCheck{
			Name:       "health1",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   10,
			Port:       8080,
			RequiredConsecutivePasses: 5,
			RetryCount:                5,
			Timeout:                   10,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress1-{{.ClusterName}}",
	}
	controller2 := model.IngressController{
		Name: "ingress2",
		Health: &model.HealthCheck{
			Name:       "health2",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   12,
			Port:       8082,
			RequiredConsecutivePasses: 7,
			RetryCount:                7,
			Timeout:                   12,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress2-{{.ClusterName}}",
	}
	controllers := []*model.IngressController{&controller1, &controller2}
	environment := model.Environment{
		Cluster:    "dc-type",
		DataCenter: "dc",
		Type:       "type",
	}

	serviceGroups := processor.BuildServiceGroups(controllers, &environment)
	suite.Assert().NotNil(serviceGroups)
	suite.Assert().Equal(1, len(serviceGroups))
	actualServiceGroup, found := serviceGroups[expectedServiceGroupName]
	suite.Assert().True(found)
	suite.Assert().Equal(expectedServiceGroupName, actualServiceGroup.Name)
	suite.Assert().NotNil(actualServiceGroup.Health)
	suite.Assert().Equal("/syntheticHealth", actualServiceGroup.Health.Endpoint)
	suite.Assert().Equal("404", actualServiceGroup.Health.ExpectCode)
	suite.Assert().Equal(expectedServiceGroupName, actualServiceGroup.Health.Name)
	suite.Assert().Equal(controller1.Health.Interval, actualServiceGroup.Health.Interval)
	suite.Assert().Equal(controller1.Health.Port, actualServiceGroup.Health.Port)
	suite.Assert().Equal(controller1.Health.RequiredConsecutivePasses, actualServiceGroup.Health.RequiredConsecutivePasses)
	suite.Assert().Equal(controller1.Health.RetryCount, actualServiceGroup.Health.RetryCount)
	suite.Assert().Equal(controller1.Health.Timeout, actualServiceGroup.Health.Timeout)
}

func (suite *K8sProcessorTestSuite) TestBuildServiceGroups_applyTemplateFails() {
	suite.helper.SetUtilApplyTemplate(func(data interface{}, tpl string) (string, error) {
		return "", errors.New("fail")
	})
	client := suite.client
	processor := suite.helper.BuildK8sProcessor(client)
	controller := model.IngressController{
		Name: "ingress1",
		Health: &model.HealthCheck{
			Name:       "health1",
			Endpoint:   "/health",
			ExpectCode: "200",
			Interval:   10,
			Port:       8080,
			RequiredConsecutivePasses: 5,
			RetryCount:                5,
			Timeout:                   10,
		},
		NodeSelectors: map[string]string{
			"ingress": "true",
		},
		Port: 80,
		ServiceGroupNameTemplate: "ingress1-{{.ClusterName}}",
	}
	controllers := []*model.IngressController{&controller}
	environment := model.Environment{
		Cluster:    "dc-type",
		DataCenter: "dc",
		Type:       "type",
	}

	serviceGroups := processor.BuildServiceGroups(controllers, &environment)
	suite.Assert().NotNil(serviceGroups)
	suite.Assert().Equal(0, len(serviceGroups))
}
