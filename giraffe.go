package giraffe

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Graph struct {
	sync.Mutex
	Nodes map[uint64]*Node

	topNodeID uint64
}

func NewGraph(name string) (*Graph, error) {
	g := &Graph{
		Nodes: make(map[uint64]*Node),
	}
	// insert root node
	g.Nodes[0] = &Node{ID: 0}

	return g, nil
}

func (g *Graph) NodeCount() int {
	g.Lock()
	defer g.Unlock()
	return len(g.Nodes)
}

func (g *Graph) LastNodeID() uint64 {
	g.Lock()
	g.Unlock()
	return g.topNodeID
}

func (g *Graph) InsertNode() *Node {
	g.Lock()
	g.Unlock()

	n := &Node{
		ID: atomic.AddUint64(&g.topNodeID, 1),
	}

	g.Nodes[n.ID] = n

	return n
}

func (g *Graph) InsertDataNode(key string, value []byte) *Node {
	n := g.InsertNode()
	n.Key = key
	n.Value = value
	return n
}

func (g *Graph) FindRoots() []uint64 {
	g.Lock()
	defer g.Unlock()
	roots := make([]uint64, 0)
	for _, node := range g.Nodes {
		node.Lock()
		if len(node.sources) == 0 {
			roots = append(roots, node.ID)
		}
		node.Unlock()
	}
	return roots
}

func (g *Graph) FindNodeIDByKey(key string) (uint64, bool) {
	for id, node := range g.Nodes {
		if node.Key == key {
			return id, true
		}
	}

	return 0, false
}

func (g *Graph) FindNodeByKey(key string) (*Node, bool) {
	id, found := g.FindNodeIDByKey(key)
	if !found {
		return nil, false
	}

	return g.Nodes[id], true
}

func (g *Graph) ToVisJS() string {
	dataSet := ""
	edges := ""

	for id, node := range g.Nodes {
		dataSet += fmt.Sprintf(`{id: %d, label: '%s'},`, id, node.Key)
		for _, nID := range node.ListDestinations() {
			edges += fmt.Sprintf(`{from: %d, to: %d},`, id, nID)
		}
	}

	js := `
<html>
<head>
    <script type="text/javascript" src="http://visjs.org/dist/vis.js"></script>
    <link href="http://visjs.org/dist/vis.css" rel="stylesheet" type="text/css" />

    <style type="text/css">
        #mynetwork {
            width: 600px;
            height: 400px;
            border: 1px solid lightgray;
        }
    </style>
</head>
<body>
<div id="mynetwork"></div>

<script type="text/javascript">
    // create an array with nodes
    var nodes = new vis.DataSet([
        ` + dataSet + `
    ]);

    // create an array with edges
    var edges = new vis.DataSet([
        ` + edges + `
    ]);

    // create a network
    var container = document.getElementById('mynetwork');

    // provide the data in the vis format
    var data = {
        nodes: nodes,
        edges: edges
    };
    var options = {};

    // initialize your network!
    var network = new vis.Network(container, data, options);
</script>
</body>
</html>
    `

	return js
}
