package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/model"
)

//IngressControllerProcessor processor responsible for processing ingresses
type IngressControllerProcessor interface {
	FindIngressControllers() ([]*model.IngressController, error)
}

type ingressControllerProcessorImpl struct {
	k8sClient apiserver.Client
	a10Client api.Client
}

func (processor ingressControllerProcessorImpl) FindIngressControllers() ([]*model.IngressController, error) {
	return processor.k8sClient.GetIngressControllers()
}
