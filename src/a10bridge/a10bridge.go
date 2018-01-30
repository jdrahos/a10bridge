package main

import (
	"a10bridge/config"
	"a10bridge/model"
	"a10bridge/processor"
	"a10bridge/util"
	"time"

	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	context, err := config.BuildConfig()
	if err != nil {
		glog.Errorf("Failed to initialize the application. error: %s", err)
		return
	}

	done := make(chan bool)
	interval := time.Minute * time.Duration(*context.Arguments.Interval)
	ticker := time.NewTicker(interval)
	executionFunc := func() bool {
		glog.Info("The execution is starting")
		tout := make(chan bool)
		go func() {
			reconcile(context)
			tout <- false
		}()
		select {
		case <-tout:
			glog.Info("The execution has finished")
			return true
		case <-time.After(interval):
			glog.Error("The execution has timed out")
			return false
		}
	}

	go func() {
		if executionFunc() {
			for _ = range ticker.C {
				if !executionFunc() {
					break
				}
			}
		}
		done <- true
	}()
	<-done
}

func reconcile(context *config.RunContext) {
	serviceGroups, nodesMap, err := buildexpectedState(context)
	if err != nil {
		glog.Errorf("Failed to build expected state by inspecting kubernetes configuration. error: %s", err)
		return
	}

	for _, a10Instance := range context.A10Instances {
		err := processContext(&a10Instance, serviceGroups, nodesMap)
		if err != nil {
			glog.Errorf("Failed to process context for a10 server %s. error: %s", a10Instance.Name, err)
		}
	}
}

func buildexpectedState(context *config.RunContext) (map[string]*model.ServiceGroup, map[string]*model.Node, error) {
	var serviceGroups map[string]*model.ServiceGroup
	nodesMap := make(map[string]*model.Node)
	k8sProcessor, err := processor.BuildK8sProcessor()
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
			continue
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

func processContext(a10instance *config.A10Instance, serviceGroups map[string]*model.ServiceGroup, nodesMap map[string]*model.Node) error {
	glog.Infof("Processing context for a10 load balancer %s", a10instance.Name)
	processors, err := processor.BuildA10Processors(a10instance)
	if err != nil {
		return err
	}
	defer processors.Destroy()
	glog.Info("Making sure servers in a10 are in sync with ingress nodes")
	failedNodeNames := make([]string, 0)

	for nodeName, node := range nodesMap {
		err := processors.Node.ProcessNode(node)
		if err != nil {
			glog.Errorf("Failed to process node %s. error: %s", nodeName, err)
			failedNodeNames = append(failedNodeNames, nodeName)
		}
	}

	glog.Info("Processing service groups")

	for serviceGroupName, serviceGroup := range serviceGroups {
		err := processors.HealthCheck.ProcessHealthCheck(serviceGroup.Health)
		if err != nil {
			glog.Errorf("Failed to process health check %s, error: %s", serviceGroupName, err)
			continue
		}

		err = processors.ServiceGroup.ProcessServiceGroup(serviceGroup, failedNodeNames)
		if err != nil {
			glog.Errorf("Failed to process service group %s, error: %s", serviceGroupName, err)
			continue
		}
	}

	glog.Infof("Done processing context for a10 load balancer %s", a10instance.Name)
	return nil
}
