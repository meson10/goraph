package goraph

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	NODES    = 10
	CHILDREN = 10
)

var g Graph

func TestGraph(t *testing.T) {
	g = MakeGraph()
}

func nodeIdent(ix int) string {
	return string(ix)
}

func childIdent(ix int) string {
	return fmt.Sprintf("child%v", ix)
}

func TestGraphAddNode(t *testing.T) {

	var wg sync.WaitGroup

	for i := 0; i < NODES; i++ {
		wg.Add(1)

		go func(ix int) {
			defer wg.Done()
			g.AddChild(NewNode(nodeIdent(ix)))
		}(i)
	}

	wg.Wait()

	children := g.Children()
	familySize := len(children)

	if familySize != NODES {
		t.Error(fmt.Sprintf("Graph should have %v children. Found %v", NODES, familySize))
	}
}

func TestGraphAddChildren(t *testing.T) {
	var mg sync.WaitGroup

	for _, c := range g.Children() {
		mg.Add(1)

		go func(c GraphNode) {
			defer mg.Done()

			for i := 0; i < CHILDREN; i++ {
				c.AddChild(NewNode(fmt.Sprintf("child%v", i)))
			}

		}(c)
	}

	mg.Wait()

	for id, c := range g.Children() {
		if len(c.Children()) != CHILDREN {
			t.Error(fmt.Sprintf("Child %v should have %v children.", id, CHILDREN))
		}
	}

}

func TestGraphFindChild(t *testing.T) {
	id := rand.Intn(NODES)
	if node := g.Find(nodeIdent(id)); node == nil {
		t.Error(fmt.Sprintf("Cannot find node %v", id))
		return
	}

	invalid_id := id + 1000
	if node := g.Find(string(invalid_id)); node != nil {
		t.Error(fmt.Sprintf("Should not find node %v", invalid_id))
		return
	}
}

func TestGraphFindNodeChild(t *testing.T) {
	id := rand.Intn(NODES)
	node := g.Find(nodeIdent(id))

	childId := rand.Intn(CHILDREN)
	child := node.Find(childIdent(childId))
	if child == nil {
		t.Error(fmt.Sprintf("Cannot find child %v for node %v", childId, id))
		return
	}
}

func TestGraphNodeParent(t *testing.T) {
	id := rand.Intn(NODES)
	node := g.Find(nodeIdent(id))

	childId := rand.Intn(CHILDREN)
	child := node.Find(childIdent(childId))

	if child.Parent() != node {
		t.Error(fmt.Sprintf("%v 's parent != %v", childId, id))
	}
}

func TestGraphNodeDirty(t *testing.T) {
	id := rand.Intn(NODES)
	node := g.Find(nodeIdent(id))
	node.SetDirty()

	if node.IsDirty() != true {
		t.Error(fmt.Sprintf("%v should have been dirty", id))
	}

	n := g.Find(nodeIdent(id))
	if n.IsDirty() != true {
		t.Error(fmt.Sprintf("%v should have been dirty", id))
	}

	n.RemoveDirty()
	if n.IsDirty() == true {
		t.Error(fmt.Sprintf("%v should not been dirty anymore", id))
	}
}

func TestGraphNodeDelete(t *testing.T) {
	id := rand.Intn(NODES)
	node := g.Find(nodeIdent(id))
	node.Delete()

	n := g.Find(nodeIdent(id))
	if n != nil {
		t.Error(fmt.Sprintf("%v should have been deleted", id))
	}
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	status := m.Run()
	os.Exit(status)
}
