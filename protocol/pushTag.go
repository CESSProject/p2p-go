package protocol

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

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
	remotePeer := s.Conn().RemotePeer().String()
	log.Println("receive push tag req: ", remotePeer)
	if e.node.GetIdleFileTee() != string(remotePeer) &&
		e.node.GetServiceFileTee() != string(remotePeer) {
		s.Reset()
		log.Println("receive invalid push tag req: ", remotePeer)
		return
	}

	respMsg := &pb.TagPushResponse{
		Code: 1,
	}
	w := pbio.NewDelimitedWriter(s)

	switch reqMsg.GetResult().(type) {
	case *pb.TagPushRequest_Ctgr:
		customTag := reqMsg.GetCtgr()
		tagpath := filepath.Join(e.node.ServiceTagDir, customTag.Tag.T.Name+".tag")
		err = saveTagFile(tagpath, customTag.Tag)
		if err != nil {
			os.Remove(tagpath)
		} else {
			respMsg.Code = 0
			e.node.PutServiceTagEventCh(tagpath)
		}
	case *pb.TagPushRequest_Itgr:
		idleTag := reqMsg.GetCtgr()
		tagpath := filepath.Join(e.node.IdleTagDir, idleTag.Tag.T.Name+".tag")
		err = saveTagFile(tagpath, idleTag.Tag)
		if err != nil {
			log.Println("file req save tag err:", err)
			os.Remove(tagpath)
		} else {
			respMsg.Code = 0
			e.node.PutIdleTagEventCh(tagpath)
		}
	case *pb.TagPushRequest_Error:
		log.Println("receive file req err")
	default:
		log.Println("receive invalid file req")
	}
	w.WriteMsg(respMsg)
	log.Printf("%s: push tag response to %s sent.", s.Conn().LocalPeer().String(), remotePeer)
}

func saveTagFile(tagpath string, tag *pb.Tag) error {
	f, err := os.OpenFile(tagpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return f.Sync()
}
