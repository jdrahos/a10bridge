package apiserver

import (
	"a10bridge/model"
	"strings"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

//Client apis server client
type K8sClient interface {
	GetNodes() ([]*model.Node, error)
	GetConfigMap(namespace string, name string) (*model.ConfigMap, error)
	GetIngressControllers() ([]*model.IngressController, error)
}

type clientImpl struct {
	corev1Impl            corev1.CoreV1Interface
	extensionsv1beta1Impl extensionsv1beta1.ExtensionsV1beta1Interface
}

//New build new client
func newClient(clientset *kubernetes.Clientset) K8sClient {
	return clientImpl{
		corev1Impl:            clientset.CoreV1(),
		extensionsv1beta1Impl: clientset.ExtensionsV1beta1(),
	}
}

//GetNodes get nodes
func (client clientImpl) GetNodes() ([]*model.Node, error) {
	var nodes []*model.Node
	nodeList, err := client.corev1Impl.Nodes().List(metav1.ListOptions{})
	if err == nil {
		for _, node := range nodeList.Items {
			node, builderr := buildNode(node)
			if builderr == nil {
				nodes = append(nodes, node)
			} else {
				err = builderr
				node = nil
				break
			}
		}
	} else {
		nodes = nil
	}

	return nodes, err
}

//GetConfigMap finds config map or returns null
func (client clientImpl) GetConfigMap(namespace string, name string) (*model.ConfigMap, error) {
	var config *model.ConfigMap
	configMapList, err := client.corev1Impl.ConfigMaps(namespace).List(metav1.ListOptions{})
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

func (client clientImpl) GetIngressControllers() ([]*model.IngressController, error) {
	var controllers []*model.IngressController
	controllerList, err := client.extensionsv1beta1Impl.DaemonSets("ingress").List(metav1.ListOptions{})
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
