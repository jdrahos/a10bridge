package apiserver

import (
	"a10bridge/model"
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func buildHealthCheck(controllerDaemonSet v1beta1.DaemonSet, mainContainer *v1.Container) (*model.HealthCheck, error) {
	var port int
	var err error

	endpoint, endpointFound := controllerDaemonSet.Annotations["a10.health.endpoint"]
	if !endpointFound {
		glog.Infof("health endpoint annotation not found for ingress controller %s, going to use liveness probe", controllerDaemonSet.GetName())
	}
	portStr, portFound := controllerDaemonSet.Annotations["a10.health.port"]
	if !portFound {
		glog.Infof("health port annotation not found for ingress controller %s, going to use liveness probe", controllerDaemonSet.GetName())
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			glog.Warningf("Failed to convert port annotation %s to int value with error %s, going to use liveness probe", portStr, err)
			portFound = false
		}
	}

	livenessProbe := mainContainer.LivenessProbe
	if livenessProbe == nil {
		return nil, fmt.Errorf("Liveness probe not found on ingress controller %s on container %s", controllerDaemonSet.Name, mainContainer.Name)
	}

	if !endpointFound {
		endpoint = livenessProbe.HTTPGet.Path
	}
	if !portFound {
		port = int(livenessProbe.HTTPGet.Port.IntVal)
	}

	return &model.HealthCheck{
		Endpoint:                  endpoint,
		Port:                      port,
		Interval:                  int(livenessProbe.PeriodSeconds),
		RetryCount:                int(livenessProbe.FailureThreshold),
		RequiredConsecutivePasses: int(livenessProbe.SuccessThreshold),
		Timeout:                   int(livenessProbe.TimeoutSeconds),
		ExpectCode:                "200",
	}, nil
}
