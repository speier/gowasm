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
	// create nodes for demo components below
	client.Render(sandbox, dom.QuerySelector("body"))

	// simple component
	client.Render(HelloMessage("World!"), dom.QuerySelector("#s1"))

	// stateful component
	go client.Mount(TimerFn(0), dom.QuerySelector("#s2"))
	client.Mount(Timer(0), dom.QuerySelector("#s3"))
}

var h = vdom.H

func sandbox() *vdom.VNode {
	return h("div", nil,
		h("div", &vdom.Attrs{Props: &vdom.Props{"id": "s1"}}),
		h("div", &vdom.Attrs{Props: &vdom.Props{"id": "s2"}}),
		h("div", &vdom.Attrs{Props: &vdom.Props{"id": "s3"}}),
	)
}

func HelloMessage(name string) *vdom.VNode {
	return h("h2", nil, h("Hello ", nil), h(name, nil))
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
		return h("h2", nil, h("Seconds: ", nil), h(strconv.Itoa(seconds), nil))
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
	return h("h2", nil, h("Seconds: ", nil), h(strconv.Itoa(t.seconds), nil))
}

func (t *TimerComp) SetUpdateHandler(updateHandler func()) {
	t.update = updateHandler
}
