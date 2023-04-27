package protocol

import (
	"log"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-msgio/pbio"
)

const MUS_PROTOCOL = "/kldr/vdf/1"

type MusProtocol struct {
	node *core.Node // local host
	//requests map[string]*pb.MusRequest // used to access request data from response handlers
}

func NewMusProtocol(node *core.Node) *MusProtocol {
	e := MusProtocol{node: node} //requests: make(map[string]*pb.MusRequest)}
	node.SetStreamHandler(MUS_PROTOCOL, e.onMusRequest)
	return &e
}

// func (e *MusProtocol) MusReq(peerId peer.ID, filename, customdata string, blocknum int64) (uint32, error) {
// 	log.Printf("Sending tag req to: %s", peerId)

// 	if err := checkFileName(filename); err != nil {
// 		return 0, err
// 	}

// 	if err := checkCustomData(customdata); err != nil {
// 		return 0, err
// 	}

// 	s, err := e.node.NewStream(context.Background(), peerId, MUS_PROTOCOL)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer s.Close()

// 	w := pbio.NewDelimitedWriter(s)
// 	reqMsg := &pb.MusRequest{}

// 	err = w.WriteMsg(reqMsg)
// 	if err != nil {
// 		s.Reset()
// 		return 0, err
// 	}

// 	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
// 	respMsg := &pb.MusResponse{}
// 	err = r.ReadMsg(respMsg)
// 	if err != nil {
// 		s.Reset()
// 		return 0, err
// 	}

// 	log.Printf("Tag req resp code: %d", respMsg.Code)
// 	return respMsg.Code, nil
// }

// remote peer requests handler
func (e *MusProtocol) onMusRequest(s network.Stream) {
	r := pbio.NewDelimitedReader(s, MusProtocolMsgBuf)
	reqMsg := &pb.MusRequest{}
	err := r.ReadMsg(reqMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}

	log.Printf("receive mus req")

	w := pbio.NewDelimitedWriter(s)
	respMsg := &pb.MusResponse{
		Code:   0,
		Name:   []string{"123", "456"},
		U:      []string{"abc", "def"},
		MuHash: "123456789abcdef",
	}
	err = w.WriteMsg(respMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	log.Printf("%s: mus response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	e.node.PutMuEventCh(s.Conn().RemotePeer().String())
}
