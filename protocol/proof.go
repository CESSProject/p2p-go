package protocol

import (
	"context"
	"log"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const AggrProof_PROTOCOL = "/kldr/apv/1"

type AggrProofProtocol struct {
	node *core.Node
}

func NewAggrProofProtocol(node *core.Node) *AggrProofProtocol {
	e := AggrProofProtocol{node: node}
	return &e
}

func (e *AggrProofProtocol) AggrProofReq(peerId peer.ID, ihash, shash []byte, qslice []*pb.Qslice, puk, sign []byte) (uint32, error) {
	log.Printf("Sending AggrProof req to: %s", peerId)

	s, err := e.node.NewStream(context.Background(), peerId, AggrProof_PROTOCOL)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.AggrProofRequest{
		IdleProofFileHash:    ihash,
		ServiceProofFileHash: shash,
		Qslice:               qslice,
		Publickey:            puk,
		Sign:                 sign,
	}

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	respMsg := &pb.AggrProofResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	log.Printf("AggrProof req resp code: %d", respMsg.Code)
	return respMsg.Code, nil
}
