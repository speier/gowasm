package pages

import (
	"github.com/speier/gowasm/pkg/router"
	"github.com/speier/gowasm/pkg/vdom"
)

func Index() *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("h1", nil,
			vdom.H("Hello 👋", nil),
		),
		router.Link("/about", "About"),
	)
}
