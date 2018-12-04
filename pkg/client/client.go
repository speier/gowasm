package client

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/component"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func Render(view *vdom.VNode, container js.Value) {
	dom.Render(view, container)
}

type client struct {
	updateScheduled bool
	renderFunction  func()
}

func Mount(i interface{}, container js.Value) {
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
	case component.Component:
		render = v.Render
		v.SetUpdateHandler(c.renderFunction)
	case func() *vdom.VNode:
		render = v
	case func(func()) *vdom.VNode:
		render = func() *vdom.VNode { return v(c.renderFunction) }
	}

	c.renderFunction()

	run()
}

func run() {
	done := make(chan bool)
	unloading := js.NewEventCallback(js.PreventDefault, func(event js.Value) {
		done <- true
	})
	defer unloading.Release()
	dom.Window.Call("addEventListener", "beforeunload", unloading)
	<-done
}
