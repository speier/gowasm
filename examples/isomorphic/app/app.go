package app

import (
	"strconv"
	"syscall/js"

	"github.com/speier/gowasm/pkg/vdom"
)

type State struct {
	Count int
}

type Actions struct {
	State  *State
	Update func()
}

func (a *Actions) Down(value int) {
	a.State.Count -= value
	a.Update()
}

func (a *Actions) Up(value int) {
	a.State.Count += value
	a.Update()
}

func View(state *State, actions *Actions) *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("h1", nil, vdom.Text(strconv.Itoa(state.Count))),
		vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func(args []js.Value) { actions.Down(1) }}}, vdom.Text("-")),
		vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func(args []js.Value) { actions.Up(1) }}}, vdom.Text("+")),
	)
}

var Home HomeComponent

type HomeComponent struct {
	Message string
}

func (h HomeComponent) Render() string {
	return `<h2>Home</h2>`
}

var About AboutComponent

type AboutComponent struct{}

func (a AboutComponent) Render() string {
	return `<h2>About</h2>`
}
