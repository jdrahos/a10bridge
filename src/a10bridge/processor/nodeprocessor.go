package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
	"a10bridge/util"
	"fmt"

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
		if processor.a10Client.IsServerNotFound(a10err) {
			a10err = processor.a10Client.CreateServer(node)
			if a10err != nil {
				return a10err
			}
		} else {
			return a10err
		}
	} else {
		fmt.Println(util.ToJSON(server))

		if !isSame(node, server) {
			glog.Infof("Server and node configurations differ, setting the server to ip %s and weight %s", node.IPAddress, node.Weight)
			server.IPAddress = node.IPAddress
			server.Weight = node.Weight
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

func isSame(node *model.Node, server *model.Node) bool {
	if node.IPAddress != server.IPAddress {
		glog.Infof("Server ip addresses '%s' and '%s' don't match", server.IPAddress, node.IPAddress)
		return false
	}

	if node.Weight != server.Weight {
		glog.Infof("Server weights '%s' and '%s' don't match", server.IPAddress, node.IPAddress)
		return false
	}

	return true
}
