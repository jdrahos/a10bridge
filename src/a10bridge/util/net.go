package util

import (
	"net"
)

var ipResolver IPResolver = new(defaultResolverImpl)

//LookupIP call to get one of the resolved ips as string
func LookupIP(hostname string) (string, error) {
	return ipResolver.LookupIP(hostname)
}

//InjectIPResolver allows injection of custom host resolver and returns the current resolver
func InjectIPResolver(resolver IPResolver) IPResolver {
	old := ipResolver
	ipResolver = resolver
	return old
}

//IPResolver implementations of this interface are responsible for ip lookups
type IPResolver interface {
	LookupIP(hostname string) (string, error)
}

type defaultResolverImpl struct{}

//LookupIP call to get the first looked up ip as string
func (resolver defaultResolverImpl) LookupIP(hostname string) (string, error) {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	return addr[0].String(), err
}
