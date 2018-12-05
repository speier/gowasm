package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o channels.wasm ."

import (
	"strconv"
	"time"

	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	go client.Render(TimerChan(0), dom.QuerySelector("#s1"))
	client.Render(AddSub(0), dom.QuerySelector("#s2"))
}

var h = vdom.H

func TimerChan(i int) chan *vdom.VNode {
	seconds := make(chan int, 1)
	seconds <- i

	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			i++
			seconds <- i
		}
	}()

	res := make(chan *vdom.VNode)
	go func() {
		for s := range seconds {
			res <- h("div", nil, h("Seconds: ", nil), h(strconv.Itoa(s), nil))
		}
	}()
	return res
}

func AddSub(count int) chan *vdom.VNode {
	res := make(chan *vdom.VNode, 1)
	update := func() {}

	update = func() {
		res <- vdom.H("div", nil,
			vdom.H("h1", nil, vdom.Text(strconv.Itoa(count))),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { count -= 1; update() }}}, vdom.Text("-")),
			vdom.H("button", &vdom.Attrs{Events: &vdom.Events{"click": func() { count += 1; update() }}}, vdom.Text("+")),
		)
	}

	update()

	return res
}
