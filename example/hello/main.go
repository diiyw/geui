package main

import (
	"github/diiyw/geui"
)

func main() {
	node := geui.LoadXML("main.xml")
	w, err := geui.NewWindow(
		node,
		geui.Title(node.Name),
		geui.Size(node.Model.Width, node.Model.Height),
	)
	if err != nil {
		panic(err)
	}
	w.Show()
}
