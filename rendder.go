package geui

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"os"
)

type Renderer struct {
	canvas     *gg.Context
	nodes      []*Node
	loadedFace map[string]font.Face
}

func NewRenderer(n *Node) *Renderer {
	return &Renderer{
		canvas:     gg.NewContext(int(n.Model.Width), int(n.Model.Height)),
		loadedFace: map[string]font.Face{},
		nodes:      n.GetNodes(),
	}
}

func (render *Renderer) Render(filename string) {
	for _, node := range render.nodes {
		switch node.Type {
		case ElementNode:
			render.canvas.DrawRectangle(node.Model.RelativeX, node.Model.RelativeY, node.Model.Width, node.Model.Height)
			render.canvas.SetHexColor(node.Style.BackgroundColor)
			render.canvas.Fill()
		case CharDataNode:
			if face, ok := render.loadedFace[node.Parent.Style.FontFamily]; ok {
				render.canvas.SetFontFace(face)
			} else {
				face, err := gg.LoadFontFace(node.Parent.Style.FontFamily, node.Parent.Style.FontSize)
				if err != nil {
					panic(err)
				}
				render.canvas.SetFontFace(face)
				render.loadedFace[node.Parent.Style.FontFamily] = face
			}
			render.canvas.SetHexColor(node.Parent.Style.FontColor)
			w, h := node.Parent.Model.Width, node.Parent.Model.Height
			render.canvas.DrawStringAnchored(node.Data, node.Parent.Model.RelativeX+w/2, node.Parent.Model.RelativeY+h/2, 0.5, 0.5)
		}
	}
	fi, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	_ = render.canvas.EncodePNG(fi)
}
