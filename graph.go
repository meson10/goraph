package goraph

import "strings"

var TAB = "│" + strings.Repeat(" ", 4)

const SEP = "─── "
const LEAF = "├"
const LASTLEAF = "└"

func NewNode(id string) *Node {
	if id == "" {
		id = "unknown"
	}

	x := &Node{id: id}
	x.SeedTag()
	x.InitChildren()
	return x
}

func MakeGraph() Graph {
	x := NewNode("root")
	x.parent = nil
	return x
}
