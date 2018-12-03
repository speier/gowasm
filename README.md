# Go WASM

Toolkit for building web apps with Go and WebAssembly.

At the core of the toolkit it's a lightweight UI library, consist of a hyperscript style api for building virtual dom nodes.

Modular components supports to render virtual dom nodes to an HTML string on the server, and DOM elements in the browser (with patching and possible hydration).

## Getting Started

Install dependencies:

```sh
$ GO111MODULE=on go get
```

Writing components:

```go
func main() {
	view := vdom.H("div", nil,
		vdom.H("h1", nil, vdom.Text("Hello World!")),
	)

	dom.Render(view, dom.QuerySelector("#root"))
}
```

## Examples

To run the [isomorphic](examples/isomorphic) example generate WebAssembly:

```sh
$ go generate ./examples/isomorphic/...
```

and run the server:

```sh
$ go run examples/isomorphic/server.go
```

## License

[MIT](LICENSE)
