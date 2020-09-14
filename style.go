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
	HoverColor      string
	BorderColor     string
	BorderWidth     float64
}

var (
	DefaultFontSize        float64 = 14
	DefaultFontColor               = "#666666"
	DefaultBackgroundColor         = "#4B4B4B"
	DefaultBorderColor             = "#666666"
)

func NewStyle() *CSStyle {
	return &CSStyle{
		FontColor:       DefaultFontColor,
		FontSize:        DefaultFontSize,
		Height:          35,
		Width:           0,
		FontFamily:      DefaultFont,
		LineHeight:      30,
		TextAlign:       CENTER,
		BackgroundColor: DefaultBackgroundColor,
		HoverColor:      DefaultBackgroundColor,
		BorderColor:     DefaultBorderColor,
		BorderWidth:     1,
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
			case "hover-color":
				s.Next()
				tok = s.Next()
				node.Style.HoverColor = tok.Value
			case "font-size":
				s.Next()
				tok = s.Next()
				node.Style.FontSize, _ = strconv.ParseFloat(tok.Value, 32)
			case "border-width":
				s.Next()
				tok = s.Next()
				node.Style.BorderWidth, _ = strconv.ParseFloat(tok.Value, 32)
			case "border-color":
				s.Next()
				tok = s.Next()
				node.Style.BorderColor = tok.Value
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
