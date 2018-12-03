package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o patch.wasm ."

import (
	"strconv"
	"syscall/js"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

// inspired by https://github.com/jorgebucaran/superfine#usage
func main() {
	var render func(state int)
	view := func(count int) *vdom.VNode {
		return vdom.H("div", nil,
			vdom.H("h1", nil, vdom.Text(strconv.Itoa(count))),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { count -= 1; render(count) }}}, vdom.Text("-")),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { count += 1; render(count) }}}, vdom.Text("+")),
		)
	}

	app := func(view func(count int) *vdom.VNode, container js.Value, node *vdom.VNode) func(state int) {
		return func(state int) {
			node = dom.Patch(node, view(state), container)
		}
	}

	container := dom.QuerySelector("#root")
	render = app(view, container, nil)
	render(0)

	select {}
}
