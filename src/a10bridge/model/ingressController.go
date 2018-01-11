package model

//IngressController ingress controller data structure
type IngressController struct {
	Name                     string
	NodeSelectors            map[string]string
	Nodes                    []*Node
	ServiceGroupNameTemplate string
	Health                   *HealthCheck
	Port                     int
}
