package giraffe

import "sync"

// Node is a key value pair that has directional relationships with other Nodes
type Node struct {
	ID    uint64
	Key   string
	Value []byte

	sync.RWMutex
	destinations []uint64
	sources      []uint64
}

// AddRelationship adds a newNode as a destination of this node
func (n *Node) AddRelationship(newNode *Node) {
	n.Lock()
	defer n.Unlock()

	n.destinations = append(n.destinations, newNode.ID)
	newNode.addSource(n)
}

// addSource provides a way to more easily traverse the graph.
// addSource is not exposed to prevent circular logic as nodes are added
// to both destination and source lists
func (n *Node) addSource(newNode *Node) {
	n.Lock()
	defer n.Unlock()

	n.sources = append(n.sources, newNode.ID)
}

// ListDestinations lists all nodes that this node points towards
func (n *Node) ListDestinations() []uint64 {
	n.Lock()
	defer n.Unlock()
	c := make([]uint64, len(n.destinations))
	copy(c, n.destinations)
	return c
}

// ListSources lists all nodes that point to this node
func (n *Node) ListSources() []uint64 {
	n.Lock()
	defer n.Unlock()
	c := make([]uint64, len(n.sources))
	copy(c, n.sources)
	return c
}
