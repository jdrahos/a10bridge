package apiserver

import (
	"k8s.io/client-go/kubernetes"
)

type TestHelper struct{}

func (helper *TestHelper) BuildClient(clientset *kubernetes.Clientset) *Client {
	return &Client{
		clientset: clientset,
	}
}
