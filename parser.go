package geui

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

func LoadXML(v string) *Node {
	fi, err := os.Open(v)
	if err != nil {
		panic(err)
	}
	p := newParser(fi)
	for {
		_, err := p.parse()
		if err == io.EOF {
			p.doc.FirstChild.NextSibling.Parent = nil
			p.doc.FirstChild.NextSibling.PrevSibling = nil
			return p.doc.FirstChild.NextSibling
		}
		if err != nil {
			panic(err)
		}
	}
}

type parser struct {
	decoder *xml.Decoder
	doc     *Node
	level   int
	prev    *Node
}

func newParser(r io.Reader) *parser {
	p := &parser{
		decoder: xml.NewDecoder(r),
		doc:     &Node{Type: DocumentNode, Model: new(Model)},
		level:   0,
	}
	p.decoder.CharsetReader = charset.NewReaderLabel
	p.prev = p.doc
	return p
}

func (p *parser) parse() (*Node, error) {
	for {
		tok, err := p.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if p.level == 0 {
				// missing XML declaration
				node := &Node{Type: DeclarationNode, Data: "xml", level: 1, Model: new(Model)}
				AddChild(p.prev, node)
				p.level = 1
				p.prev = node
			}

			node := &Node{
				Type:  ElementNode,
				Model: new(Model),
				Data:  tok.Name.Local,
				Style: NewStyle(),
				level: p.level,
			}
			for _, attr := range tok.Attr {
				parseAttr(node, attr.Name.Local, attr.Value)
			}
			if p.level == p.prev.level {
				AddSibling(p.prev, node)
			} else if p.level > p.prev.level {
				AddChild(p.prev, node)
			} else if p.level < p.prev.level {
				for i := p.prev.level - p.level; i > 1; i-- {
					p.prev = p.prev.Parent
				}
				AddSibling(p.prev.Parent, node)
			}
			parse(node)
			p.prev = node
			p.level++
		case xml.EndElement:
			p.level--
		case xml.CharData:
			v := strings.TrimSpace(string(tok))
			if v == "" {
				continue
			}
			node := &Node{Type: CharDataNode, Data: v, level: p.level, Model: new(Model)}
			if p.level == p.prev.level {
				AddSibling(p.prev, node)
			} else if p.level > p.prev.level {
				AddChild(p.prev, node)
			} else if p.level < p.prev.level {
				for i := p.prev.level - p.level; i > 1; i-- {
					p.prev = p.prev.Parent
				}
				AddSibling(p.prev.Parent, node)
			}
		case xml.Comment:
			node := &Node{Type: CommentNode, Data: string(tok), level: p.level, Model: new(Model)}
			if p.level == p.prev.level {
				AddSibling(p.prev, node)
			} else if p.level > p.prev.level {
				AddChild(p.prev, node)
			} else if p.level < p.prev.level {
				for i := p.prev.level - p.level; i > 1; i-- {
					p.prev = p.prev.Parent
				}
				AddSibling(p.prev.Parent, node)
			}
		case xml.Directive:
		}
	}
}

func parserXY(v string) (x, y float64) {
	xy := strings.Split(v, ",")
	if len(xy) == 0 {
		return
	}
	if len(xy) == 1 {
		x, _ = strconv.ParseFloat(xy[0], 64)
	}
	if len(xy) == 2 {
		y, _ = strconv.ParseFloat(xy[1], 64)
	}
	return
}

func parseAttr(node *Node, key, val string) {
	switch key {
	case "name":
		node.Name = val
	case "id":
		node.ID = val
	case "xy":
		node.Model.X, node.Model.Y = parserXY(val)
	case "rel-xy":
		node.Model.RelativeX, node.Model.RelativeY = parserXY(val)
	case "width":
		node.Model.Width, _ = strconv.ParseFloat(val, 64)
	case "height":
		node.Model.Height, _ = strconv.ParseFloat(val, 64)
	case "value":
		node.Value = []rune(val)
	case "style":
		parseInlineStyle(node, val)
	}
}

func parse(node *Node) {
	// position
	if node.Model.X == 0 && node.Model.Y == 0 && node.Data != "window" {
		if node.Parent != nil {
			node.Model.RelativeX = node.Parent.Model.X + 10
			node.Model.RelativeY = node.Parent.Model.Y + 10
			if node.PrevSibling != nil {
				node.Model.RelativeY += node.PrevSibling.Model.Height + node.PrevSibling.Model.RelativeY
			}
		}
	}
	// width
	if node.Model.Width == 0 {
		if node.Style.Width != 0 {
			node.Model.Width = node.Style.Width
		}
		// extend parent
		if node.Model.Width == 0 && node.Parent != nil && node.Parent.Model != nil {
			node.Model.Width = node.Parent.Model.Width - 10*2
		}
	}
	// height
	if node.Model.Height == 0 {
		if node.Style.Height != 0 {
			node.Model.Height = node.Style.Height
		}
	}
}
