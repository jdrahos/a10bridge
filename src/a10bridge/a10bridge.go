package main

import (
	"a10bridge/config"
	"a10bridge/model"
	"a10bridge/processor"
	"a10bridge/util"
	"os"
	"sort"
	"time"

	"github.com/golang/glog"
)

var processorBuildK8sProcessor = processor.BuildK8sProcessor
var configBuildConfig = config.BuildConfig
var processorBuildA10Processors = processor.BuildA10Processors

type exitCode = int

const (
	Normal                     exitCode = 0
	FailedToBuildConfig        exitCode = 1
	ExcutionTimedOut           exitCode = 2
	FailedToBuildExpectedState exitCode = 3
	FailedToProcessA10Instance exitCode = 4
)

func main() {
	os.Exit(mainInternal())
}

func mainInternal() exitCode {
	defer glog.Flush()
	context, err := configBuildConfig()
	if err != nil {
		glog.Errorf("Failed to initialize the application. error: %s", err)
		return FailedToBuildConfig
	}

	done := make(chan exitCode)
	interval := time.Second * time.Duration(*context.Arguments.Interval)
	ticker := time.NewTicker(interval)
	executionFunc := func() exitCode {
		glog.Info("The execution is starting")
		exitCodeChan := make(chan exitCode)
		go func() {
			exitCodeChan <- reconcile(context)
		}()
		select {
		case code := <-exitCodeChan:
			glog.Info("The execution has finished")
			return code
		case <-time.After(interval):
			glog.Error("The execution has timed out")
			return ExcutionTimedOut
		}
	}

	if *context.Arguments.Daemon {
		go func() {
			code := executionFunc()
			if code == Normal {
				for _ = range ticker.C {
					code = executionFunc()
					if code != Normal {
						break
					}
				}
			}
			done <- code
		}()
		return <-done
	}

	return executionFunc()
}

func reconcile(context *config.RunContext) exitCode {
	serviceGroups, nodesMap, err := buildexpectedState(context)
	if err != nil {
		glog.Errorf("Failed to build expected state by inspecting kubernetes configuration. error: %s", err)
		return FailedToBuildExpectedState
	}

	if *context.Arguments.Sort {
		sort.Sort(context.A10Instances)
	}

	result := Normal

	for _, a10Instance := range context.A10Instances {
		err := processContext(context, &a10Instance, serviceGroups, nodesMap)
		if err != nil {
			glog.Errorf("Failed to process context for a10 server %s. error: %s", a10Instance.Name, err)
			result = FailedToProcessA10Instance
		}
	}

	return result
}

func buildexpectedState(context *config.RunContext) (map[string]*model.ServiceGroup, map[string]*model.Node, error) {
	var serviceGroups map[string]*model.ServiceGroup
	nodesMap := make(map[string]*model.Node)
	k8sProcessor, err := processorBuildK8sProcessor()
	if err != nil {
		glog.Errorf("Failed to build kubernetes processor. error: %s", err)
		return serviceGroups, nodesMap, err
	}

	environment, err := k8sProcessor.BuildEnvironment()
	if err != nil {
		glog.Errorf("Failed to build environment. error: %s", err)
		return serviceGroups, nodesMap, err
	}
	glog.Infof("Using environment: %s", util.ToJSON(environment))

	controllers, err := k8sProcessor.FindIngressControllers()
	if err != nil {
		glog.Errorf("Failed to get ingress controllers. error: %s", err)
		return serviceGroups, nodesMap, err
	}

	glog.Infof("Ingress controllers: %s", util.ToJSON(controllers))

	for _, controller := range controllers {
		glog.Infof("Looking up nodes for ingress controller %s", controller.Name)
		nodes, err := k8sProcessor.FindNodes(controller.NodeSelectors)
		if err != nil {
			glog.Errorf("Failed to get nodes for controller %s. error: %s", controller.Name, err)
			return serviceGroups, nodesMap, err
		}

		for _, node := range nodes {
			controller.Nodes = append(controller.Nodes, node)
			nodesMap[node.Name] = node
		}
	}

	glog.Info("Generating service groups based on ingress controllers")

	serviceGroups = k8sProcessor.BuildServiceGroups(controllers, environment)
	glog.Infof("Service groups: %s", util.ToJSON(serviceGroups))

	return serviceGroups, nodesMap, nil
}

func processContext(context *config.RunContext, a10instance *config.A10Instance, serviceGroups map[string]*model.ServiceGroup, nodesMap map[string]*model.Node) error {
	nodesSlice := make(model.Nodes, 0)
	for _, node := range nodesMap {
		nodesSlice = append(nodesSlice, node)
	}
	serviceGroupSlice := make(model.ServiceGroups, 0)
	for _, serviceGroup := range serviceGroups {
		serviceGroupSlice = append(serviceGroupSlice, serviceGroup)
	}
	if *context.Arguments.Sort {
		sort.Sort(nodesSlice)
		sort.Sort(serviceGroupSlice)
	}

	glog.Infof("Processing context for a10 load balancer %s", a10instance.Name)
	processors, err := processorBuildA10Processors(a10instance)
	if err != nil {
		return err
	}
	defer processors.Destroy()
	glog.Info("Making sure servers in a10 are in sync with ingress nodes")
	failedNodeNames := make([]string, 0)
	for _, node := range nodesSlice {
		err := processors.Node.ProcessNode(node)
		if err != nil {
			glog.Errorf("Failed to process node %s. error: %s", node.Name, err)
			failedNodeNames = append(failedNodeNames, node.Name)
		}
	}

	glog.Info("Processing service groups")
	for _, serviceGroup := range serviceGroupSlice {
		err := processors.HealthCheck.ProcessHealthCheck(serviceGroup.Health)
		if err != nil {
			glog.Errorf("Failed to process health check %s, error: %s", serviceGroup.Name, err)
			continue
		}

		err = processors.ServiceGroup.ProcessServiceGroup(serviceGroup, failedNodeNames)
		if err != nil {
			glog.Errorf("Failed to process service group %s, error: %s", serviceGroup.Name, err)
			continue
		}
	}

	glog.Infof("Done processing context for a10 load balancer %s", a10instance.Name)
	return nil
}
