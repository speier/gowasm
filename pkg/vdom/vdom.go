package vdom

import (
	"fmt"
	"strings"

	"golang.org/x/net/html/atom"
)

type NodeType uint32

const (
	TextNode NodeType = iota
	ElementNode
)

type Props = map[string]string

type Events = map[string]func()

type Attrs struct {
	Props  *Props
	Events *Events
}

type VNode struct {
	Type     NodeType
	TagName  string
	Attrs    *Attrs
	Children []*VNode
}

// hyperscript-style API: h(tagName, attrs, children)
func H(tagName string, attrs *Attrs, children ...*VNode) *VNode {
	a := atom.Lookup([]byte(strings.ToLower(tagName)))
	if a == 0 {
		return Text(tagName)
	}
	return &VNode{Type: ElementNode, TagName: strings.ToLower(tagName), Attrs: attrs, Children: children}
}

func Text(text string) *VNode {
	return &VNode{Type: TextNode, TagName: text}
}

// TODO
func (vnode *VNode) HashCode() string {
	if vnode.Type == TextNode {
		return vnode.TagName
	}
	if vnode.Attrs != nil && vnode.Attrs.Props != nil {
		return fmt.Sprintf("%s/%v", vnode.TagName, *vnode.Attrs.Props)
	}
	return fmt.Sprintf("%s/%v", vnode.TagName, Attrs{})
}
