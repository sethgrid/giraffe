package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/sethgrid/giraffe"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	// make a new graph
	g, _ := giraffe.NewGraph("example")

	// we automatically have access to the root node
	root := g.Nodes[0]

	// let's make a bunch of new nodes
	nodeCount := 50
	nodes := make([]*giraffe.Node, nodeCount)
	nodes[0] = root
	for i := 1; i < nodeCount; i++ {
		nodes[i] = g.InsertDataNode(fmt.Sprintf("node %d", i), []byte{})
	}

	// let's make random relationships between the nodes
	for i := 0; i < nodeCount*1; i++ {
		nodeA := int(rand.Intn(nodeCount - 1))
		nodeB := int(rand.Intn(nodeCount - 1))
		if nodeA == nodeB {
			nodeA = 0
		}

		nodes[nodeA].AddRelationship(nodes[nodeB])
	}

	// and we can now visualize the relationships
	fmt.Println(g.ToVisJS())

	// while we started with one root node, we can have many roots
	log.Println("root nodes: ", g.FindRoots())
}
