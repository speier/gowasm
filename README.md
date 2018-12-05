# Go WASM

Toolkit for building web apps with Go and WebAssembly.

Modular packages, supports building virtual dom nodes with a hyperscript style api, rendering to HTML string on the server, and DOM elements in the browser.

Features:

- virtual dom _(with patching and possible hydration)_
- hyperscript
- server-side rendering
- routing

## Getting Started

Install dependencies:

```sh
$ go get github.com/speier/gowasm
```

Write a component:

```go
package main

import (
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

func main() {
	view := vdom.H("div", nil,
		vdom.H("h1", nil, vdom.Text("Hello World!")),
	)

	dom.Render(view, dom.QuerySelector("body"))
}
```

Run with dev server:

```sh
$ gowasm serve
```


## Examples

To run the [basic](examples/basic) example with dev server:

```sh
$ gowasm serve examples/basic
```

To run the [isomorphic](examples/isomorphic) example:

```sh
# generate wasm
$ go generate ./examples/isomorphic/...
# run the server
$ go run examples/isomorphic/server.go
```

## License

[MIT](LICENSE)
