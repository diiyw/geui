package geui

import "testing"

func TestWindow(t *testing.T) {
	node := LoadXML("testdata/main.xml")
	w, err := NewWindow(
		node,
		Title(node.Name),
		Size(node.Model.Width, node.Model.Height),
	)
	if err != nil {
		panic(err)
	}
	w.Show()
}
