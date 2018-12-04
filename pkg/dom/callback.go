package dom

import (
	"syscall/js"
)

// callbacks wrapper for future usage maybe

type (
	Value             = js.Value
	Callback          = js.Callback
	EventCallbackFlag = js.EventCallbackFlag
)

const (
	PreventDefault           = js.PreventDefault
	StopPropagation          = js.StopPropagation
	StopImmediatePropagation = js.StopImmediatePropagation
)

func NewCallback(f func([]Value)) Callback {
	return js.NewCallback(f)
}

func NewEventCallback(flags EventCallbackFlag, fn func(event Value)) Callback {
	return js.NewEventCallback(flags, fn)
}
