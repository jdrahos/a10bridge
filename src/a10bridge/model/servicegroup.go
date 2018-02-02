package model

type ServiceGroup struct {
	Name               string
	Health             *HealthCheck
	IngressControllers []*IngressController
	Members            []*Member
}

type ServiceGroups []*ServiceGroup

func (s ServiceGroups) Len() int {
	return len(s)
}
func (s ServiceGroups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ServiceGroups) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
