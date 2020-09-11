package geui

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	n := LoadXML("testdata/main.xml")
	r := NewRenderer(n)
	r.Render("testdata/main.png")
}
