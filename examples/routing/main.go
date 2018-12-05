package main

import (
	"github.com/speier/gowasm/pkg/client"
	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/router"

	"github.com/speier/gowasm/examples/routing/pages"
)

func main() {
	r := router.New()

	r.Route("/", pages.Index)
	r.Route("/about", pages.About)

	client.Mount(r.Switch(), dom.QuerySelector("body"))
}
