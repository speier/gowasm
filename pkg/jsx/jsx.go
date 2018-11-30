package jsx

import (
	"github.com/speier/gowasm/pkg/babel"
)

func TransformJSX(jsx string) (string, error) {
	js, err := babel.Transform(jsx, map[string]interface{}{
		"plugins": []string{"transform-react-jsx"},
	})
	if err != nil {
		return "", err
	}
	return js, nil
}
