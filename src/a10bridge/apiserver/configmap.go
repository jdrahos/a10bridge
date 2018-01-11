package apiserver

import (
	"a10bridge/model"

	"k8s.io/api/core/v1"
)

func buildConfigMap(k8sConfigMap v1.ConfigMap) *model.ConfigMap {
	return &model.ConfigMap{
		Name:      k8sConfigMap.GetName(),
		Namespace: k8sConfigMap.GetNamespace(),
		Data:      k8sConfigMap.Data,
	}
}
