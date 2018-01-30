package apiserver_test

import (
	"a10bridge/apiserver"
	"errors"
	"net"
	"strconv"
	"testing"

	k8stesting "k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

type ClientTestSuite struct {
	suite.Suite
	helper *apiserver.TestHelper
}

func (suite *ClientTestSuite) SetupTest() {
}

func TestClient(t *testing.T) {
	tests := new(ClientTestSuite)
	tests.helper = new(apiserver.TestHelper)
	suite.Run(t, tests)
}

func (suite *ClientTestSuite) TestGetNodes() {
	ip1 := []byte{10, 10, 10, 1}
	original := suite.helper.SetNetLookupIPFunc(func(host string) ([]net.IP, error) {
		return []net.IP{ip1}, nil
	})
	defer suite.helper.SetNetLookupIPFunc(original)

	expectedLabels := map[string]string{
		"test":    "label",
		"another": "test label",
	}
	expectedName := "node1"
	defaultA10Server := expectedName
	defaultWeight := "1"
	node1 := corev1.Node{}
	node1.SetLabels(expectedLabels)
	node1.SetName(expectedName)
	nodeList := corev1.NodeList{
		Items: []corev1.Node{node1},
	}

	clientset := fake.NewSimpleClientset(&nodeList)
	client := suite.helper.BuildClient(clientset)

	nodes, err := client.GetNodes()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(nodes)
	suite.Assert().Equal(1, len(nodes))
	suite.Assert().Equal(expectedLabels, nodes[0].Labels)
	suite.Assert().Equal(expectedName, nodes[0].Name)
	suite.Assert().Equal("10.10.10.1", nodes[0].IPAddress)
	suite.Assert().Equal(defaultA10Server, nodes[0].A10Server)
	suite.Assert().Equal(defaultWeight, nodes[0].Weight)
}

func (suite *ClientTestSuite) TestGetNodes_serverNameFromAnnotation() {
	ip1 := []byte{10, 10, 10, 1}
	original := suite.helper.SetNetLookupIPFunc(func(host string) ([]net.IP, error) {
		return []net.IP{ip1}, nil
	})
	defer suite.helper.SetNetLookupIPFunc(original)

	expectedA10Server := "server1"
	annotations := map[string]string{
		"a10.server": expectedA10Server,
	}
	expectedLabels := map[string]string{
		"test":    "label",
		"another": "test label",
	}
	expectedName := "node1"
	node1 := corev1.Node{}
	node1.SetAnnotations(annotations)
	node1.SetLabels(expectedLabels)
	node1.SetName(expectedName)
	nodeList := corev1.NodeList{
		Items: []corev1.Node{node1},
	}

	clientset := fake.NewSimpleClientset(&nodeList)
	client := suite.helper.BuildClient(clientset)

	nodes, err := client.GetNodes()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(nodes)
	suite.Assert().Equal(expectedA10Server, nodes[0].A10Server)
}

func (suite *ClientTestSuite) TestGetNodes_serverWeightFromAnnotation() {
	ip1 := []byte{10, 10, 10, 1}
	original := suite.helper.SetNetLookupIPFunc(func(host string) ([]net.IP, error) {
		return []net.IP{ip1}, nil
	})
	defer suite.helper.SetNetLookupIPFunc(original)

	expectedWeight := "152"
	annotations := map[string]string{
		"a10.server.weight": expectedWeight,
	}
	expectedLabels := map[string]string{
		"test":    "label",
		"another": "test label",
	}
	expectedName := "node1"
	node1 := corev1.Node{}
	node1.SetAnnotations(annotations)
	node1.SetLabels(expectedLabels)
	node1.SetName(expectedName)
	nodeList := corev1.NodeList{
		Items: []corev1.Node{node1},
	}

	clientset := fake.NewSimpleClientset(&nodeList)
	client := suite.helper.BuildClient(clientset)

	nodes, err := client.GetNodes()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(nodes)
	suite.Assert().Equal(expectedWeight, nodes[0].Weight)
}

func (suite *ClientTestSuite) TestGetNodes_ipResolutionFails() {
	original := suite.helper.SetNetLookupIPFunc(func(host string) ([]net.IP, error) {
		return nil, errors.New("fail")
	})
	defer suite.helper.SetNetLookupIPFunc(original)

	nodeList := corev1.NodeList{Items: []corev1.Node{corev1.Node{}}}
	clientset := fake.NewSimpleClientset(&nodeList)
	client := suite.helper.BuildClient(clientset)

	nodes, err := client.GetNodes()

	suite.Assert().NotNil(err)
	suite.Assert().Nil(nodes)
}

func (suite *ClientTestSuite) TestGetNodes_apiCallFails() {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("*", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("fail")
	})
	client := suite.helper.BuildClient(clientset)
	nodes, err := client.GetNodes()

	suite.Assert().NotNil(err)
	suite.Assert().Nil(nodes)
}

