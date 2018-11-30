package babel

//go:generate mkdir -p tmp
//go:generate curl https://unpkg.com/babel-standalone@6.26.0/babel.min.js -o tmp/babel.min.js
//go:generate rice embed-go

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/dop251/goja"
)

func Transform(src string, opts map[string]interface{}) (string, error) {
	babelJS := rice.MustFindBox("tmp").MustString("babel.min.js")

	vm := goja.New()
	_, err := vm.RunString(babelJS)
	if err != nil {
		return "", err
	}

	var transform goja.Callable
	babel := vm.Get("Babel")
	err = vm.ExportTo(babel.ToObject(vm).Get("transform"), &transform)
	if err != nil {
		return "", err
	}

	res, err := transform(babel, vm.ToValue(src), vm.ToValue(opts))
	if err != nil {
		return "", err
	}

	return res.ToObject(vm).Get("code").String(), nil
}
