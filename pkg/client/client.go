package client

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/component"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

// supports:
//  - `*vdom.VNode`
//  - `func() *vdom.VNode`
//  - `chan *vdom.VNode`
func Render(i interface{}, container js.Value) {
	switch v := i.(type) {
	case *vdom.VNode:
		dom.Render(v, container)
	case func() *vdom.VNode:
		dom.Render(v(), container)
	case chan *vdom.VNode:
		go func() {
			var node *vdom.VNode
			for n := range v {
				node = dom.Patch(node, n, container)
			}
		}()
		run()
	}
}

// supports:
//  - `func() *vdom.VNode
//  - `Component`
func Mount(i interface{}, container js.Value) {
	// default render function, not rendering anything
	render := func() *vdom.VNode { return nil }

	c := &renderContext{}
	var reqAnimFrame js.Callback
	var node *vdom.VNode
	c.renderFunction = func() {
		if !c.updateScheduled {
			c.updateScheduled = true
			reqAnimFrame = js.NewCallback(func(args []js.Value) {
				node = dom.Patch(node, render(), container)
				reqAnimFrame.Release()
				c.updateScheduled = false
			})
			dom.Window.Call("requestAnimationFrame", reqAnimFrame)
		}
	}

	switch v := i.(type) {
	case func() *vdom.VNode:
		render = v
	case func(func()) *vdom.VNode:
		if v != nil {
			render = func() *vdom.VNode { return v(c.renderFunction) }
		}
	case component.Component:
		render = v.Render
		v.Init(c)
	}

	c.renderFunction()
	run()
}

type renderContext struct {
	updateScheduled bool
	renderFunction  func()
}

func (c *renderContext) Update() {
	c.renderFunction()
}

// keep mounted app runing
func run() {
	done := make(chan bool)
	unloading := js.NewEventCallback(js.PreventDefault, func(event js.Value) {
		done <- true
	})
	defer unloading.Release()
	dom.Window.Call("addEventListener", "beforeunload", unloading)
	<-done
}
