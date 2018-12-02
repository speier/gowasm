package vdom

import (
	"fmt"

	"golang.org/x/net/html/atom"
)

type NodeType uint32

const (
	TextNode NodeType = iota
	ElementNode
)

type Props = map[string]string

type VNode struct {
	Type     NodeType
	TagName  string
	Attrs    *Attrs
	Children []*VNode
}

// hyperscript-style API: h(tagName, attrs, children)
func H(tagName string, attrs *Attrs, children ...*VNode) *VNode {
	a := atom.Lookup([]byte(tagName))
	if a == 0 {
		return &VNode{Type: TextNode, TagName: tagName}
	}
	return &VNode{Type: ElementNode, TagName: tagName, Attrs: attrs, Children: children}
}

func Text(text string) *VNode {
	return &VNode{TagName: text}
}

func (vnode *VNode) HashCode() string {
	if vnode.Type == TextNode {
		return vnode.TagName
	}
	if vnode.Attrs != nil && vnode.Attrs.Props != nil {
		return fmt.Sprintf("%s/%v", vnode.TagName, *vnode.Attrs.Props)
	}
	return fmt.Sprintf("%s/%v", vnode.TagName, Attrs{})
}
