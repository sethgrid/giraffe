package giraffe

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	ErrKeyExists = "key exists"
)

// Graph is a data structure composed of Nodes that have directional relationships to one another
type Graph struct {
	Name string

	duplicateKeys        bool
	circularRelationship bool

	sync.Mutex
	Nodes map[uint64]*Node
	keys  map[string]bool

	topNodeID uint64
}

// NewGraph creates a graph with default properties
func NewGraph(name string) (*Graph, error) {
	return NewConstraintGraph(name, true, true)
}

// NewConstraintGraph allows you to control if duplicate keys or circular relationships can exist
// TODO: implement circular relationship check
func NewConstraintGraph(name string, duplicateKeys, circularRelationship bool) (*Graph, error) {
	g := &Graph{
		Name:                 name,
		Nodes:                make(map[uint64]*Node),
		keys:                 make(map[string]bool),
		duplicateKeys:        duplicateKeys,
		circularRelationship: circularRelationship,
	}
	// insert root node
	g.Nodes[0] = &Node{ID: 0}

	return g, nil
}

// NodeCount returns the number of nodes in the graph
func (g *Graph) NodeCount() int {
	g.Lock()
	defer g.Unlock()
	return len(g.Nodes)
}

// LastNodeID returns the id of the last inserted node
func (g *Graph) LastNodeID() uint64 {
	g.Lock()
	g.Unlock()
	return g.topNodeID
}

// InsertNode inserts an empty default node into the graph. See InsertDataNode().
func (g *Graph) InsertNode() *Node {
	g.Lock()
	g.Unlock()

	n := &Node{
		ID: atomic.AddUint64(&g.topNodeID, 1),
	}

	g.Nodes[n.ID] = n

	return n
}

// InsertDataNode is an alternate constructor to InsertNode() allowing you to pass in a key and value
func (g *Graph) InsertDataNode(key string, value []byte) (*Node, error) {
	if !g.duplicateKeys {
		if _, ok := g.keys[key]; ok {
			return nil, errors.New(ErrKeyExists)
		}
		g.keys[key] = true
	}

	n := g.InsertNode()
	n.Key = key
	n.Value = value
	return n, nil
}

// FindRoots finds all nodes that do not have a source nodes below them
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

// FindNodeIDByKey returns the first node's ID with a matching key
func (g *Graph) FindNodeIDByKey(key string) (uint64, bool) {
	for id, node := range g.Nodes {
		if node.Key == key {
			return id, true
		}
	}

	return 0, false
}

// FindNodeByKey returns the the first node with a matching key
func (g *Graph) FindNodeByKey(key string) (*Node, bool) {
	id, found := g.FindNodeIDByKey(key)
	if !found {
		return nil, false
	}

	return g.Nodes[id], true
}

// ToVisJS generates a simple HTML/Javascript view of the data
// See http://visjs.org/ for information and styles
func (g *Graph) ToVisJS() string {
	dataSet := ""
	edges := ""

	for id, node := range g.Nodes {
		dataSet += fmt.Sprintf(`{id: %d, label: 'node %d'},`, id, id)
		for _, nID := range node.ListDestinations() {
			edges += fmt.Sprintf(`{from: %d, to: %d, arrows:'middle',},`, id, nID)
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
