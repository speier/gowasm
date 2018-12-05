package router

import (
	"sync"
	"syscall/js"

	"github.com/speier/gowasm/pkg/dom"
	"github.com/speier/gowasm/pkg/vdom"
)

var (
	// singleton
	r    *Router
	once sync.Once
)

type Router struct {
	path   string
	routes map[string]func() *vdom.VNode
	update func()
}

func New() *Router {
	once.Do(func() {
		r = &Router{
			path:   "/",
			routes: make(map[string]func() *vdom.VNode),
			update: func() {},
		}
		cb := js.NewEventCallback(0, r.handleLocationChange)
		dom.Window.Call("removeEventListener", "popstate", cb) // normally this will be in app's cleanup
		dom.Window.Call("addEventListener", "popstate", cb)
	})
	return r
}

func (r *Router) Route(path string, view func() *vdom.VNode) {
	r.routes[path] = view
}

func (r *Router) Switch() func(func()) *vdom.VNode {
	return func(update func()) *vdom.VNode {
		r.update = update
		return r.routes[r.path]()
	}
}

func (r *Router) Navigate(path string) {
	r.path = path
	r.update()
}

func (r *Router) handleLocationChange(e js.Value) {
	pathname := dom.Window.Get("location").Get("pathname").String()
	r.Navigate(pathname)
}
