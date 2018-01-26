package processor

import (
	"a10bridge/a10/api"
	"a10bridge/model"
	"a10bridge/util"
	"fmt"

	"github.com/golang/glog"
)

//ServiceGroupProcessor processor responsible for processing nodes
type ServiceGroupProcessor interface {
	ProcessServiceGroup(serviceGroup *model.ServiceGroup, failedNodeNames []string) error
}

type serviceGroupProcessorImpl struct {
	a10Client api.Client
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
		if processor.a10Client.IsServiceGroupNotFound(a10err) {
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

		missingMembers := findMissingMembers(serviceGroup.Members, a10ServiceGroup.Members)

		if len(missingMembers) > 0 {
			for _, member := range missingMembers {
				err := processor.a10Client.CreateMember(member)
				if err != nil && !processor.a10Client.IsMemberAlreadyExists(err) {
					glog.Errorf("Failed to create member %s:%d for service group %s. error: %s", member.ServerName, member.Port, member.ServiceGroupName, err)
					a10err = err
				}
			}
		}

		extraMembers := findExtraMembers(serviceGroup.Members, a10ServiceGroup.Members)

		if len(extraMembers) > 0 {
			for _, member := range extraMembers {
				err := processor.a10Client.DeleteMember(member)
				if err != nil {
					glog.Errorf("Failed to delete member %s:%d for service group %s. error: %s", member.ServerName, member.Port, member.ServiceGroupName, err)
					a10err = err
				}
			}
		}
	}

	return a10err
}

func sameGroupConfigs(serviceGroup *model.ServiceGroup, a10ServiceGroup *model.ServiceGroup) bool {
	if serviceGroup.Health.Name != a10ServiceGroup.Health.Name {
		glog.Infof("Health monitor names '%s' and '%s' don't match", serviceGroup.Health.Name, a10ServiceGroup.Health.Name)
		return false
	}

	return true
}

func findExtraMembers(expected []*model.Member, members []*model.Member) []*model.Member {
	extraMembers := make([]*model.Member, 0)

	for _, member := range members {
		if !containsMemeber(expected, member) {
			glog.Infof("'%s' should not be a member", member)
			extraMembers = append(extraMembers, member)
		}
	}

	return extraMembers
}

func findMissingMembers(expected []*model.Member, members []*model.Member) []*model.Member {
	missingMembers := make([]*model.Member, 0)

	for _, member := range expected {
		if !containsMemeber(members, member) {
			glog.Infof("'%s' should be a member", member)
			missingMembers = append(missingMembers, member)
		}
	}

	return missingMembers
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
