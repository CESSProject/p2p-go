package protocol

import (
	"log"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-msgio/pbio"
)

const PushTag_Protocol = "/kldr/tagpush/1"

type PushTagProtocol struct {
	node *core.Node
}

func NewPushTagProtocol(node *core.Node) *PushTagProtocol {
	e := PushTagProtocol{node: node}
	node.SetStreamHandler(PushTag_Protocol, e.onPushTagRequest)
	return &e
}

// remote peer requests handler
func (e *PushTagProtocol) onPushTagRequest(s network.Stream) {
	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	reqMsg := &pb.TagPushRequest{}
	err := r.ReadMsg(reqMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}

	log.Printf("receive push tag req: %s", s.Conn().RemotePeer())

	w := pbio.NewDelimitedWriter(s)
	respMsg := &pb.TagPushResponse{
		Code: 0,
	}
	w.WriteMsg(respMsg)

	log.Printf("%s: push tag response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
}
