package model

//ConfigMap data holder
type ConfigMap struct {
	Name      string
	Namespace string
	Data      map[string]string
}
