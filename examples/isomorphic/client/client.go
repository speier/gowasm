package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o ../static/app.wasm ."

import (
	"strconv"
	"syscall/js"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	// state := &app.State{Count: 0}
	// client.Render(state, func() *vdom.VNode { return app.Init(state) }, "root")

	root := dom.QuerySelector("#root")

	var render func(state int)
	view := func(count int) *vdom.VNode {
		return vdom.H("div", nil,
			vdom.H("h1", nil, vdom.Text(strconv.Itoa(count))),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Ev{"click": func(args []js.Value) { count -= 1; render(count) }}}, vdom.Text("-")),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Ev{"click": func(args []js.Value) { count += 1; render(count) }}}, vdom.Text("+")),
		)
	}

	app := func(view func(count int) *vdom.VNode, container js.Value, node *vdom.VNode) func(state int) {
		return func(state int) {
			// node = dom.Patch(node, view(state), container)
			dom.Render(view(state), container)
		}
	}
	render = app(view, root, nil)
	render(0)

	js.Global().Get("window").Call("requestAnimationFrame", js.NewCallback(func(args []js.Value) {
		println("requestAnimationFrame")
	}))
	select {}
	// done := make(chan struct{}, 0)
	// <-done
}

var (
	UpdateHandler func()
)

func Action(f func()) {
	f()
	Update()
}

func Update() {
	UpdateHandler()
}

type State struct {
	count int
}

var state = State{}

func increment() {
	state.count++
}

func decrement() {
	state.count--
}

func Counter() *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("div", nil, vdom.Text(strconv.Itoa(state.count))),
		vdom.H("div", nil), // vdom.H("button", vdom.Attributes{"onclick": Action(increment)}, vdom.Text("+")),
		// vdom.H("button", vdom.Attributes{"onclick": Action(decrement)}, vdom.Text("-")),

	)
}
