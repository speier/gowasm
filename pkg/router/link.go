package router

import (
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func Link(href string, text string) *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("a", &vdom.Attrs{Props: &vdom.Props{"href": href}, Events: &vdom.Events{"click": func() { navigate(href) }}},
			vdom.H(text, nil),
		),
	)
}

func navigate(to string) {
	pathname := dom.Window.Get("location").Get("pathname").String()
	if to != pathname {
		r.Navigate(to)
		dom.Window.Get("history").Call("pushState", pathname, "", to)
	}
}
