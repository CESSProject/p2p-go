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
	remotePeer := s.Conn().RemotePeer()
	tagpath := ""
	log.Printf("receive push tag req: %s", remotePeer)

	if e.node.GetIdleFileTee() == string(remotePeer) {
		tagpath = filepath.Join(e.node.IdleTagDir, reqMsg.Tag.T.Name+".tag")
	} else if e.node.GetServiceFileTee() == string(remotePeer) {
		tagpath = filepath.Join(e.node.ServiceTagDir, reqMsg.Tag.T.Name+".tag")
	} else {
		s.Reset()
		log.Printf("receive invalid push tag req: %s", remotePeer)
		return
	}

	respMsg := &pb.TagPushResponse{
		Code: 0,
	}

	w := pbio.NewDelimitedWriter(s)

	f, err := os.OpenFile(tagpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		respMsg.Code = 1
		log.Println(err)
		w.WriteMsg(respMsg)
		return
	}
	defer f.Close()

	b, err := json.Marshal(reqMsg.Tag)
	if err != nil {
		os.Remove(tagpath)
		respMsg.Code = 1
		log.Println(err)
		w.WriteMsg(respMsg)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		os.Remove(tagpath)
		respMsg.Code = 1
		log.Println(err)
		w.WriteMsg(respMsg)
		return
	}
	err = f.Sync()
	if err != nil {
		os.Remove(tagpath)
		respMsg.Code = 1
		log.Println(err)
		w.WriteMsg(respMsg)
		return
	}
	e.node.PutTagEventCh(tagpath)
	w.WriteMsg(respMsg)
	log.Printf("%s: push tag response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
}
