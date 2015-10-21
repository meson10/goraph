package goraph

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	dirty     = "*"
	blank     = ""
	comma     = ","
	semiColon = ":"
)

var MachineReadable bool

func FieldSeperator() string {
	if MachineReadable {
		return comma
	} else {
		return semiColon
	}
}

type Node struct {
	id        string
	node_type string
	parent    GraphNode
	children  map[string]GraphNode
	Taggable
}

func (self *Node) SetType(ntype string) {
	self.Lock()
	defer self.Unlock()
	self.node_type = ntype
}

func (self *Node) Delete() {
	parent := self.Parent()
	if parent != nil {
		parent.RemoveChild(self.Id())
	}
}

func (self *Node) RemoveChild(id string) {
	self.Lock()
	defer self.Unlock()
	delete(self.children, id)
}

func (self *Node) Children() map[string]GraphNode {
	self.RLock()
	defer self.RUnlock()
	return self.children
}

func (self *Node) Find(id string) GraphNode {
	self.RLock()
	defer self.RUnlock()
	return self.children[id]
}

func (self *Node) Parent() GraphNode {
	return self.parent
}

func (self *Node) SetId(ident string) GraphNode {
	self.id = ident
	return self
}

func (self *Node) Id() string {
	return self.id
}

func (self *Node) InitChildren() {
	if self.children == nil {
		self.children = map[string]GraphNode{}
	}
}

func (self *Node) IsDirty() bool {
	return self.HasTag(dirty)
}

func (self *Node) RemoveDirty() GraphNode {
	self.RemoveTag(dirty)
	return self
}

func (self *Node) SetDirty() GraphNode {
	self.SetTag(dirty, blank)
	return self
}

func (self *Node) AddChild(n GraphNode) {
	n.SetParent(self)

	self.Lock()
	defer self.Unlock()
	self.children[n.Id()] = n
}

func (self *Node) SetParent(p GraphNode) {
	self.Lock()
	defer self.Unlock()
	self.parent = p
}

func (self *Node) Append(parent GraphNode, child GraphNode) {
	parent.AddChild(child)
}

func (self *Node) DisplayName() string {
	return self.node_type + FieldSeperator() + self.Id()
}

func (self *Node) MachineOutput(prefix string, w io.Writer) {
	var log string

	if self.Parent() != nil {
		if prefix != blank {
			log = fmt.Sprintf("%v,%v", prefix, self.DisplayName())
		} else {
			log = self.DisplayName()
		}

		arr := []string{log}

		for k, v := range self.GetTags() {
			if k == "*" {
				continue
			}
			arr = append(arr, k)
			arr = append(arr, v)
		}

		if self.IsDirty() {
			arr = append(arr, "dirty")
			arr = append(arr, "true")
		}

		w.Write([]byte(strings.Join(arr, comma)))
		w.Write([]byte("\n"))
	}

	for _, x := range self.Children() {
		x.MachineOutput(log, w)
	}
}

func (self *Node) Output(indent int, last bool, w io.Writer) {
	sep := blank
	if indent > 0 {
		sep += strings.Repeat(TAB, indent-1)
		if last == true {
			sep += LASTLEAF
		} else {
			sep += LEAF
		}
		sep += SEP
	}

	log := fmt.Sprintf("%v%v", sep, self.Id())

	//Display a sweet little star next to the node, if dirty.
	is_dirty := self.IsDirty()

	if is_dirty == true {
		log += "[*]"
	}

	//Display Tags associated with this node.
	tags := self.GetTags()
	if len(tags) > 0 {
		tagstr := blank
		for k, v := range tags {
			if k == "*" {
				continue
			}
			tagstr += fmt.Sprintf(" %v=>%v", k, v)
		}
		log += tagstr
	}

	if is_dirty {
		color.Set(color.FgRed)
	}

	w.Write([]byte(log))
	w.Write([]byte("\n"))

	if indent >= 0 {
		indent++
	}

	//Display Notes associated with this node.
	notes := self.Notes()
	if len(notes) > 0 {
		reasons := NewNode("Notes")
		self.AddChild(reasons)
		for _, r := range notes {
			reasons.AddChild(NewNode(r))
		}
	}

	// Keep doing this for all children.
	children := self.Children()
	length := len(children)

	//Need to maintain a separate index, since iteration through maps
	//do not yield an inbuilt iteratoion index.
	ix := 0
	for _, x := range children {
		ix++
		if ix == length {
			last = true
		} else {
			last = false
		}
		x.Output(indent, last, w)
	}

	color.Unset()
}

func (self *Node) Walk() {
	if MachineReadable == true {
		self.MachineOutput(blank, os.Stdout)
	} else {
		self.Output(0, false, os.Stdout)
	}
}
