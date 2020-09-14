package geui

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
	"os"
)

type Renderer struct {
	canvas      *gg.Context
	nodes       []*Node
	loadedFonts map[string]*truetype.Font
}

func NewRenderer(n *Node) *Renderer {
	return &Renderer{
		canvas:      gg.NewContext(int(n.Model.Width), int(n.Model.Height)),
		loadedFonts: map[string]*truetype.Font{},
		nodes:       n.GetNodes(),
	}
}

func (render *Renderer) Render(filename string) {

	var setFontFace = func(n *Node) {
		f, ok := render.loadedFonts[n.Parent.Style.FontFamily]
		if !ok {
			fontBytes, err := ioutil.ReadFile(n.Parent.Style.FontFamily)
			if err != nil {
				log.Println(err)
				return
			}
			f, err = truetype.Parse(fontBytes)
			if err != nil {
				log.Println(err)
				return
			}
			render.loadedFonts[n.Parent.Style.FontFamily] = f
		}
		face := truetype.NewFace(f, &truetype.Options{
			Size: n.Parent.Style.FontSize,
		})
		render.canvas.SetFontFace(face)
	}
	for _, n := range render.nodes {
		switch n.Type {
		case ElementNode:
			switch n.Data {
			case "input":
				render.canvas.SetHexColor(n.Style.BorderColor)
				render.canvas.SetLineWidth(n.Style.BorderWidth)
				x, y := n.Model.RelativeX+0.5, n.Model.RelativeY+0.5
				w, h := n.Model.Width, n.Model.Height
				render.canvas.DrawLine(x, y, x+w, y)
				render.canvas.DrawLine(x, y, x, y+h)
				render.canvas.DrawLine(x+w, y, x+w, y+h)
				render.canvas.DrawLine(x, y+h, x+w, y+h)
				fw, _ := render.canvas.MeasureString(string(n.Value))
				fw /= 2
				afw := fw / float64(len(n.Value))
				v := n.Value
				if fw > w {
					v = n.Value[len(n.Value)-int(w/afw)-1:]
					fw = w - 5
				}
				if len(n.Value) != 0 {
					render.canvas.DrawLine(fw+x, y+5, fw+x, h+y-5)
				}
				render.canvas.Stroke()
				setFontFace(n)
				render.canvas.SetHexColor(n.Style.FontColor)
				if len(n.Value) != 0 {
					render.canvas.DrawStringAnchored(string(v), n.Model.RelativeX+6, n.Model.RelativeY+h/2, 0, 0.4)
				}
			case "checkbox":

			default:
				render.canvas.DrawRectangle(n.Model.RelativeX, n.Model.RelativeY, n.Model.Width, n.Model.Height)
				render.canvas.SetHexColor(n.Style.BackgroundColor)
				render.canvas.Fill()
			}
		case CharDataNode:
			setFontFace(n)
			render.canvas.SetHexColor(n.Parent.Style.FontColor)
			w, h := n.Parent.Model.Width, n.Parent.Model.Height
			render.canvas.DrawStringAnchored(n.Data, n.Parent.Model.RelativeX+w/2, n.Parent.Model.RelativeY+h/2, 0.5, 0.5)
		}
	}
	fi, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	_ = render.canvas.EncodePNG(fi)
}
