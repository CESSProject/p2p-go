package protocol

import (
	"context"
	"log"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const IdleDataTag_Protocol = "/kldr/idtg/1"

type IdleDataTagProtocol struct {
	node *core.Node
}

func NewIdleDataTagProtocol(node *core.Node) *IdleDataTagProtocol {
	e := IdleDataTagProtocol{node: node}
	return &e
}

func (e *IdleDataTagProtocol) IdleReq(peerId peer.ID, filesize, blocknum uint64, pubkey, sign []byte) (uint32, error) {
	log.Printf("Sending file req to: %s", peerId)

	reqMsg := &pb.IdleDataTagRequest{
		FileSize:  filesize,
		BlockNum:  blocknum,
		Publickey: pubkey,
		Sign:      sign,
	}

	s, err := e.node.NewStream(context.Background(), peerId, IdleDataTag_Protocol)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	var respMsg = &pb.IdleDataTagResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}
	if respMsg.Code == 0 {
		e.node.SetIdleFileTee(string(peerId))
	}
	log.Printf("Idle req suc")
	return respMsg.Code, nil
}