func (suite *ClientTestSuite) TestGetConfigMap() {
	namespace := "testns"
	name := "configMap"
	configMap1 := corev1.ConfigMap{}
	configMap1.SetName(name)
	configMap1.SetNamespace(namespace)
	configMap1.Data = map[string]string{
		"data": "value",
	}
	configMap2 := corev1.ConfigMap{}
	configMap2.SetName(name + "_not_it")
	configMap2.SetNamespace(namespace + "_not_it")
	configMap2.Data = map[string]string{
		"data": "value _not_it",
	}
	configMapList := corev1.ConfigMapList{
		Items: []corev1.ConfigMap{configMap1, configMap2},
	}

	clientset := fake.NewSimpleClientset(&configMapList)
	client := suite.helper.BuildClient(clientset)
	configMap, err := client.GetConfigMap(namespace, name)
	suite.Assert().Nil(err)
	suite.Assert().NotNil(configMap)
	suite.Assert().Equal(name, configMap.Name)
	suite.Assert().Equal(namespace, configMap.Namespace)
	suite.Assert().Equal(configMap1.Data, configMap.Data)
}

func (suite *ClientTestSuite) TestGetConfigMap_apiCallFails() {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("*", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("fail")
	})
	client := suite.helper.BuildClient(clientset)
	configMap, err := client.GetConfigMap("", "")

	suite.Assert().NotNil(err)
	suite.Assert().Nil(configMap)
}

func (suite *ClientTestSuite) TestGetIngressControllers() {
	expectedName := "test-ingress-controller-80"
	expectedNodeSelector := map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	expectedPort := 8080
	expectedServiceGroupTemplate := "svc grp name template"
	defaultHttpStatusCode := "200"
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName(expectedName)
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group": expectedServiceGroupTemplate,
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = expectedNodeSelector
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{})
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_jmx",
				HostPort:      8181,
				ContainerPort: 9191,
			},
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      int32(expectedPort),
				ContainerPort: 9090,
			},
		},
		LivenessProbe: &livenessProbe,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(controllers)
	suite.Assert().Equal(1, len(controllers))
	suite.Assert().Equal(expectedName, controllers[0].Name)
	suite.Assert().Equal(expectedNodeSelector, controllers[0].NodeSelectors)
	suite.Assert().Equal(expectedPort, controllers[0].Port)
	suite.Assert().Equal(expectedServiceGroupTemplate, controllers[0].ServiceGroupNameTemplate)
	suite.Assert().NotNil(controllers[0].Health)
	suite.Assert().Equal(livenessProbe.HTTPGet.Path, controllers[0].Health.Endpoint)
	suite.Assert().Equal(int(livenessProbe.HTTPGet.Port.IntVal), controllers[0].Health.Port)
	suite.Assert().Equal(int(livenessProbe.PeriodSeconds), controllers[0].Health.Interval)
	suite.Assert().Equal(int(livenessProbe.FailureThreshold), controllers[0].Health.RetryCount)
	suite.Assert().Equal(int(livenessProbe.SuccessThreshold), controllers[0].Health.RequiredConsecutivePasses)
	suite.Assert().Equal(int(livenessProbe.TimeoutSeconds), controllers[0].Health.Timeout)
	suite.Assert().Equal(defaultHttpStatusCode, controllers[0].Health.ExpectCode)
}

func (suite *ClientTestSuite) TestGetIngressControllers_healthCheckFromAnnotations() {
	expectedHealthCheckPath := "/healtz"
	expectedHealthCheckPort := 8088
	expectedName := "test-ingress-controller"
	expectedNodeSelector := map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	expectedPort := 8080
	expectedServiceGroupTemplate := "svc grp name template"
	defaultHttpStatusCode := "200"
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName(expectedName)
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group":   expectedServiceGroupTemplate,
		"a10.health.endpoint": expectedHealthCheckPath,
		"a10.health.port":     strconv.Itoa(expectedHealthCheckPort),
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = expectedNodeSelector
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_jmx",
				HostPort:      8181,
				ContainerPort: 9191,
			},
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      int32(expectedPort),
				ContainerPort: 9090,
			},
		},
		LivenessProbe: &livenessProbe,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(controllers)

	suite.Assert().NotNil(controllers[0].Health)
	suite.Assert().Equal(expectedHealthCheckPath, controllers[0].Health.Endpoint)
	suite.Assert().Equal(expectedHealthCheckPort, controllers[0].Health.Port)
	suite.Assert().Equal(int(livenessProbe.PeriodSeconds), controllers[0].Health.Interval)
	suite.Assert().Equal(int(livenessProbe.FailureThreshold), controllers[0].Health.RetryCount)
	suite.Assert().Equal(int(livenessProbe.SuccessThreshold), controllers[0].Health.RequiredConsecutivePasses)
	suite.Assert().Equal(int(livenessProbe.TimeoutSeconds), controllers[0].Health.Timeout)
	suite.Assert().Equal(defaultHttpStatusCode, controllers[0].Health.ExpectCode)
}

