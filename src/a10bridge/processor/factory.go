package processor

import (
	"a10bridge/a10"
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/config"
)

//A10Processors a10 processors holder
type A10Processors struct {
	Node         NodeProcessor
	ServiceGroup ServiceGroupProcessor
	HealthCheck  HealthCheckProcessor
	client       api.Client
}

func (processors A10Processors) Destroy() {
	processors.client.Close()
}

//BuildK8sProcessor builds kubernetes processor
func BuildK8sProcessor() (K8sProcessor, error) {
	client, err := apiserver.CreateClient()
	if err != nil {
		return nil, err
	}
	return &k8sProcessorImpl{
		k8sClient: client,
	}, nil
}

//BuildA10Processors builds a10 processors
func BuildA10Processors(a10instance *config.A10Instance) (A10Processors, error) {
	var processors A10Processors
	a10Client, err := a10.BuildClient(a10instance)
	if err != nil {
		return processors, err
	}

	return A10Processors{
		Node: &nodeProcessorImpl{
			a10Client: a10Client,
		},

		ServiceGroup: &serviceGroupProcessorImpl{
			a10Client: a10Client,
		},

		HealthCheck: &healthCheckProcessorImpl{
			a10Client: a10Client,
		},

		client: a10Client,
	}, err
}
