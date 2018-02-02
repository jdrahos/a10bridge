package model

//Node node information holder
type Node struct {
	Name      string
	A10Server string
	Weight    string
	IPAddress string
	Labels    map[string]string
}

type Nodes []*Node

func (s Nodes) Len() int {
	return len(s)
}
func (s Nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Nodes) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
