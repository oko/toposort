package toposort

// Edge represents an edge from one node to another
type Edge struct {
	From Node
	To   Node
}

// EdgesFromNode represents a set of Edges leaving a node
type EdgesFromNode map[string]Edge

// Copy returns a copy of EdgesFromNode
func (ne EdgesFromNode) Copy() EdgesFromNode {
	e := make(EdgesFromNode)
	for k, v := range ne {
		e[k] = v
	}
	return e
}

// Edges represents the set of all Edges in a topology
// grouped by origin node
type Edges map[string]EdgesFromNode

// AddEdge adds an edge to this edge set
func (e Edges) AddEdge(from, to Node) {
	if _, ok := e[from.Id()]; !ok {
		e[from.Id()] = make(EdgesFromNode)
	}
	e[from.Id()][to.Id()] = Edge{from, to}
}

// HasIncoming returns whether a node has incoming Edges
func (e Edges) HasIncoming(n Node) bool {
	for _, ne := range e {
		if ne == nil {
			continue
		}
		if _, ok := ne[n.Id()]; ok {
			return true
		}
	}
	return false
}

// Count returns the total count of Edges
func (e Edges) Count() int {
	count := 0
	for _, ne := range e {
		if ne == nil {
			continue
		}
		count += len(ne)
	}
	return count
}

// Copy returns a copy of the entire Edges structure
func (e Edges) Copy() Edges {
	c := make(Edges)
	for k, v := range e {
		c[k] = v.Copy()
	}
	return c
}
