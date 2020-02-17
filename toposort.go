package toposort

import (
	"errors"

	"github.com/emirpasic/gods/utils"

	"github.com/emirpasic/gods/trees/binaryheap"
)

var (
	ErrNodeExists       = errors.New("node already exists in topology")
	ErrNodeDoesNotExist = errors.New("node does not exist in topology")
	ErrRuntimeExceeded  = errors.New("sort runtime exceeded bound")
)

type ErrCycleInTopology struct {
	OriginalEdges  Edges
	RemainingEdges Edges
}

func (e *ErrCycleInTopology) Error() string {
	return "cycle in topology"
}

// Topology represents an entire graph topology
type Topology struct {
	Nodes Nodes
	Edges Edges
}

// NewTopology returns a new topology
func NewTopology() *Topology {
	return &Topology{
		Nodes: make(Nodes),
		Edges: make(Edges),
	}
}

func (t *Topology) AddNode(n Node) error {
	return t.Nodes.Add(n)
}

func (t *Topology) AddEdge(from, to Node) error {
	if !t.Nodes.Has(from) || !t.Nodes.Has(to) {
		return ErrNodeDoesNotExist
	}
	t.Edges.AddEdge(from, to)
	return nil
}

func comparator(a, b interface{}) int {
	return utils.StringComparator(a.(Node).Id(), b.(Node).Id())
}

func edgeComparator(a, b interface{}) int {
	ae := a.(Edge)
	be := b.(Edge)
	if ae.From.Id() != be.From.Id() {
		return utils.StringComparator(ae.From.Id(), be.From.Id())
	}
	return utils.StringComparator(ae.To.Id(), be.To.Id())
}

// Sort returns a valid topological sorting of this topology's Nodes
func (t *Topology) Sort() ([]Node, error) {

	/*
		Implementation of Kahn's algorithm: Wikipedia pseudocode

			L ← Empty list that will contain the sorted elements
			S ← Set of all Nodes with no incoming edge
			while S is non-empty do
			    remove a node n from S
			    add n to tail of L
			    for each node m with an edge e from n to m do
			        remove edge e from the graph
			        if m has no other incoming Edges then
			            insert m into S
			if graph has Edges then
			    return error   (graph has at least one cycle)
			else
			    return L   (a topologically sorted order)
	*/
	L := make([]Node, 0, len(t.Nodes))
	Sq := binaryheap.NewWith(comparator)
	for _, x := range t.starts() {
		Sq.Push(x)
	}
	edges := t.Edges.Copy()

	i := 0
	for {
		if Sq.Empty() {
			break
		}
		var n Node

		ni, _ := Sq.Pop()
		n = ni.(Node)
		L = append(L, n)

		eq := binaryheap.NewWith(edgeComparator)
		for _, e := range edges[n.Id()] {
			eq.Push(e)
		}
		for !eq.Empty() {
			mi, _ := eq.Pop()
			m := mi.(Edge).To
			delete(edges[n.Id()], m.Id())
			if !edges.HasIncoming(m) {
				Sq.Push(m)
			}
		}
		i++

		// in case of bugs...
		if i > 2*t.bound() {
			return nil, ErrRuntimeExceeded
		}
	}
	if edges.Count() > 0 {
		return nil, &ErrCycleInTopology{OriginalEdges: t.Edges.Copy(), RemainingEdges: edges.Copy()}
	}
	return L, nil
}

func (t *Topology) bound() int {
	sum := len(t.Nodes)
	for _, ne := range t.Edges {
		sum += len(ne)
	}
	return sum
}

func (t *Topology) starts() []Node {
	ret := make([]Node, 0)
	for _, n := range t.Nodes {
		if t.Edges.HasIncoming(n) {
			continue
		}
		ret = append(ret, n)
	}
	return ret
}
