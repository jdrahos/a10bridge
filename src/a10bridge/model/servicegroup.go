package model

type ServiceGroup struct {
	Name               string
	Health             *HealthCheck
	IngressControllers []*IngressController
	Members            []*Member
}
