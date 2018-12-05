package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o basic.wasm ."

import (
	"strconv"
	"time"

	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	// simple component
	client.Render(HelloMessage("World!"), dom.QuerySelector("#s1"))
	client.Render(TextBox("type here..."), dom.QuerySelector("#s2"))
	// stateful component
	go client.Mount(TimerFn(0), dom.QuerySelector("#s3"))
	client.Mount(Timer(0), dom.QuerySelector("#s4"))
}

var h = vdom.H

func HelloMessage(name string) *vdom.VNode {
	return h("div", nil, h("Hello ", nil), h(name, nil))
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

func Timer(seconds int) *TimerComp {
	t := &TimerComp{seconds: seconds}
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			t.seconds++
			t.update()
		}
	}()
	return t
}

type TimerComp struct {
	seconds int
	update  func()
}

func (t *TimerComp) Render() *vdom.VNode {
	return h("div", nil, h("Seconds: ", nil), h(strconv.Itoa(t.seconds), nil))
}

func (t *TimerComp) SetUpdateHandler(updateHandler func()) {
	t.update = updateHandler
}

func TextBox(placeholder string) *vdom.VNode {
	textbox := h("input", &vdom.Attrs{Props: &vdom.Props{"type": "text", "placeholder": placeholder}})
	textbox.OnCreate = func() {
		println("oncreate")
	}
	return textbox
}
