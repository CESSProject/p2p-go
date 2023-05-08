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
	node *core.Node // local host
	//requests map[string]*pb.IdleDataTagRequest // used to access request data from response handlers
}

func NewIdleDataTagProtocol(node *core.Node) *IdleDataTagProtocol {
	e := IdleDataTagProtocol{node: node} // requests: make(map[string]*pb.IdleDataTagRequest)}
	//node.SetStreamHandler(FILE_PROTOCOL, e.onIdleRequest)
	return &e
}

func (e *IdleDataTagProtocol) IdleReq(peerId peer.ID, filesize, blocknum uint64, sign []byte) (uint32, error) {
	log.Printf("Sending file req to: %s", peerId)

	reqMsg := &pb.IdleDataTagRequest{
		FileSize: filesize,
		BlockNum: blocknum,
		Sign:     sign,
	}

	s, err := e.node.NewStream(context.Background(), peerId, IdleDataTag_Protocol)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}
	var respMsg = &pb.IdleDataTagResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	log.Printf("Idle req suc")
	return respMsg.Code, nil
}

// // remote peer requests handler
// func (e *IdleProtocol) onIdleRequest(s network.Stream) {
// 	r := pbio.NewDelimitedReader(s, IdleProtocolMsgBuf)
// 	reqMsg := &pb.IdleRequest{}
// 	err := r.ReadMsg(reqMsg)
// 	if err != nil {
// 		s.Reset()
// 		log.Println(err)
// 		return
// 	}

// 	log.Printf("receive file req: %d", reqMsg.PeerIndex)

// 	w := pbio.NewDelimitedWriter(s)
// 	respMsg := &pb.TagResponse{
// 		Code: P2PResponseOK,
// 	}

// 	defer func() {
// 		respMsg.MessageData.Timestamp = time.Now().UnixMilli()
// 		err = w.WriteMsg(respMsg)
// 		if err != nil {
// 			s.Reset()
// 			log.Println(err)
// 		}
// 	}()

// 	log.Printf("%s: File response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
// }
