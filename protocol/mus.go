package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	log.Printf("receive mus[%v] req", reqMsg.Type)

	w := pbio.NewDelimitedWriter(s)
	respMsg := &pb.MusResponse{
		Code: 0,
	}
	var names Names
	var us Us
	var namespath = ""
	var uspath = ""
	var muPath = ""

	switch reqMsg.Type {
	case pb.MusType_CustomData:
		namespath = e.node.SnamesPath
		uspath = e.node.SusPath
		muPath = e.node.SmuPath
		respMsg.MuHash, err = CalcPathSHA256(e.node.SmuPath)
		if err != nil {
			s.Reset()
			log.Println(err)
			return
		}

	case pb.MusType_IdleData:
		namespath = e.node.InamesPath
		uspath = e.node.IusPath
		muPath = e.node.ImuPath
		respMsg.MuHash, err = CalcPathSHA256(e.node.ImuPath)
		if err != nil {
			s.Reset()
			log.Println(err)
			return
		}
	default:
		s.Reset()
		log.Println(err)
		return
	}

	content, err := os.ReadFile(namespath)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	err = json.Unmarshal(content, &names)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	content, err = os.ReadFile(uspath)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	err = json.Unmarshal(content, &us)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	respMsg.Name = names.Name
	respMsg.U = us.U
	err = w.WriteMsg(respMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	log.Printf("%s: mus response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	e.node.PutMuEventCh(fmt.Sprintf("%s#%s", s.Conn().RemotePeer().String(), muPath))
}
