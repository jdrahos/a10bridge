package apiserver

import (
	"a10bridge/model"
	"net"

	"k8s.io/api/core/v1"
)

var netLookupIP = net.LookupIP

//BuildNode builds a apiserver node with relevant information
func buildNode(k8sNode v1.Node) (*model.Node, error) {
	var node model.Node
	name := k8sNode.GetName()
	addr, err := netLookupIP(name)

	if err == nil {
		node = model.Node{
			Name:      name,
			IPAddress: addr[0].String(),
			A10Server: findA10ServerName(k8sNode),
			Weight:    findNodeWeight(k8sNode, "1"),
			Labels:    k8sNode.Labels,
		}
	}

	return &node, err
}

func findNodeWeight(k8sNode v1.Node, defWeight string) string {
	weight, exists := k8sNode.Annotations["a10.server.weight"]
	if !exists {
		weight = defWeight
	}

	return weight
}

func findA10ServerName(k8sNode v1.Node) string {
	serverName, exists := k8sNode.Annotations["a10.server"]
	if !exists {
		serverName = k8sNode.GetName()
	}

	return serverName
}
