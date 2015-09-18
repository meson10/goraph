package goraph

import (
	"strings"
	"sync"
)

const NoteSeperator = "->"

type Taggable struct {
	tags  map[string]string
	notes []string
	sync.RWMutex
}

//Since go doesn't have constructors, initialize self.tags with a valid struct.
func (self *Taggable) SeedTag() {
	if self.tags == nil {
		self.tags = map[string]string{}
	}
}

func (self *Taggable) RemoveTag(key string) {
	self.Lock()
	defer self.Unlock()
	delete(self.tags, key)
}

//Set a Tag on the instance. You can only attach 10 of these.
func (self *Taggable) SetTag(key string, value string) Tagger {
	self.Lock()
	defer self.Unlock()
	self.tags[key] = value
	return self
}

//Fetch a previously saved Tag. Always returns a string, could be empty.
func (self *Taggable) HasTag(key string) bool {
	self.RLock()
	defer self.RUnlock()
	_, ok := self.tags[key]
	return ok
}

//Fetch a previously saved Tag. Always returns a string, could be empty.
func (self *Taggable) GetTag(key string) string {
	self.RLock()
	defer self.RUnlock()
	return self.tags[key]
}

//Get all the tags attached to this Instance.
func (self *Taggable) GetTags() map[string]string {
	self.RLock()
	defer self.RUnlock()
	return self.tags
}

func (self *Taggable) Notes() []string {
	self.RLock()
	defer self.RUnlock()
	return self.notes
}

func (self *Taggable) AddNote(checks ...string) Tagger {
	self.Lock()
	defer self.Unlock()
	self.notes = append(self.notes, strings.Join(checks, NoteSeperator))
	return self
}
