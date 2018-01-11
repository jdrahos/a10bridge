package apiserver

import (
	"a10bridge/model"
	"strings"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Client apis server client
type Client struct {
	clientset *kubernetes.Clientset
}

//New build new client
func newClient(clientset *kubernetes.Clientset) Client {
	return Client{
		clientset: clientset,
	}
}

//GetNodes get nodes
func (client Client) GetNodes() ([]*model.Node, error) {
	var nodes []*model.Node
	nodeList, err := client.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err == nil {
		nodes = make([]*model.Node, len(nodeList.Items))

		for idx, node := range nodeList.Items {
			nodes[idx], err = buildNode(node)
		}
	}

	return nodes, err
}

//GetConfigMap finds config map or returns null
func (client Client) GetConfigMap(namespace string, name string) (*model.ConfigMap, error) {
	var config *model.ConfigMap
	configMapList, err := client.clientset.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, configMap := range configMapList.Items {
		if configMap.GetName() == name {
			config = buildConfigMap(configMap)
			break
		}
	}

	return config, err
}

func (client Client) GetIngressControllers() ([]*model.IngressController, error) {
	var controllers []*model.IngressController
	controllerList, err := client.clientset.ExtensionsV1beta1().DaemonSets("ingress").List(metav1.ListOptions{})
	if err == nil {
		for _, controller := range controllerList.Items {
			if strings.HasSuffix(controller.GetName(), "ingress-controller") {
				ingressController, err := buildIngressController(controller)
				if err != nil {
					glog.Errorf("Failed to build ingress controller %s. error: %s", controller.GetName(), err)
					continue
				}
				controllers = append(controllers, ingressController)
			}

		}
	}
	return controllers, err
}
