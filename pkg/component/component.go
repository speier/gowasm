package component

import (
	"github.com/speier/gowasm/pkg/vdom"
)

type Component interface {
	SetState()
	Render() *vdom.VNode
}
