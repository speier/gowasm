// $ go test -v pkg/jsx/*.go
package jsx

import (
	"testing"
)

func TestTransformJSX(t *testing.T) {
	js, err := TransformJSX(
		`/** @jsx h */
	<div>
		Hello {state.name}
	</div>`)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(js)
}
