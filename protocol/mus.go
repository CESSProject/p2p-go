package protocol

import (
	"context"
	"log"
	"time"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const MUS_PROTOCOL = "/mus/req/1"

type MusProtocol struct {
	node     *core.Node                // local host
	requests map[string]*pb.MusRequest // used to access request data from response handlers
}

func NewMusProtocol(node *core.Node) *MusProtocol {
	e := MusProtocol{node: node, requests: make(map[string]*pb.MusRequest)}
	node.SetStreamHandler(MUS_PROTOCOL, e.onMusRequest)
	return &e
}

func (e *MusProtocol) MusReq(peerId peer.ID, filename, customdata string, blocknum int64) (uint32, error) {
	log.Printf("Sending tag req to: %s", peerId)

	if err := checkFileName(filename); err != nil {
		return 0, err
	}

	if err := checkCustomData(customdata); err != nil {
		return 0, err
	}

	s, err := e.node.NewStream(context.Background(), peerId, TAG_PROTOCOL)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.TagRequest{
		MessageData: &pb.Messagedata{
			Timestamp: time.Now().UnixMilli(),
			Id:        uuid.NewString(),
		},
		FileName:   filename,
		CustomData: customdata,
		BlockNum:   blocknum,
	}

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	respMsg := &pb.TagResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	log.Printf("Tag req resp code: %d", respMsg.Code)
	return respMsg.Code, nil
}

// remote peer requests handler
func (e *MusProtocol) onMusRequest(s network.Stream) {
	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	reqMsg := &pb.MusRequest{}
	err := r.ReadMsg(reqMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}

	log.Printf("receive mus req, msg id: %s", reqMsg.MessageData.Id)

	w := pbio.NewDelimitedWriter(s)
	respMsg := &pb.MusResponse{
		MessageData: &pb.Messagedata{
			Timestamp: time.Now().UnixMilli(),
			Id:        reqMsg.MessageData.Id,
		},
		Code: P2PResponseOK,
	}
	w.WriteMsg(respMsg)

	log.Printf("%s: mus response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())

	//TODO: send files(mus,u,name) to s.Conn().RemotePeer()
}
