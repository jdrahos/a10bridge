package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
	"a10bridge/util"

	"github.com/golang/glog"
)

//NodeProcessor processor responsible for processing nodes
type ServiceGroupProcessor interface {
	BuildServiceGroups(controllers []*model.IngressController, environment *model.Environment) map[string]*model.ServiceGroup
}

type serviceGroupProcessorImpl struct {
	k8sClient apiserver.Client
	a10Client api.Client
}

func (processor serviceGroupProcessorImpl) BuildServiceGroups(controllers []*model.IngressController, environment *model.Environment) map[string]*model.ServiceGroup {
	serviceGroups := make(map[string]*model.ServiceGroup)

	for _, controller := range controllers {
		serviceGroupName, err := util.ApplyTemplate(environment, controller.ServiceGroupNameTemplate, controller.ServiceGroupNameTemplate)
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
			serviceGroup.Health.Endpoint = "/syntheticHealh"
			serviceGroup.Health.ExpectCode = "404"
		}
	}

	return serviceGroups
}
