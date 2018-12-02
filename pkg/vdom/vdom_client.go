// +build js,wasm

package vdom

import (
	"syscall/js"
)

type Events = map[string]func([]js.Value)

type Attrs struct {
	Props  *Props
	Events *Events
}
