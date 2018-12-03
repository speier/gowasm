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
	window := js.Global().Get("window")
	initState := []byte(window.Get("initialState").String())

	var state *app.State
	json.Unmarshal(initState, &state)

	actions := &app.Actions{}
	App(state, actions, app.View, dom.QuerySelector("#root"))
}

func App(state *app.State, actions *app.Actions, view func(state *app.State, actions *app.Actions) *vdom.VNode, container js.Value) {
	renderFactory := func(view func(state *app.State, actions *app.Actions) *vdom.VNode, container js.Value) func(state *app.State) {
		return func(state *app.State) {
			dom.Render(view(state, actions), container)
		}
	}

	render := renderFactory(view, container)

	actions.State = state
	actions.Update = func() { render(state) }

	render(state)

	select {}
}
