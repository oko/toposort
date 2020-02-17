package main

import (
	"fmt"
	"github.com/oko/toposort"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	topofile = pflag.StringP("topo", "T", "", "topology file")
)

type SimpleNode struct{ ID string }

func (sn *SimpleNode) Id() string { return sn.ID }

func parseTopo(topo string) (*toposort.Topology, []toposort.Edge, error) {
	t := toposort.NewTopology()
	e := make([]toposort.Edge, 0)
	for _, line := range strings.Split(topo, "\n") {
		trim := strings.TrimSpace(line)
		if len(trim) == 0 {
			continue
		}
		split := strings.Split(trim, "->")
		if len(split) != 2 {
			return nil, nil, fmt.Errorf("invalid line %s", line)
		}
		from := &SimpleNode{ID: strings.TrimSpace(split[0])}
		to := &SimpleNode{ID: strings.TrimSpace(split[1])}
		_ = t.AddNode(from)
		_ = t.AddNode(to)
		_ = t.AddEdge(from, to)
		e = append(e, toposort.Edge{From: from, To: to})
	}
	return t, e, nil
}

func main() {
	pflag.Parse()
	var input io.Reader
	if *topofile == "-" {
		input = os.Stdin
		*topofile = "STDIN"
	} else {
		f, err := os.Open(*topofile)
		if err != nil {
			log.Fatalf("failed to open %s: %s", *topofile, err)
		}
		defer f.Close()
		input = f
	}
	data, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("error reading from %s: %s", *topofile, err)
	}
	topo, edges, err := parseTopo(string(data))
	if err != nil {
		log.Fatalf("failed to parse topology: %s", err)
	}
	sort, err := topo.Sort()
	if err != nil {
		if cy, ok := err.(*toposort.ErrCycleInTopology); ok {
			log.Printf("finding candidate cycles...")
			findCycles(cy)
		}
		log.Fatalf("failed to sort topo: %s", err)
	}

	printBaseTopo(sort, edges)

}

func findCycles(e *toposort.ErrCycleInTopology) {
	ex := make(map[string]bool)
	for _, edges := range e.OriginalEdges {
		for _, e := range edges {
			es := fmt.Sprintf("%s->%s", e.From.Id(), e.To.Id())
			if _, ok := ex[es]; ok {
				log.Printf("possible dupe? %s", es)
			}
			ex[es] = true
		}
	}
	for _, edges := range e.RemainingEdges {
		for _, e := range edges {
			es := fmt.Sprintf("%s->%s", e.To.Id(), e.From.Id())
			if _, ok := ex[es]; ok {
				log.Printf("possible dupe? %s", es)
			}
			ex[es] = true
		}
	}
}

// this should really be a template.........
func printBaseTopo(sort []toposort.Node, edges []toposort.Edge) {
	fmt.Println("digraph G {")
	fmt.Println("  subgraph cluster_orig {")
	for _, e := range edges {
		fmt.Printf("    %s -> %s\n", e.From.Id(), e.To.Id())
	}
	fmt.Println("  }")
	fmt.Println("  subgraph cluster_sort {")
	fmt.Println("    rankdir=LR;")
	var last toposort.Node
	fmt.Println("    edge[style=invis];")
	for _, n := range sort {
		fmt.Printf("    %s_sort [label=\"%s\"]\n", n.Id(), n.Id())
		if last != nil {
			fmt.Printf("    %s_sort -> %s_sort\n", last.Id(), n.Id())
		}
		last = n
	}
	fmt.Println("    edge[style=solid, constraint=false];")
	for _, e := range edges {
		fmt.Printf("    %s_sort -> %s_sort\n", e.From.Id(), e.To.Id())
	}
	fmt.Println("")
	fmt.Println("  }")
	fmt.Println("}")
}
