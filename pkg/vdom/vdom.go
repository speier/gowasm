package vdom

import (
	"golang.org/x/net/html/atom"
)

type NodeType uint32

const (
	TextNode NodeType = iota
	ElementNode
)

type Attribute struct {
	Key, Val string
}

type VNode struct {
	Type     NodeType
	Data     string
	Attrs    []*Attribute
	Children []*VNode
}

func H(tagName string, attrs []*Attribute, children ...*VNode) *VNode {
	a := atom.Lookup([]byte(tagName))
	if a == 0 {
		return &VNode{Type: TextNode, Data: tagName}
	}
	return &VNode{Type: ElementNode, Data: tagName, Attrs: attrs, Children: children}
}

func HText(data string) *VNode {
	return &VNode{Data: data}
}
