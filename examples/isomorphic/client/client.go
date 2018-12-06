package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o ../static/app.wasm ."

import (
	"encoding/json"

	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/component"
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"

	"github.com/speier/gowasm/examples/isomorphic/app"
)

func main() {
	state := &app.State{}
	actions := &app.Actions{State: state}

	initState := dom.Window.Get("initialState").String()
	json.Unmarshal([]byte(initState), state)

	app := AppWrapper(state, actions, app.View)
	client.Mount(app, dom.QuerySelector("#root"))

	// App(state, actions, app.View, dom.QuerySelector("#root"))
}

// func App(state *app.State, actions *app.Actions, view func(state *app.State, actions *app.Actions) *vdom.VNode, container js.Value) {
// 	renderFactory := func() func(state *app.State) {
// 		var node *vdom.VNode
// 		return func(state *app.State) {
// 			node = dom.Patch(node, view(state, actions), container)
// 		}
// 	}

// 	render := renderFactory()
// 	actions.Update = func() { render(state) }
// 	render(state)

// 	select {}
// }

func AppWrapper(state *app.State, actions *app.Actions, view func(state *app.State, actions *app.Actions) *vdom.VNode) *AppComponent {
	a := &AppComponent{state, actions, view}
	actions.Update = func() { a.Render() }
	return a
}

type AppComponent struct {
	state   *app.State
	actions *app.Actions
	view    func(state *app.State, actions *app.Actions) *vdom.VNode
}

func (a *AppComponent) Init(ctx component.Context) {
	a.actions.Update = ctx.Update
}

func (a *AppComponent) Render() *vdom.VNode {
	return a.view(a.state, a.actions)
}
