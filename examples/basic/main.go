package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o basic.wasm ."

import (
	"strconv"
	"time"

	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/component"
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	// simple components
	dom.Render(HelloMessage("World!"), dom.QuerySelector("#s1"))
	dom.Render(TextBox("type here..."), dom.QuerySelector("#s2"))

	// stateful components
	go client.Mount(TimerFn(0), dom.QuerySelector("#s3"))
	client.Mount(Timer(0), dom.QuerySelector("#s4"))
}

var h = vdom.H

func HelloMessage(name string) *vdom.VNode {
	return h("div", nil, h("Hello ", nil), h(name, nil))
}

func TextBox(placeholder string) *vdom.VNode {
	textbox := h("input", &vdom.Attrs{Props: &vdom.Props{"type": "text", "placeholder": placeholder}})
	textbox.OnCreate = func() {
		println("oncreate")
	}
	return textbox
}

func TimerFn(seconds int) func(updateHandler func()) *vdom.VNode {
	update := func() {}

	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			seconds++
			update()
		}
	}()

	return func(updateHandler func()) *vdom.VNode {
		update = updateHandler
		return h("div", nil, h("Seconds: ", nil), h(strconv.Itoa(seconds), nil))
	}
}

// component

type TimerComp struct {
	seconds int
	ctx     component.Context
}

func Timer(seconds int) *TimerComp {
	t := &TimerComp{seconds: seconds}
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			t.seconds++
			t.ctx.Update()
		}
	}()
	return t
}

func (t *TimerComp) Init(ctx component.Context) {
	t.ctx = ctx
}

func (t *TimerComp) Render() *vdom.VNode {
	return h("div", nil, h("Seconds: ", nil), h(strconv.Itoa(t.seconds), nil))
}
