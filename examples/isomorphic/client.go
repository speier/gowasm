package main

import (
	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	state := &State{Count: 0}
	client.Render(state, func() *vdom.VNode { return App(state) }, "root")
}
