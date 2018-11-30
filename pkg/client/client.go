package client

import (
	"syscall/js"

	"github.com/speier/gowasm/pkg/server"
	"github.com/speier/gowasm/pkg/vdom"
)

func Render(state interface{}, view func() *vdom.VNode, el string) {
	html := server.RenderToString(view())
	js.Global().Get("document").Call("getElementById", el).Set("innerHTML", html)

	// done := make(chan struct{}, 0)
	// <-done
}
