package main

import (
	"fmt"

	"github.com/sethgrid/giraffe"
)

func main() {
	g, _ := giraffe.NewGraph("example")

	root := g.Nodes[0]
	n1 := g.InsertDataNode("node 1", []byte("data"))
	n2 := g.InsertDataNode("node 2", []byte("data"))
	n3 := g.InsertDataNode("node 3", []byte("data"))
	n4 := g.InsertDataNode("node 4", []byte("data"))
	n5 := g.InsertDataNode("node 5", []byte("data"))

	root.AddRelationship(n1)
	root.AddRelationship(n2)
	root.AddRelationship(n3)
	n3.AddRelationship(n4)
	n4.AddRelationship(n5)
	n5.AddRelationship(n3)

	fmt.Println(g.ToVisJS())
}
