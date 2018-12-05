package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o ../static/app.wasm ."

import (
	"encoding/json"
	"syscall/js"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"

	"github.com/speier/gowasm/examples/isomorphic/app"
)

func main() {
	var state *app.State
	initState := []byte(dom.Window.Get("initialState").String())
	json.Unmarshal(initState, &state)

	actions := &app.Actions{}
	App(state, actions, app.View, dom.QuerySelector("#root"))
}

func App(state *app.State, actions *app.Actions, view func(state *app.State, actions *app.Actions) *vdom.VNode, container js.Value) {
	renderFactory := func(view func(state *app.State, actions *app.Actions) *vdom.VNode, container js.Value, node *vdom.VNode) func(state *app.State) {
		return func(state *app.State) {
			node = dom.Patch(node, view(state, actions), container)
		}
	}

	render := renderFactory(view, container, nil)

	actions.State = state
	actions.Update = func() { render(state) }

	render(state)

	select {}
}
