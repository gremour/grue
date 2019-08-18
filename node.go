package grue

// node is basic node structure containing reference
// to one parent and multiple children.
type node struct {
	parent   Node
	children []Node
}

// Node is interface for operating tree-like structure nodes.
type Node interface {
	Close()
	Foster(ch Node)
	SubNodes() []Node

	setParent(p Node)
	getParent() Node
	addChild(ch Node)
	removeChildren()
	removeChild(ch Node)
}

// Close removes node from tree.
func (n *node) Close() {
	n.removeChildren()
	if n.parent != nil {
		n.parent.removeChild(n)
		n.parent = nil
	}
}

// Foster adds child node to this, detaching from previous parent,
// if needed.
func (n *node) Foster(ch Node) {
	if ch == nil {
		return
	}
	p := ch.getParent()
	if p == n {
		return
	}
	if p != nil {
		p.removeChild(ch)
	}

	ch.setParent(n)
	n.addChild(ch)
}

func (n *node) SubNodes() []Node {
	return n.children
}

func (n *node) setParent(p Node) {
	n.parent = p
}

func (n *node) getParent() Node {
	return n.parent
}

func (n *node) addChild(ch Node) {
	for _, c := range n.children {
		if c == ch {
			// already in children.
			return
		}
	}
	n.children = append(n.children, ch)
}

func (n *node) removeChild(ch Node) {
	pch := n.children
	l := len(pch)
	for i, c := range pch {
		if c == ch {
			if i < l-1 {
				copy(pch[i:l-1], pch[i+1:])
				// 0 1 2 3 4 5
				// a b c d e f
				//     ^i=2 l=6
				//     [2:5](c to e) replaced by [3:6](d to f)
			}
			pch[l-1] = nil
			break
		}
	}
}

func (n *node) removeChildren() {
	for _, c := range n.children {
		c.setParent(nil)
		c.Close()
	}
	n.children = []Node{}
}
