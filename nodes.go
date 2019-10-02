package toposort

// Node is the interface type for a single node in a topology
type Node interface {
	Id() string
}

// Nodes represents the set of Nodes in a topology
type Nodes map[string]Node

// Add adds a node to this Nodes set
func (n Nodes) Add(node Node) error {
	if n.Has(node) {
		return ErrNodeExists
	}
	n[node.Id()] = node
	return nil
}

// Has checks whether a given node is in a Nodes set
func (n Nodes) Has(node Node) bool {
	_, ok := n[node.Id()]
	return ok
}
