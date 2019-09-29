package toposort

import "testing"

func TestEdgesFromNode_Copy(t *testing.T) {
	a := &NodeTest{}
	b := &NodeTest{}
	e := make(EdgesFromNode)
	e["a"] = Edge{a, b}
	e2 := e.Copy()
	e2["a"] = Edge{b, a}
	if e["a"].From == b {
		t.Errorf("copy did not produce a copy of e: key a value changed in e")
	}
	if e["a"].To == a {
		t.Errorf("copy did not produce a copy of e: key b found in e")
	}
}

func TestEdges_AddEdge(t *testing.T) {
	a := &NodeTest{ID: "a"}
	b := &NodeTest{ID: "b"}
	edges := make(Edges)
	edges.AddEdge(a, b)
	if _, ok := edges["a"]; !ok {
		t.Fatalf("Edges.AddEdge(a, b) did not add edges[a]")
	}
	if _, ok := edges["a"]["b"]; !ok {
		t.Fatalf("Edges.AddEdge(a, b) did not add edges[a][b]")
	}
}
func TestEdges_Count(t *testing.T) {
	a := &NodeTest{ID: "a"}
	b := &NodeTest{ID: "b"}
	edges := make(Edges)
	edges.AddEdge(a, b)
	edges["c"] = nil
	if edges.Count() != 1 {
		t.Errorf("expected one edge")
	}

	if edges.HasIncoming(a) {
		t.Errorf("b should not have incoming edges")
	}
}
