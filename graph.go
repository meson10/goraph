package goraph

import "strings"

var TAB = "│" + strings.Repeat(" ", 4)

const SEP = "─── "
const LEAF = "├"
const LASTLEAF = "└"
const UNKNOWN = "unknown"

func NewNode(id string) *Node {
	if id == "" {
		id = UNKNOWN
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
