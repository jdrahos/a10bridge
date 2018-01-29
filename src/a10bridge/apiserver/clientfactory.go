package apiserver

import (
	"os"
	"path/filepath"

	"github.com/golang/glog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//CreateClient creates kubernetes apiserver client
func CreateClient() (K8sClient, error) {
	//assume we are running inside the pod, if we fail lets try to build kubectl client
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Warningf("Failed to create in cluster config. Error: %v", err)
		return createKubectlClient()
	}

	if err != nil {
		glog.Warningf("Failed to create in cluster client. Error: %v", err)
		return createKubectlClient()
	}

	clientset, err := kubernetes.NewForConfig(config)
	client := newClient(clientset)
	glog.Info("Created in-cluser client")

	return client, err
}

func createKubectlClient() (K8sClient, error) {
	kubectlConfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	glog.Infof("Trying to create kubectl client using %s", kubectlConfigPath)

	var client K8sClient

	config, err := clientcmd.BuildConfigFromFlags("", kubectlConfigPath)
	if err != nil {
		glog.Errorf("Failed to create kubectl client. Error: %v ", err)
		return client, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Failed to create kubectl client. Error: %v ", err)
		return client, err
	}

	client = newClient(clientset)
	glog.Info("Created kubectl client")

	return client, nil
}
