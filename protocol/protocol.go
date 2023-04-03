package protocol

import "github.com/CESSProject/p2p-go/core"

type Protocol struct {
	*core.Node
	*WriteFileProtocol
	*ReadFileProtocol
	// add other protocols here...
}

func NewProtocol(node *core.Node) *Protocol {
	return &Protocol{Node: node}
}
