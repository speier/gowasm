package client

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func Render(view func() *vdom.VNode, container js.Value) {
	done := make(chan struct{}, 0)

	// just to keep go wasm alive even when there is no event subscription in the app
	js.Global().Get("window").Call("requestAnimationFrame", js.NewCallback(func(args []js.Value) {
		// println("requestAnimationFrame")
	}))

	dom.Render(view(), container)
	<-done
}
