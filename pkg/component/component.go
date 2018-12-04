package component

import (
	"github.com/speier/gowasm/pkg/vdom"
)

type Component interface {
	Render() *vdom.VNode
	SetUpdateHandler(func())
}

type BaseComponent struct{}
