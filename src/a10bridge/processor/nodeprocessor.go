package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
	"a10bridge/util"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

//NodeProcessor processor responsible for processing nodes
type NodeProcessor interface {
	FindNodes(nodeSelectors map[string]string) ([]*model.Node, error)
	ProcessNode(node *model.Node) error
}

type nodeProcessorImpl struct {
	k8sClient apiserver.Client
	a10Client api.Client
}

func (processor nodeProcessorImpl) FindNodes(nodeSelectors map[string]string) ([]*model.Node, error) {
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

func (processor nodeProcessorImpl) ProcessNode(node *model.Node) error {
	glog.Infof("Processing node %s", util.ToJSON(node))

	server, a10err := processor.a10Client.GetServer(node.Name)
	if a10err != nil {
		//server not found
		if a10err.Code() == 67174402 {
			server, converr := buildServer(node)
			if converr != nil {
				return converr
			}
			a10err = processor.a10Client.CreateServer(server)
			if a10err != nil {
				return a10err
			}
		} else {
			return a10err
		}
	} else {
		fmt.Println(util.ToJSON(server))

		nodeIP := node.IPAddress.String()
		nodeWeight, converr := strconv.Atoi(node.Weight)
		if converr != nil {
			return converr
		}

		if nodeIP != server.IP || nodeWeight != server.Weight {
			glog.Infof("Server and node configurations differ, setting the server to ip %s and weight %s", nodeIP, nodeWeight)
			server.IP = nodeIP
			server.Weight = nodeWeight
			a10err = processor.a10Client.UpdateServer(server)
			if a10err != nil {
				return a10err
			}
			glog.Info("Server configuration synced with node configuration")
		} else {
			glog.Info("Server and node configurations are in sync")
		}
	}

	return a10err
}

func buildServer(node *model.Node) (*model.Server, error) {
	var err error
	server := model.Server{}
	server.IP = node.IPAddress.String()
	server.Name = node.Name
	server.Weight, err = strconv.Atoi(node.Weight)

	return &server, err
}
