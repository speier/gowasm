package component

import (
	"github.com/speier/gowasm/pkg/vdom"
)

type Context interface {
	Update()
}

type Component interface {
	Init(ctx Context)
	Render() *vdom.VNode
}
