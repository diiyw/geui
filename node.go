package geui

// A NodeType is the type of a Node.
type NodeType uint

const (
	// DocumentNode is a document object that, as the root of the document tree,
	// provides access to the entire XML document.
	DocumentNode NodeType = iota
	// DeclarationNode is the document type declaration, indicated by the following
	// tag (for example, <!DOCTYPE...> ).
	DeclarationNode
	// ElementNode is an element (for example, <item> ).
	ElementNode
	// TextNode is the text content of a node.
	TextNode
	// CharDataNode node <![CDATA[content]]>
	CharDataNode
	// CommentNode a comment (for example, <!-- my comment --> ).
	CommentNode
)

// A Node consists of a NodeType and some Data (tag name for
// element nodes, content for text) and are part of a tree of Nodes.
type Node struct {
	Model *Model

	ID, Name string

	Parent                   *Node
	PrevSibling, NextSibling *Node
	FirstChild, LastChild    *Node

	Type  NodeType
	Data  string
	Value []rune

	Style *CSStyle

	level int // node level in the tree
}

// AddChild adds a new node 'n' to a node 'parent' as its last child.
func AddChild(parent, n *Node) {
	n.Parent = parent
	n.NextSibling = nil
	if parent.FirstChild == nil {
		parent.FirstChild = n
		n.PrevSibling = nil
	} else {
		parent.LastChild.NextSibling = n
		n.PrevSibling = parent.LastChild
	}

	parent.LastChild = n
}

// AddSibling adds a new node 'n' as a sibling of a given node 'sibling'.
// Note it is not necessarily true that the new node 'n' would be added
// immediately after 'sibling'. If 'sibling' isn't the last child of its
// parent, then the new node 'n' will be added at the end of the sibling
// chain of their parent.
func AddSibling(sibling, n *Node) {
	for t := sibling.NextSibling; t != nil; t = t.NextSibling {
		sibling = t
	}
	n.Parent = sibling.Parent
	sibling.NextSibling = n
	n.PrevSibling = sibling
	n.NextSibling = nil
	if sibling.Parent != nil {
		sibling.Parent.LastChild = n
	}
}

func (n *Node) GetNodes() (nodes []*Node) {
	nodes = make([]*Node, 0)
	nodes = append(nodes, n)
	var f func(*Node)
	f = func(n *Node) {
		if n != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				nodes = append(nodes, c)
				f(c)
			}
		}
	}
	f(n)
	return
}

func (n *Node) GetActiveNode(x, y float64) *Node {
	var f func(*Node) *Node
	f = func(n *Node) *Node {
		if n != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Focused(x, y) {
					return c
				}
				f(c)
			}
		}
		return nil
	}
	return f(n)
}
