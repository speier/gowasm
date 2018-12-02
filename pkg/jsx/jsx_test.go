// $ go test -v pkg/jsx/*.go
package jsx

import (
	"io/ioutil"
	"testing"
)

var testFiles = []string{
	"hello.jsx",
	"todo_item.jsx",
	"todo_list.jsx",
}

func TestTransformJSX(t *testing.T) {
	for _, f := range testFiles {
		transformFile(t, "test/"+f)
	}
}

func transformFile(t *testing.T, filename string) {
	jsx, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := TransformJSX(string(jsx))
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(res)
}
