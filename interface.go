package goraph

import "io"

type GraphNode interface {
	Graph
	IsDirty() bool
	RemoveDirty() GraphNode
	SetDirty() GraphNode
	Output(int, bool, io.Writer)
	SetParent(GraphNode)
	Parent() GraphNode
	Delete()
	RemoveChild(string)
}

type Graph interface {
	Tagger
	SetId(string) GraphNode
	Id() string
	InitChildren()
	Children() map[string]GraphNode
	AddChild(GraphNode)
	Walk()
	Find(string) GraphNode
}

type Tagger interface {
	SeedTag()
	Notes() []string
	AddNote(...string) Tagger
	HasTag(string) bool
	RemoveTag(string)
	SetTag(string, string) Tagger
	GetTag(string) string
	GetTags() map[string]string
}
