package giraffe

import "sync"

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

func (n *Node) addSource(newNode *Node) {
	n.Lock()
	defer n.Unlock()

	n.sources = append(n.sources, newNode.ID)
}

func (n *Node) ListDestinations() []uint64 {
	n.Lock()
	defer n.Unlock()
	c := make([]uint64, len(n.destinations))
	copy(c, n.destinations)
	return c
}
