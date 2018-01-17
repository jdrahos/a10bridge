package processor

import (
	"a10bridge/a10/api"
	"a10bridge/model"
	"a10bridge/util"
	"fmt"

	"github.com/golang/glog"
)

//HealthCheckProcessor processor responsible for processing ingresses
type HealthCheckProcessor interface {
	ProcessHealthCheck(healthCheck *model.HealthCheck) error
}

type healthCheckProcessorImpl struct {
	a10Client api.Client
}

func (processor healthCheckProcessorImpl) ProcessHealthCheck(healthCheck *model.HealthCheck) error {
	glog.Infof("Processing healht check %s", util.ToJSON(healthCheck))

	healthMonitor, a10err := processor.a10Client.GetHealthMonitor(healthCheck.Name)
	if a10err != nil {
		//health monitor not found
		if processor.a10Client.IsHealthMonitorNotFound(a10err) {
			healthMonitor = healthCheck
			a10err = processor.a10Client.CreateHealthMonitor(healthMonitor)
			if a10err != nil {
				return a10err
			}
		} else {
			return a10err
		}
	} else {
		fmt.Println(util.ToJSON(healthMonitor))

		if !sameHealthConfigs(healthCheck, healthMonitor) {
			glog.Info("Health monitor configuration in a10 differs from healthcheck configuration in kubernetes, resetting monitor in a10")
			a10err = processor.a10Client.UpdateHealthMonitor(healthCheck)
			if a10err != nil {
				return a10err
			}
			glog.Info("Health monitor configuration synced with kubernetes configuration")
		} else {
			glog.Info("Health monitor configuration is in sync with kubernetes configuration")
		}
	}

	return a10err
}

func sameHealthConfigs(healthCheck *model.HealthCheck, healthMonitor *model.HealthCheck) bool {
	if healthCheck.Endpoint != healthMonitor.Endpoint {
		glog.Infof("Endpoints '%s' and '%s' don't match", healthCheck.Endpoint, healthMonitor.Endpoint)
		return false
	}
	if healthCheck.ExpectCode != healthMonitor.ExpectCode {
		glog.Infof("Expected codes '%s' and '%s' don't match", healthCheck.ExpectCode, healthMonitor.ExpectCode)
		return false
	}
	if healthCheck.Interval != healthMonitor.Interval {
		glog.Infof("Intervals '%s' and '%s' don't match", healthCheck.Interval, healthMonitor.Interval)
		return false
	}
	if healthCheck.Port != healthMonitor.Port {
		glog.Infof("Ports '%d' and '%d' don't match", healthCheck.Port, healthMonitor.Port)
		return false
	}
	if healthCheck.RequiredConsecutivePasses != healthMonitor.RequiredConsecutivePasses {
		glog.Infof("Required consecutive passes '%d' and '%d' don't match", healthCheck.RequiredConsecutivePasses, healthMonitor.RequiredConsecutivePasses)
		return false
	}
	if healthCheck.RetryCount != healthMonitor.RetryCount {
		glog.Infof("Retry counts '%d' and '%d' don't match", healthCheck.RetryCount, healthMonitor.RetryCount)
		return false
	}
	if healthCheck.Timeout != healthMonitor.Timeout {
		glog.Infof("Timeouts '%d' and '%d' don't match", healthCheck.Timeout, healthMonitor.Timeout)
		return false
	}

	return true
}
