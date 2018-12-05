package client

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/component"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

// for stateless components
// supports:
//  - `*vdom.VNode`
//  - `func() *vdom.VNode`
func Render(i interface{}, container js.Value) {
	switch v := i.(type) {
	case *vdom.VNode:
		dom.Render(v, container)
	case func() *vdom.VNode:
		dom.Render(v(), container)
	}
}

// for stateful components
// supports:
//  - ...
func Mount(i interface{}, container js.Value) {
	// default render function, not rendering anything
	render := func() *vdom.VNode { return nil }

	c := &client{}
	var reqAnimFrame js.Callback
	c.renderFunction = func() {
		if !c.updateScheduled {
			c.updateScheduled = true
			reqAnimFrame = js.NewCallback(func(args []js.Value) {
				dom.Render(render(), container)
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
		render = func() *vdom.VNode { return v(c.renderFunction) }
	case component.Component:
		render = v.Render
		v.SetUpdateHandler(c.renderFunction)
	}

	c.renderFunction()
	run()
}

type client struct {
	updateScheduled bool
	renderFunction  func()
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
