package giraffe

import "testing"

func TestRootNodeCreate(t *testing.T) {
	g, err := NewGraph("testGraph")
	if err != nil {
		t.Fatalf("unable to create graph, error: %v", err)
	}

	if g.NodeCount() != 1 {
		t.Error("expected a root node to be created")
	}
	if g.LastNodeID() != 0 {
		t.Errorf("got id %d, want id %d for root node", g.LastNodeID(), 0)
	}
}

func TestAddNewNodesWithDestinationsAndSources(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()
	n4 := g.InsertNode()

	root := g.Nodes[0]
	root.AddRelationship(n1)
	root.AddRelationship(n2)
	n2.AddRelationship(n3)
	n4.AddRelationship(n3)

	destinations := root.ListDestinations()

	if len(destinations) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(destinations), 2)
	}

	if want := []uint64{n1.ID, n2.ID}; !ContainsAll(destinations, want) {
		t.Errorf("actual destinations %v not expected %v", destinations, want)
	}

	sources := n3.ListSources()

	if len(sources) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(sources), 2)
	}

	if want := []uint64{n2.ID, n4.ID}; !ContainsAll(sources, want) {
		t.Errorf("actual sources %v not expected %v", sources, want)
	}

}

func TestFindRoots(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1 := g.InsertNode()
	n2 := g.InsertNode()
	n3 := g.InsertNode()
	n4 := g.InsertNode()

	root := g.Nodes[0]
	root.AddRelationship(n1)
	root.AddRelationship(n2)
	n2.AddRelationship(n3)

	roots := g.FindRoots()
	expectedRoots := []uint64{root.ID, n4.ID}

	if len(roots) != len(expectedRoots) {
		t.Errorf("got %d roots, wat %d roots", len(roots), len(expectedRoots))
	}

	if !ContainsAll(roots, expectedRoots) {
		t.Errorf("actual roots %v not expected %v", roots, expectedRoots)
	}
}

func TestAddNewDataNodesAndFindNode(t *testing.T) {
	g, _ := NewGraph("testGraph")
	n1, _ := g.InsertDataNode("key1", []byte("value1"))
	n2, _ := g.InsertDataNode("key2", []byte("value2"))
	n3, _ := g.InsertDataNode("key3", []byte("value3"))

	root := g.Nodes[0]
	root.AddRelationship(n1)
	root.AddRelationship(n2)
	n2.AddRelationship(n3)

	destinations := root.ListDestinations()

	if len(destinations) != 2 {
		t.Errorf("got %d, want %d destination nodes", len(destinations), 2)
	}

	if want := []uint64{n1.ID, n2.ID}; !ContainsAll(destinations, want) {
		t.Errorf("actual destinations %v not expected %v", destinations, want)
	}

	foundNode, ok := g.FindNodeByKey("key3")
	if !ok {
		t.Fatal("unable to find node")
	}

	if string(foundNode.Value) != "value3" {
		t.Errorf("got `%s`, want `value3` data", foundNode.Value)
	}
}

func TestDuplicateKeyError(t *testing.T) {
	g, _ := NewConstraintGraph("testGraph", false, true)

	_, err := g.InsertDataNode("key", []byte("data"))
	if err != nil {
		t.Errorf("expected no error, got `%s`", err)
	}

	_, err = g.InsertDataNode("key", []byte("oops, same key"))
	if err == nil {
		t.Error("no error when there should have been")
	}

	if err.Error() != ErrKeyExists {
		t.Errorf("unexpected error message. got `%s`, want `%s`", err.Error(), ErrKeyExists)
	}
}

func ContainsAll(superset, subset []uint64) bool {
	for _, sub := range subset {
		found := false
		for _, super := range superset {
			if sub == super {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
