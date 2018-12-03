package main

//go:generate sh -c "GOOS=js GOARCH=wasm go build -o basic.wasm ."

import (
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	view := vdom.H("div", nil,
		vdom.H("h1", nil, vdom.Text("Hello World!")),
	)

	dom.Render(view, dom.QuerySelector("#root"))
}
