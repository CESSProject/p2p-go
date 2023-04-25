package protocol

import (
	"context"
	"fmt"
	"log"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const IDLE_PROTOCOL = "/idle/req/1"

type IdleProtocol struct {
	node     *core.Node                 // local host
	requests map[string]*pb.IdleRequest // used to access request data from response handlers
}

func NewIdleProtocol(node *core.Node) *IdleProtocol {
	e := IdleProtocol{node: node, requests: make(map[string]*pb.IdleRequest)}
	//node.SetStreamHandler(FILE_PROTOCOL, e.onIdleRequest)
	return &e
}

func (e *IdleProtocol) IdleReq(peerId peer.ID, peerindex uint64, sign []byte) error {
	log.Printf("Sending file req to: %s", peerId)

	reqMsg := &pb.IdleRequest{
		PeerIndex: peerindex,
		Sign:      sign,
	}

	s, err := e.node.NewStream(context.Background(), peerId, IDLE_PROTOCOL)
	if err != nil {
		return err
	}
	defer s.Close()

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return err
	}
	var respMsg = &pb.FileResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return err
	}
	if respMsg.Code != P2PResponseOK {
		return fmt.Errorf("File req returns failed: %d", respMsg.Code)
	}

	log.Printf("Idle req suc")
	return nil
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
