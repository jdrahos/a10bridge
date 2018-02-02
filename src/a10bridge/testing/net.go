package testing

import (
	"fmt"

	"github.com/golang/glog"
)

//ConfigurableResolver configurable implementation of HostResolver interface
type ConfigurableResolver struct {
	data map[string]string
}

//LookupIP call to get the first looked up ip as string
func (resolver ConfigurableResolver) LookupIP(hostname string) (string, error) {
	addr, found := resolver.data[hostname]
	if found {
		return addr, nil
	}
	return "", fmt.Errorf("Failed to resolve ip for host %s", hostname)
}

//AddRecord add new record for lookups
func (resolver *ConfigurableResolver) AddRecord(hostname, addr string) {
	resolver.data[hostname] = addr
}

//Reset reset record configuration
func (resolver *ConfigurableResolver) Reset() {
	resolver.data = make(map[string]string)
}

func (resolver *ConfigurableResolver) init() {
	glog.Error("in init")
	resolver.data = make(map[string]string)
}
