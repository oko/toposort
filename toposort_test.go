package toposort

import (
	"encoding/json"
	"strings"
	"testing"
)

func parseTopology(topo string) *Topology {
	t := NewTopology()
	lines := strings.Split(topo, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		parts := strings.Split(l, "->")
		if len(parts) != 2 {
			panic("need format X -> Y for topo")
		}
		from := &NodeTest{
			strings.TrimSpace(parts[0]),
		}
		to := &NodeTest{
			strings.TrimSpace(parts[1]),
		}
		_ = t.AddNode(from)
		_ = t.AddNode(to)
		_ = t.AddEdge(from, to)
	}
	return t
}

type NodeTest struct {
	ID string
}

func (n *NodeTest) Id() string {
	return n.ID
}

var _ Node = &NodeTest{}

func TestNewTopology(t *testing.T) {
	nt := NewTopology()
	if nt == nil {
		t.Fatalf("topology returned was nil")
	}
	if nt.nodes == nil {
		t.Fatalf("did not initialize Topology.nodes in NewTopology")
	}
	if nt.edges == nil {
		t.Fatalf("did not initialize Topology.edges in NewTopology")
	}
}
func TestTopology_AddNode(t *testing.T) {
	nt := NewTopology()
	if err := nt.AddNode(&NodeTest{ID: "a"}); err != nil {
		t.Fatalf("failed to add node to topology: %s", err)
	}
	if _, ok := nt.nodes["a"]; !ok {
		t.Fatalf("Topology.AddNode didn't actually add node")
	}
	if len(nt.nodes) != 1 {
		t.Fatalf("expected a node to be added to topology")
	}
	if nt.nodes["a"].Id() != "a" {
		t.Fatalf("somehow ended up with a different node added?")
	}
}

func TestTopology_AddEdge(t *testing.T) {
	nt := NewTopology()
	a := &NodeTest{ID: "a"}
	b := &NodeTest{ID: "b"}
	if err := nt.AddEdge(a, b); err == nil {
		t.Fatalf("should have errored with edges added after nodes")
	}
}

func TestTopology_Sort(t *testing.T) {
	ns := make(Nodes)
	for _, x := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"} {
		ns[x] = &NodeTest{ID: x}
	}
	nt := NewTopology()
	for _, n := range ns {
		if err := nt.AddNode(n); err != nil {
			t.Fatalf("error adding node %s", n.Id())
		}
	}
	_ = nt.AddEdge(ns["a"], ns["b"])
	_ = nt.AddEdge(ns["a"], ns["c"])
	_ = nt.AddEdge(ns["b"], ns["d"])
	_ = nt.AddEdge(ns["b"], ns["e"])
	_ = nt.AddEdge(ns["c"], ns["f"])
	_ = nt.AddEdge(ns["d"], ns["g"])
	_ = nt.AddEdge(ns["e"], ns["h"])
	_ = nt.AddEdge(ns["f"], ns["h"])
	_ = nt.AddEdge(ns["g"], ns["i"])
	_ = nt.AddEdge(ns["h"], ns["i"])

	if sort, err := nt.Sort(); err != nil {
		t.Errorf("error sorting %s", err)
	} else {
		js, err := json.MarshalIndent(sort, "", "  ")
		if err == nil {
			t.Logf("SORT: %s", string(js))
		}
	}

	_ = nt.AddEdge(ns["i"], ns["g"])

	if _, err := nt.Sort(); err == nil {
		t.Errorf("did not get expected sorting error with cycle")
	}
}

func TestTopology_SortWorst(t *testing.T) {
	ns := make(Nodes)
	for _, x := range []string{"d", "c", "b", "a"} {
		ns[x] = &NodeTest{ID: x}
	}
	nt := NewTopology()
	for _, n := range ns {
		if err := nt.AddNode(n); err != nil {
			t.Fatalf("error adding node %s", n.Id())
		}
	}
	_ = nt.AddEdge(ns["c"], ns["d"])
	_ = nt.AddEdge(ns["b"], ns["d"])
	_ = nt.AddEdge(ns["a"], ns["c"])
	_ = nt.AddEdge(ns["a"], ns["b"])

	if sort, err := nt.Sort(); err != nil {
		t.Errorf("error sorting %s", err)
	} else {
		js, err := json.MarshalIndent(sort, "", "  ")
		if err == nil {
			t.Logf("SORT: %s", string(js))
		}
	}
}

var testCases = map[string][]string{
	`
a -> b
a -> c
`: []string{"a", "b", "c"},
	`
a -> b
a -> c
c -> d
`: []string{"a", "b", "c", "d"},
	`
a -> b
b -> a
`: nil,
	`
c -> d
b -> d
a -> c
a -> b
`: []string{"a", "b", "c", "d"},
}

func sort2string(ns []Node) string {
	ss := make([]string, 0, len(ns))
	for _, n := range ns {
		ss = append(ss, string(n.Id()))
	}
	return strings.Join(ss, " ")
}

func doTestCaseBench(t *testing.B, topo string, result []string) {
	tp := parseTopology(topo)
	sort, err := tp.Sort()
	if result == nil && err == nil {
		t.Errorf("got valid sort from unsortable test case: %s => %#v got %s", topo, result, sort2string(sort))
	}
	if result == nil && err != nil {
		return
	}
	if err != nil {
		t.Errorf("error sorting: %s", err)
		return
	}
	if len(result) != len(sort) {
		t.Errorf("mismatched length from sort: %d != %d", len(result), len(sort))
	}
}
func doTestCase(t *testing.T, topo string, result []string) {
	tp := parseTopology(topo)
	sort, err := tp.Sort()
	if result == nil && err == nil {
		t.Errorf("got valid sort from unsortable test case: %s => %#v got %s", topo, result, sort2string(sort))
	}
	if result == nil && err != nil {
		return
	}
	if err != nil {
		t.Errorf("error sorting: %s", err)
		return
	}
	if len(result) != len(sort) {
		t.Errorf("mismatched length from sort: %d != %d", len(result), len(sort))
	}
}
func TestTopology_SortVarious(t *testing.T) {
	for topo, res := range testCases {
		doTestCase(t, topo, res)
	}
}

func BenchmarkTopology_SortVarious(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for topo, res := range testCases {
			doTestCaseBench(b, topo, res)
		}
	}
}
func BenchmarkTopology_Sort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make(Nodes)
		for _, x := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"} {
			ns[x] = &NodeTest{ID: x}
		}
		nt := NewTopology()
		for _, n := range ns {
			_ = nt.AddNode(n)
		}
		_ = nt.AddEdge(ns["a"], ns["b"])
		_ = nt.AddEdge(ns["a"], ns["c"])
		_ = nt.AddEdge(ns["b"], ns["d"])
		_ = nt.AddEdge(ns["b"], ns["e"])
		_ = nt.AddEdge(ns["c"], ns["f"])
		_ = nt.AddEdge(ns["d"], ns["g"])
		_ = nt.AddEdge(ns["e"], ns["h"])
		_ = nt.AddEdge(ns["f"], ns["h"])
		_ = nt.AddEdge(ns["g"], ns["i"])
		_ = nt.AddEdge(ns["h"], ns["i"])

		if _, err := nt.Sort(); err != nil {
			b.Errorf("error sorting %s", err)
		}
	}
}
