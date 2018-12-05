## API design drafts

Main design goals:

- using virtual dom, supports rendering on the server
- hyperscript style api `H` as a core building block for the virtual dom _(templating could based on this in the future)_
- components start small, with a simple functional approach, each component is a single function returns a `VNode`

	```go
	func Welcome(name string) *vdom.VNode {
		return h("h2", nil, h("Hello ", nil), h(name, nil))
	}
	```

Lifecycle:

1. `oncreate` _== componentDidMount_
2. `onupdate` _== componentDidUpdate_
3. `onremove` _== componentWillUnmount_

Questions:

- when and how to trigger re-render on state changes?
  1. trigger manually from the component, with message passing for example
  2. trigger re-render on event handlers callback, how about non UI interactions like fetch/xhr?
  3. wire actions like hyperapp: in go reflection might involved with this approach

### Separation of layers

app:

```go
package main

import (
	".../vdom"
)

func main() {
	app := NewApp()
	app.Route('/', home)
	app.Render("#root")
}

func home(state *State) *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("h1", nil, vdom.H(state.Message)),
	)
}
```

server:

```go
..
initState := &app.State{Count: 1}
view := app.View(initState, nil)
html := server.RenderToString(view)
..
```

client:

```go
var state *app.State
initState := []byte(dom.Window.Get("initialState").String())
json.Unmarshal(initState, &state)

actions := &app.Actions{}
App(state, actions, app.View, dom.QuerySelector("#root"))
```
