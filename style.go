package geui

import (
	"github.com/gorilla/css/scanner"
	"image"
	"strconv"
)

type Align uint8

const (
	LEFT Align = iota
	RIGHT
	CENTER
)

type CSStyle struct {
	Width, Height   float64
	LineHeight      int
	FontFamily      string
	FontSize        float64
	FontColor       string
	TextAlign       Align
	BackgroundColor string
}

var (
	DefaultFontSize        float64 = 12
	DefaultFontColor               = "#E6E6E6"
	DefaultBackgroundColor         = "#4B4B4B"
)

func NewStyle() *CSStyle {
	return &CSStyle{
		FontColor:       DefaultFontColor,
		FontSize:        DefaultFontSize,
		Height:          30,
		Width:           0,
		FontFamily:      DefaultFont,
		LineHeight:      30,
		TextAlign:       CENTER,
		BackgroundColor: DefaultBackgroundColor,
	}
}

func parseInlineStyle(node *Node, v string) {
	s := scanner.New(v)
	for {
		tok := s.Next()
		if tok.Type == scanner.TokenEOF {
			break
		}
		switch tok.Type {
		case scanner.TokenIdent:
			switch tok.Value {
			case "background-color":
				s.Next()
				tok = s.Next()
				node.Style.BackgroundColor = tok.Value
			case "font-color":
				s.Next()
				tok = s.Next()
				node.Style.FontColor = tok.Value
			case "font-size":
				s.Next()
				tok = s.Next()
				node.Style.FontSize, _ = strconv.ParseFloat(tok.Value, 32)
			}
		}
	}
}

func (n *Node) Bounds() image.Rectangle {
	x, y, w, h := int(n.Model.RelativeX), int(n.Model.RelativeY),
		int(n.Model.Width), int(n.Model.Height)
	return image.Rect(
		x,
		y,
		x+w,
		y+h,
	)
}

func (n *Node) Focused(x, y float64) bool {
	p := image.Point{
		X: int(x),
		Y: int(y),
	}
	return p.In(n.Bounds())
}
