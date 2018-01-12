package api

import (
	"a10bridge/model"
)

//Client a10 client
type Client interface {
	Close() A10Error

	GetServer(serverName string) (*model.Server, A10Error)
	CreateServer(server *model.Server) A10Error
	UpdateServer(server *model.Server) A10Error

	GetHealthMonitor(monitorName string) (*model.HealthCheck, A10Error)
	CreateHealthMonitor(monitor *model.HealthCheck) A10Error
	UpdateHealthMonitor(monitor *model.HealthCheck) A10Error

	GetServiceGroup(serviceGroupName string) (*model.ServiceGroup, A10Error)
	CreateServiceGroup(serviceGroup *model.ServiceGroup) A10Error
	UpdateServiceGroup(serviceGroup *model.ServiceGroup) A10Error
}