func (suite *ClientTestSuite) TestGetIngressControllers_healthCheckFromAnnotations_unparseablePort() {
	expectedHealthCheckPath := "/healtz"
	corruptHealthCheckPort := "8088_doh"
	expectedName := "test-ingress-controller"
	expectedNodeSelector := map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	expectedPort := 8080
	expectedServiceGroupTemplate := "svc grp name template"
	defaultHttpStatusCode := "200"
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName(expectedName)
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group":   expectedServiceGroupTemplate,
		"a10.health.endpoint": expectedHealthCheckPath,
		"a10.health.port":     corruptHealthCheckPort,
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = expectedNodeSelector
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_jmx",
				HostPort:      8181,
				ContainerPort: 9191,
			},
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      int32(expectedPort),
				ContainerPort: 9090,
			},
		},
		LivenessProbe: &livenessProbe,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(controllers)

	suite.Assert().NotNil(controllers[0].Health)
	suite.Assert().Equal(expectedHealthCheckPath, controllers[0].Health.Endpoint)
	suite.Assert().Equal(int(livenessProbe.HTTPGet.Port.IntVal), controllers[0].Health.Port)
	suite.Assert().Equal(int(livenessProbe.PeriodSeconds), controllers[0].Health.Interval)
	suite.Assert().Equal(int(livenessProbe.FailureThreshold), controllers[0].Health.RetryCount)
	suite.Assert().Equal(int(livenessProbe.SuccessThreshold), controllers[0].Health.RequiredConsecutivePasses)
	suite.Assert().Equal(int(livenessProbe.TimeoutSeconds), controllers[0].Health.Timeout)
	suite.Assert().Equal(defaultHttpStatusCode, controllers[0].Health.ExpectCode)
}

func (suite *ClientTestSuite) TestGetIngressControllers_apiCallFails() {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("*", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("fail")
	})
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().NotNil(err)
	suite.Assert().Nil(controllers)
}

func (suite *ClientTestSuite) TestGetIngressControllers_missingServiceNameTemplateInAnnotations() {
	annotations := map[string]string{}
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName("test-ingress-controller")
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(annotations)
	daemonSet1.Spec.Template.Spec.NodeSelector = map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_jmx",
				HostPort:      8181,
				ContainerPort: 9191,
			},
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      8080,
				ContainerPort: 9090,
			},
		},
		LivenessProbe: &livenessProbe,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().Equal(0, len(controllers))
}

func (suite *ClientTestSuite) TestGetIngressControllers_mainContainerNotFound() {
	expectedName := "test-ingress-controller"
	expectedNodeSelector := map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName(expectedName)
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group": "svc grp 1",
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = expectedNodeSelector
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_jmx",
				HostPort:      8181,
				ContainerPort: 9191,
			},
		},
		LivenessProbe: &livenessProbe,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().Equal(0, len(controllers))
}

func (suite *ClientTestSuite) TestGetIngressControllers_missingLivenessProbe() {
	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName("test-ingress-controller")
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group": "svc grp 1",
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      8080,
				ContainerPort: 9090,
			},
		},
		LivenessProbe: nil,
	})

	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().Equal(0, len(controllers))
}

func (suite *ClientTestSuite) TestGetIngressControllers_noNameMatchesIngressController() {
	brokenName := "broken-name"
	livenessProbe := corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{IntVal: 8080},
			},
		},
		PeriodSeconds:    15,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		TimeoutSeconds:   10,
	}

	daemonSet1 := extensionsv1beta1.DaemonSet{}
	daemonSet1.SetName(brokenName)
	daemonSet1.SetNamespace("ingress")
	daemonSet1.SetAnnotations(map[string]string{
		"a10.service_group": "cvs grp template",
	})
	daemonSet1.Spec.Template.Spec.NodeSelector = map[string]string{
		"ingress_node": "true",
		"class":        "10g",
	}
	daemonSet1.Spec.Template.Spec.Containers = append(daemonSet1.Spec.Template.Spec.Containers, corev1.Container{
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "jetty_http",
				HostPort:      8080,
				ContainerPort: 9090,
			},
		},
		LivenessProbe: &livenessProbe,
	})
	daemonSetList := extensionsv1beta1.DaemonSetList{
		Items: []extensionsv1beta1.DaemonSet{daemonSet1},
	}
	clientset := fake.NewSimpleClientset(&daemonSetList)
	client := suite.helper.BuildClient(clientset)

	controllers, err := client.GetIngressControllers()

	suite.Assert().Nil(err)
	suite.Assert().Equal(0, len(controllers))
}
