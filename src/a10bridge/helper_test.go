package main

import (
	"a10bridge/config"
	"a10bridge/processor"
	"sync"

	"github.com/golang/glog"
)

type TestHelper struct{}

type BuildK8sProcessorFunc func() (processor.K8sProcessor, error)
type BuildConfigFunc func() (*config.RunContext, error)
type BuildA10ProcessorsFunc func(a10instance *config.A10Instance) (*processor.A10Processors, error)

var syncMutex = new(sync.Mutex)

func (helper TestHelper) SetBuildK8sProcessorFunc(replacement BuildK8sProcessorFunc) BuildK8sProcessorFunc {
	syncMutex.Lock()
	old := processorBuildK8sProcessor
	processorBuildK8sProcessor = replacement
	syncMutex.Unlock()
	return old
}

func (helper TestHelper) SetBuildConfigFunc(replacement BuildConfigFunc) BuildConfigFunc {
	old := configBuildConfig
	configBuildConfig = replacement
	return old
}

func (helper TestHelper) SetBuildA10ProcessorsFunc(replacement BuildA10ProcessorsFunc) BuildA10ProcessorsFunc {
	glog.Error("Replacing  BuildA10Processors function")
	old := processorBuildA10Processors
	processorBuildA10Processors = replacement
	return old
}
