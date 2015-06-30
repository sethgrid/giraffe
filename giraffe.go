package giraffe

import (
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
