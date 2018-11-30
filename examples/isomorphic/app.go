package main

import (
	"strconv"

	"github.com/speier/gowasm/pkg/vdom"
)

var h = vdom.H

func App(state *State) *vdom.VNode {
	return h("div", nil,
		h("h1", nil, vdom.HText(strconv.Itoa(state.Count))),
	)
}

var Home HomeComponent

type HomeComponent struct {
	Message string
}

func (h HomeComponent) Render() string {
	return `home`
}

var About AboutComponent

type AboutComponent struct{}

func (a AboutComponent) Render() string {
	return `<div>
		<h2>About Component</h2>
	</div>`
}
