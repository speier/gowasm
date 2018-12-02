package dom

import (
	"strings"
	"syscall/js"

	"github.com/speier/gowasm/pkg/vdom"
)

var (
	document = js.Global().Get("document")
	_oldNode *vdom.VNode
)

func QuerySelector(selector string) js.Value {
	return document.Call("querySelector", selector)
}

func Render(node *vdom.VNode, container js.Value) {
	recycle(container)
	_oldNode = Patch(_oldNode, node, container)
}

func Patch(oldNode *vdom.VNode, newNode *vdom.VNode, container js.Value) *vdom.VNode {
	patchElement(oldNode, newNode, container, 0)
	return newNode
}

func patchElement(oldNode *vdom.VNode, newNode *vdom.VNode, parent js.Value, index int) {
	if oldNode == nil {
		el := createElement(newNode)
		parent.Call("appendChild", el)
	} else if newNode == nil {
		parent.Call("removeChild", parent.Get("childNodes").Index(index))
	} else if changed(newNode, oldNode) {
		parent.Call("replaceChild", createElement(newNode), parent.Get("childNodes").Index(index))
	} else {
		newLength := len(newNode.Children)
		oldLength := len(oldNode.Children)
		for i := 0; i < newLength || i < oldLength; i++ {
			patchElement(
				oldNode.Children[i],
				newNode.Children[i],
				parent.Get("childNodes").Index(index),
				i,
			)
		}
	}
}

func changed(a *vdom.VNode, b *vdom.VNode) bool {
	return a.HashCode() != b.HashCode()
}

func createElement(node *vdom.VNode) js.Value {
	if node.Type == vdom.TextNode {
		return document.Call("createTextNode", node.TagName)
	}

	el := document.Call("createElement", node.TagName)

	if node.Attrs != nil {
		if node.Attrs.Props != nil {
			for attr, attrValue := range *node.Attrs.Props {
				el.Call("setAttribute", attr, attrValue)
			}
		}

		if node.Attrs.Events != nil {
			for eventName, handler := range *node.Attrs.Events {
				callback := js.NewCallback(func(args []js.Value) { handler() })
				el.Call("addEventListener", eventName, callback)
			}
		}
	}

	for _, child := range node.Children {
		el.Call("appendChild", createElement(child))
	}

	return el
}

func recycle(container js.Value) {
	if _oldNode != nil {
		return
	}

	el := container.Get("childNodes").Index(0)
	if el.Type() == js.TypeUndefined {
		return
	}

	_oldNode = &vdom.VNode{
		TagName: strings.ToLower(el.Get("nodeName").String()),
	}
	childNodes := el.Get("childNodes")
	childLen := childNodes.Length()
	for i := 0; i < childLen; i++ {
		c := childNodes.Index(i)
		if c.Get("nodeType").Int() == 3 {
			// text node
			_oldNode.Children = append(_oldNode.Children, &vdom.VNode{
				Type:    vdom.TextNode,
				TagName: c.Get("nodeValue").String(),
			})
		} else {
			recycle(c)
		}
	}
}
