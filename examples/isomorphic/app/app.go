package app

import (
	"strconv"

	"github.com/speier/gowasm/pkg/vdom"
)

type State struct {
	Count int `json:"count"`
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
		vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { actions.Down(1) }}}, vdom.Text("-")),
		vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { actions.Up(1) }}}, vdom.Text("+")),
	)
}
