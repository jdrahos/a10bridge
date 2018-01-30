package processor

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/config"
)

type TestHelper struct{}
type ApiserverCreateClientFunc func() (apiserver.K8sClient, error)
type A10BuildClientFunc func(a10Instance *config.A10Instance) (api.Client, api.A10Error)
type UtilApplyTemplateFunc func(data interface{}, tpl string) (string, error)

func (helper TestHelper) SetApiserverCreateClient(createClientFunc ApiserverCreateClientFunc) ApiserverCreateClientFunc {
	old := apiserverCreateClient
	apiserverCreateClient = createClientFunc
	return old
}

func (helper TestHelper) SetA10BuildClient(buildClientFunc A10BuildClientFunc) A10BuildClientFunc {
	old := a10BuildClient
	a10BuildClient = buildClientFunc
	return old
}

func (helper TestHelper) SetUtilApplyTemplate(utilApplyTemplateFunc UtilApplyTemplateFunc) UtilApplyTemplateFunc {
	old := utilApplyTemplate
	utilApplyTemplate = utilApplyTemplateFunc
	return old
}

func (helper TestHelper) BuildNodeProcessor(client api.Client) NodeProcessor {
	return nodeProcessorImpl{a10Client: client}
}

func (helper TestHelper) BuildHealthcheckProcessor(client api.Client) HealthCheckProcessor {
	return healthCheckProcessorImpl{a10Client: client}
}

func (helper TestHelper) BuildServiceGroupProcessor(client api.Client) ServiceGroupProcessor {
	return serviceGroupProcessorImpl{a10Client: client}
}

func (helper TestHelper) BuildK8sProcessor(client apiserver.K8sClient) K8sProcessor {
	return k8sProcessorImpl{k8sClient: client}
}
