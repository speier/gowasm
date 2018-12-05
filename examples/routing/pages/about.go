package pages

import (
	"github.com/speier/gowasm/pkg/router"
	"github.com/speier/gowasm/pkg/vdom"
)

func About() *vdom.VNode {
	return vdom.H("div", nil,
		vdom.H("p", nil,
			vdom.H("This is the about page", nil),
		),
		router.Link("/", "Index"),
	)
}
