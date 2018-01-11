package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
	"errors"
	"strings"
)

//EnvironmentProcessor processor responsible for discovering environment variables
type EnvironmentProcessor interface {
	BuildEnvironment() (*model.Environment, error)
}

type environmentProcessorImpl struct {
	k8sClient apiserver.Client
	a10Client api.Client
}

func (processor environmentProcessorImpl) BuildEnvironment() (*model.Environment, error) {
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
