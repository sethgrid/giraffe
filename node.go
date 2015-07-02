package giraffe

import (
	"errors"
	"sync"
)

const (
	// ErrCircular returns when adding a relationship would result in a circular relationship
	ErrCircular = "circular relationship"
)

// Node is a key value pair that has directional relationships with other Nodes
type Node struct {
	ID    uint64
	Key   string
	Value []byte

	circularRelationship bool

	sync.RWMutex
	destinations []*Node
	sources      []*Node
}

// AddRelationship adds a newNode as a destination of this node
func (n *Node) AddRelationship(newNode *Node) error {
	n.Lock()
	defer n.Unlock()

	if !n.circularRelationship && newNode.DepthFirstSearch(n) {
		return errors.New(ErrCircular)
	}

	n.destinations = append(n.destinations, newNode)
	newNode.addSource(n)

	return nil
}

// addSource provides a way to more easily traverse the graph.
// addSource is not exposed to prevent circular logic as nodes are added
// to both destination and source lists
func (n *Node) addSource(newNode *Node) {
	n.Lock()
	defer n.Unlock()

	n.sources = append(n.sources, newNode)
}

// ListDestinations lists all nodes that this node points towards
func (n *Node) ListDestinations() []*Node {
	n.Lock()
	defer n.Unlock()

	return n.destinations
}

// ListSources lists all nodes that point to this node
func (n *Node) ListSources() []*Node {
	n.Lock()
	defer n.Unlock()

	return n.sources
}

// DepthFirstSearch traverses the graph starting at this node to find the otherNode
func (n *Node) DepthFirstSearch(otherNode *Node) bool {
	for _, node := range n.destinations {
		if node.ID == otherNode.ID || node.DepthFirstSearch(otherNode) {
			return true
		}
	}
	return false
}

// BreadthFirstSearch traverses the graph starting at this node to find the otherNode
func (n *Node) BreadthFirstSearch(otherNode *Node) bool {
	// initialize the bfs queue
	return bfs(otherNode, n.destinations)
}

// bfs maintains the search queue and is the logic for BreadthFirstSearch()
func bfs(otherNode *Node, queue []*Node) bool {
	if len(queue) == 0 {
		return false
	}

	nextQueue := make([]*Node, 0)
	for _, node := range queue {
		if node.ID == otherNode.ID {
			return true
		}
		nextQueue = append(nextQueue, node.destinations...)
	}

	if bfs(otherNode, nextQueue) {
		return true
	}
	return false
}

// extractIDs grabs all the ids for the nodes. useful for debugging and in tests.
func extractIDs(nodes []*Node) []uint64 {
	IDs := make([]uint64, len(nodes))
	for i, node := range nodes {
		IDs[i] = node.ID
	}
	return IDs
}
