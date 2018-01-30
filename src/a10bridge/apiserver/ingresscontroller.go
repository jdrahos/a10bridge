package apiserver

import (
	"a10bridge/model"
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func buildIngressController(controller v1beta1.DaemonSet) (*model.IngressController, error) {
	serviceGroup, exists := controller.Annotations["a10.service_group"]
	if !exists {
		return nil, fmt.Errorf("Missing service group name tamplate on ingress controller %s", controller.GetName())
	}

	mainContainer, httpPort := findMainContainer(controller.Spec.Template.Spec.Containers)
	if mainContainer == nil {
		return nil, fmt.Errorf("Unable to find main container for ingress controller %s", controller.GetName())
	}

	healthCheck, err := buildHealthCheck(controller, mainContainer)
	if err != nil {
		return nil, fmt.Errorf("Failed to build health check for ingress controller %s", controller.GetName())
	}

	return &model.IngressController{
		Name:          controller.GetName(),
		NodeSelectors: controller.Spec.Template.Spec.NodeSelector,
		Health:        healthCheck,
		Port:          httpPort,
		ServiceGroupNameTemplate: serviceGroup,
	}, err
}

func findMainContainer(containers []v1.Container) (*v1.Container, int) {
	for _, container := range containers {
		if container.Ports == nil || len(container.Ports) == 0 {
			continue
		}

		for _, port := range container.Ports {
			if strings.HasSuffix(port.Name, "http") {
				return &container, int(port.HostPort)
			}
		}
	}

	return nil, 0
}
