package main

import (
	"a10bridge/a10"
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/args"
	"a10bridge/model"
	"a10bridge/processor"
	"a10bridge/util"

	"github.com/golang/glog"
)

var k8sClient apiserver.Client
var a10Client api.Client
var arguments args.Args

func initialize() error {
	var err error
	arguments, err = args.Build()
	if err != nil {
		return err
	}

	k8sClient, err = apiserver.CreateClient()
	if err != nil {
		return err
	}

	a10Client, err = a10.BuildClient(arguments)
	if err != nil {
		return err
	}

	return nil
}

func cleanup() {
	a10Client.Close()
	defer glog.Flush()
}

func main() {
	err := initialize()
	defer cleanup()
	if err != nil {
		glog.Errorf("Failed to initialize the application. error: %s", err)
		return
	}

	processors := processor.Build(k8sClient, a10Client)

	environment, err := processors.Environment.BuildEnvironment()
	if err != nil {
		glog.Errorf("Failed to build environment. error: %s", err)
		return
	}
	glog.Infof("Using environment: %s", util.ToJSON(environment))

	controllers, err := processors.IngressController.FindIngressControllers()
	if err != nil {
		glog.Errorf("Failed to get ingress controllers. error: %s", err)
		return
	}
	glog.Infof("Ingress controllers: %s", util.ToJSON(controllers))

	nodesMap := make(map[string]*model.Node)

	for _, controller := range controllers {
		glog.Infof("Looking up nodes for ingress controller %s", controller.Name)
		nodes, err := processors.Node.FindNodes(controller.NodeSelectors)
		if err != nil {
			glog.Errorf("Failed to get nodes for controller %s. error: %s", controller.Name, err)
			continue
		}

		for _, node := range nodes {
			controller.Nodes = append(controller.Nodes, node)
			nodesMap[node.Name] = node
		}
	}

	glog.Info("Making sure servers in a10 are in sync with ingress nodes")

	failedNodeNames := make([]string, 10)

	for nodeName, node := range nodesMap {
		err := processors.Node.ProcessNode(node)
		if err != nil {
			glog.Errorf("Failed to process node %s. error: %s", nodeName, err)
			failedNodeNames = append(failedNodeNames, nodeName)
		}
	}

	glog.Info("Generating service groups based on ingress controllers")

	serviceGroups := processors.ServiceGroup.BuildServiceGroups(controllers, environment)
	glog.Infof("Service groups: %s", util.ToJSON(serviceGroups))

	glog.Info("Processing service groups")

	for serviceGroupName, serviceGroup := range serviceGroups {
		err := processors.HealthCheck.ProcessHealthCheck(serviceGroup.Health)
		if err != nil {
			glog.Errorf("Failed to process health check %s, error: %s", serviceGroupName, err)
			continue
		}
	}
}
