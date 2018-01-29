package processor

import (
	"a10bridge/apiserver"
	"a10bridge/model"
	"a10bridge/util"
	"errors"
	"strings"

	"github.com/golang/glog"
)

var utilApplyTemplate = util.ApplyTemplate

type K8sProcessor interface {
	BuildEnvironment() (*model.Environment, error)
	FindNodes(nodeSelectors map[string]string) ([]*model.Node, error)
	FindIngressControllers() ([]*model.IngressController, error)
	BuildServiceGroups(controllers []*model.IngressController, environment *model.Environment) map[string]*model.ServiceGroup
}

type k8sProcessorImpl struct {
	k8sClient apiserver.K8sClient
}

func (processor k8sProcessorImpl) BuildEnvironment() (*model.Environment, error) {
	config, err := processor.k8sClient.GetConfigMap("ingress", "cluster-configs")
	if err != nil {
		return nil, err
	}
	clusterName, exists := config.Data["name"]
	if !exists || len(clusterName) == 0 {
		return nil, errors.New("Cluster name not found in config map")
	}

	parts := strings.Split(clusterName, "-")
	dataCenter := parts[0]
	clusterType := ""

	if len(parts) > 1 {
		clusterType = parts[1]
	}

	return &model.Environment{
		Cluster:    clusterName,
		DataCenter: dataCenter,
		Type:       clusterType,
	}, nil
}

func (processor k8sProcessorImpl) FindNodes(nodeSelectors map[string]string) ([]*model.Node, error) {
	nodes, err := processor.k8sClient.GetNodes()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("Found %d nodes", len(nodes))
	var matchingNodes []*model.Node

	for _, node := range nodes {
		matches := true
		for label, val := range nodeSelectors {
			k8sValue, exists := node.Labels[label]
			if !exists || k8sValue != val {
				matches = false
				break
			}
		}

		if matches {
			matchingNodes = append(matchingNodes, node)
		}
	}

	glog.Infof("Found %d matching nodes", len(matchingNodes))
	return matchingNodes, err
}

func (processor k8sProcessorImpl) FindIngressControllers() ([]*model.IngressController, error) {
	return processor.k8sClient.GetIngressControllers()
}

func (processor k8sProcessorImpl) BuildServiceGroups(controllers []*model.IngressController, environment *model.Environment) map[string]*model.ServiceGroup {
	serviceGroups := make(map[string]*model.ServiceGroup)

	for _, controller := range controllers {
		serviceGroupName, err := utilApplyTemplate(environment, controller.ServiceGroupNameTemplate)
		if err != nil {
			glog.Errorf("Failed to build service group name for ingress controller %s. error: %s", controller.Name, err)
			continue
		}
		glog.Infof("Ingress controller %s belongs to service group %s", controller.Name, serviceGroupName)
		serviceGroup, existed := serviceGroups[serviceGroupName]
		if !existed {
			healthCheck := *controller.Health
			healthCheck.Name = serviceGroupName
			serviceGroup := model.ServiceGroup{
				Health:             &healthCheck,
				Name:               serviceGroupName,
				IngressControllers: []*model.IngressController{controller},
			}
			serviceGroups[serviceGroupName] = &serviceGroup
		} else {
			serviceGroup.IngressControllers = append(serviceGroup.IngressControllers, controller)
		}
	}

	for _, serviceGroup := range serviceGroups {
		if len(serviceGroup.IngressControllers) > 1 {
			//we will need to fall back to ingress controller's serving port and just check something is replying with 404
			serviceGroup.Health.Endpoint = "/syntheticHealth"
			serviceGroup.Health.ExpectCode = "404"
		}
	}

	return serviceGroups
}
