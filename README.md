# `toposort`: a topological sorting library

This package implements a basic topological sorting library using Khan's algorithm.

## Topological Sorting

> In computer science, a topological sort or topological ordering of a directed graph is a linear ordering of its vertices such that for every directed edge uv from vertex u to vertex v, u comes before v in the ordering. &mdash;[Wikipedia](https://en.wikipedia.org/wiki/Topological_sorting)

![Image of graph and its topological sorting](https://raw.githubusercontent.com/oko/toposort/master/example/bigger.topo.png)

## Example

Your node types must implement the `toposort.Node` interface, which is very simple:

```
type Node interface {
	Id() string
}
```

`Id()` should return an identifier that uniquely identifies a given node in the graph.

A simple program with a topological sort:

```
package main

import "github.com/oko/toposort"

type TopoNode struct {
	ID string
}
func (t *TopoNode) Id() string {
	return t.ID
}
var _ toposort.Node = &TopoNode{}

func main() {
	t := toposort.NewTopology()
	a := &TopoNode{ID: "a"}
	b := &TopoNode{ID: "b"}
	t.AddNode(a)
	t.AddNode(b)
	t.AddEdge(a, b)
	sorted, err := t.Sort()
	if err != nil {
		panic("error sorting")
	}
	for _, s := range sorted {
		print(s.Id() + "\n")
	}
}
```

### Example Sorts

If you install Graphviz's `dot` binary, you can use `scripts/example.sh` to generate images like the one used at the top of this readme:

```
$ cd toposort
$ sudo apt install graphviz
$ scripts/example.sh example/bigger.topo
```

This will generate and open (using `xdg-open`) a rendering of the `example/bigger.topo` topology.
