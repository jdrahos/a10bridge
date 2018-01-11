package model

import "net"

//Node node information holder
type Node struct {
	Name      string
	A10Server string
	Weight    string
	IPAddress net.IP
	Labels    map[string]string
}
