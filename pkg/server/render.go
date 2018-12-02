package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/speier/gowasm/pkg/vdom"
)

func RenderToString(n *vdom.VNode) string {
	var buf bytes.Buffer
	render(&buf, n)
	return buf.String()
}

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

func render(w io.Writer, n *vdom.VNode) error {
	if x, ok := w.(writer); ok {
		return _render(x, n)
	}
	buf := bufio.NewWriter(w)
	if err := _render(buf, n); err != nil {
		return err
	}
	return buf.Flush()
}

func _render(w writer, n *vdom.VNode) error {
	if n.Type == vdom.TextNode {
		return escape(w, n.TagName)
	}

	// element node
	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(n.TagName); err != nil {
		return err
	}
	if n.Attrs != nil && n.Attrs.Props != nil {
		for k, v := range *n.Attrs.Props {
			if err := w.WriteByte(' '); err != nil {
				return err
			}
			if _, err := w.WriteString(k); err != nil {
				return err
			}
			if _, err := w.WriteString(`="`); err != nil {
				return err
			}
			if err := escape(w, v); err != nil {
				return err
			}
			if err := w.WriteByte('"'); err != nil {
				return err
			}
		}
	}
	if voidElements[n.TagName] {
		if len(n.Children) > 0 {
			return fmt.Errorf("html: void element <%s> has child nodes", n.TagName)
		}
		_, err := w.WriteString("/>")
		return err
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline beging ignored.
	if len(n.Children) > 0 {
		c := n.Children[0]
		if c != nil && c.Type == vdom.TextNode && strings.HasPrefix(c.TagName, "\n") {
			switch n.TagName {
			case "pre", "listing", "textarea":
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
		}
	}

	// Render any child nodes.
	switch n.TagName {
	case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "xmp":
		for _, c := range n.Children {
			if c.Type == vdom.TextNode {
				if _, err := w.WriteString(c.TagName); err != nil {
					return err
				}
			} else {
				if err := _render(w, c); err != nil {
					return err
				}
			}
		}
	default:
		for _, c := range n.Children {
			if err := _render(w, c); err != nil {
				return err
			}
		}
	}

	// Render the </xxx> closing tag.
	if _, err := w.WriteString("</"); err != nil {
		return err
	}
	if _, err := w.WriteString(n.TagName); err != nil {
		return err
	}
	return w.WriteByte('>')
}

const escapedChars = "&'<>\"\r"

func escape(w writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}

var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}
