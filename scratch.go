package main

import (
	"fmt"
	"strconv"

	"github.com/speier/gowasm/pkg/server"
	"github.com/speier/gowasm/pkg/vdom"
)

// user
type State struct {
	Count int
}

func main() {
	render(v1)
	render(v2("bar"))

	c := &TestComponent{1}
	render(c.Render)
	c.Inc(4)
}

func v1() *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("foo", nil),
	)
}

func v2(s string) func() *vdom.VNode {
	return func() *vdom.VNode {
		return vdom.H("div", nil,
			vdom.H(s, nil),
		)
	}
}

// pkg
type Context interface {
	Update()
}

func render(view func() *vdom.VNode) {
	str := server.RenderToString(view())
	fmt.Println(str)
}

func renderComp(c Component) {
	str := server.RenderToString(c.Render())
	fmt.Println(str)
}

type Component interface {
	Render() *vdom.VNode
}

type TestComponent struct {
	Count int
}

func (c *TestComponent) Inc(by int) {
	c.Count += by
	render(c.Render)
}

func (c *TestComponent) Render() *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H(strconv.Itoa(c.Count), nil),
	)
}
