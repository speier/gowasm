package dom

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/vdom"
)

var (
	Window   = js.Global().Get("window")
	Document = js.Global().Get("document")
)

func QuerySelector(selector string) js.Value {
	return Document.Call("querySelector", selector)
}

var oldNode *vdom.VNode

func Render(node *vdom.VNode, container js.Value) {
	oldNode = Patch(oldNode, node, container)
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
		return Document.Call("createTextNode", node.TagName)
	}

	el := Document.Call("createElement", node.TagName)

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
