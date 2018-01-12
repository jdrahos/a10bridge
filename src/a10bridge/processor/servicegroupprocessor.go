package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
	"a10bridge/util"
	"fmt"

	"github.com/golang/glog"
)

//ServiceGroupProcessor processor responsible for processing nodes
type ServiceGroupProcessor interface {
	BuildServiceGroups(controllers []*model.IngressController, environment *model.Environment) map[string]*model.ServiceGroup
	ProcessServiceGroup(serviceGroup *model.ServiceGroup, failedNodeNames []string) error
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
			serviceGroup.Health.Endpoint = "/syntheticHealth"
			serviceGroup.Health.ExpectCode = "404"
		}
	}

	return serviceGroups
}

func (processor serviceGroupProcessorImpl) ProcessServiceGroup(serviceGroup *model.ServiceGroup, failedNodeNames []string) error {
	glog.Infof("Processing service group %s", serviceGroup.Name /* util.ToJSON(serviceGroup) */)

	members := buildMembers(serviceGroup, failedNodeNames)

	if len(members) == 0 {
		return fmt.Errorf("There were no members found for service group %s", serviceGroup.Name)
	}

	serviceGroup.Members = members

	a10ServiceGroup, a10err := processor.a10Client.GetServiceGroup(serviceGroup.Name)
	if a10err != nil {
		//health monitor not found
		if a10err.Code() == 67305473 {
			a10err = processor.a10Client.CreateServiceGroup(serviceGroup)
		}
	} else {
		fmt.Println(util.ToJSON(a10ServiceGroup))

		if !sameGroupConfigs(serviceGroup, a10ServiceGroup) {
			glog.Info("Service group configuration in a10 differs from configuration in kubernetes, resetting service group in a10")
			a10err = processor.a10Client.UpdateServiceGroup(serviceGroup)
			if a10err != nil {
				return a10err
			}
			glog.Info("A10 Service group configuration synced with kubernetes")
		} else {
			glog.Info("A10 Service group configuration is in sync with kubernetes")
		}
	}

	return a10err
}

func sameGroupConfigs(serviceGroup *model.ServiceGroup, a10ServiceGroup *model.ServiceGroup) bool {
	if serviceGroup.Health.Name != a10ServiceGroup.Health.Name {
		glog.Infof("Health monitor names '%s' and '%s' don't match", serviceGroup.Health.Name, a10ServiceGroup.Health.Name)
		return false
	}

	if len(serviceGroup.Members) != len(a10ServiceGroup.Members) {
		glog.Infof("Numbers of memmbers in kubernetes '%d' and a10 '%d' don't match", len(serviceGroup.Members), len(a10ServiceGroup.Members))
		return false
	}

	for _, member := range serviceGroup.Members {
		if !containsMemeber(a10ServiceGroup.Members, member) {
			glog.Infof("Memeber '%s' is missing in a10", member)
			return false
		}
	}

	return true
}

func containsMemeber(members []*model.Member, lookFor *model.Member) bool {
	for _, item := range members {
		if item.ServerName == lookFor.ServerName && item.Port == lookFor.Port {
			return true
		}
	}
	return false
}

func buildMembers(serviceGroup *model.ServiceGroup, excludedNodeNames []string) []*model.Member {
	members := make([]*model.Member, 0)

	for _, controller := range serviceGroup.IngressControllers {
		port := controller.Port
		for _, node := range controller.Nodes {
			if util.Contains(excludedNodeNames, node.Name) {
				continue
			}
			members = append(members, &model.Member{
				Port:             port,
				ServerName:       node.A10Server,
				ServiceGroupName: serviceGroup.Name,
			})
		}
	}

	return members
}
