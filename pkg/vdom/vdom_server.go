// +build !js,!wasm

package vdom

type Props = map[string]string

type Attrs struct {
	Props *Props
}
