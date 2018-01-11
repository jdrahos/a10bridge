package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
)

//Processors processors holder
type Processors struct {
	Environment       EnvironmentProcessor
	IngressController IngressControllerProcessor
	Node              NodeProcessor
	ServiceGroup      ServiceGroupProcessor
	HealthCheck       HealthCheckProcessor
}

//Build builds processors
func Build(k8sclient apiserver.Client, a10Client api.Client) Processors {
	return Processors{
		Environment: &environmentProcessorImpl{
			a10Client: a10Client,
			k8sClient: k8sclient,
		},

		IngressController: &ingressControllerProcessorImpl{
			a10Client: a10Client,
			k8sClient: k8sclient,
		},

		Node: &nodeProcessorImpl{
			a10Client: a10Client,
			k8sClient: k8sclient,
		},

		ServiceGroup: &serviceGroupProcessorImpl{
			a10Client: a10Client,
			k8sClient: k8sclient,
		},

		HealthCheck: &healthCheckProcessorImpl{
			a10Client: a10Client,
			k8sClient: k8sclient,
		},
	}
}
