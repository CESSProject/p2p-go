/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const PushTag_Protocol = "/kldr/tagpush/1"

type PushTagProtocol struct {
	*Node
}

func (n *Node) NewPushTagProtocol() *PushTagProtocol {
	e := PushTagProtocol{Node: n}
	n.SetStreamHandler(PushTag_Protocol, e.onPushTagRequest)
	return &e
}

// remote peer requests handler
func (e *protocols) TagPushReq(peerid peer.ID) (uint32, error) {
	log.Printf("Sending TagPushReq req to: %s", peerid)

	s, err := e.PushTagProtocol.NewStream(context.Background(), peerid, PushTag_Protocol)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.IdleTagGenResult{}

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

	log.Printf("TagPushReq resp code: %d", respMsg.Code)
	return respMsg.Code, nil
}

// remote peer requests handler
func (e *PushTagProtocol) onPushTagRequest(s network.Stream) {
	defer s.Close()
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
	if e.PushTagProtocol.GetIdleFileTee() != string(remotePeer) &&
		e.PushTagProtocol.GetServiceFileTee() != string(remotePeer) {
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
		tagpath := filepath.Join(e.PushTagProtocol.ServiceTagDir, customTag.Tag.T.Name+".tag")
		err = saveTagFile(tagpath, customTag.Tag)
		if err != nil {
			os.Remove(tagpath)
		} else {
			respMsg.Code = 0
			e.PushTagProtocol.PutServiceTagEventCh(tagpath)
		}
	case *pb.TagPushRequest_Itgr:
		idleTag := reqMsg.GetItgr()
		tagpath := filepath.Join(e.PushTagProtocol.IdleTagDir, idleTag.Tag.T.Name+".tag")
		err = saveTagFile(tagpath, idleTag.Tag)
		if err != nil {
			log.Println("file req save tag err:", err)
			os.Remove(tagpath)
		} else {
			respMsg.Code = 0
			e.PushTagProtocol.PutIdleTagEventCh(tagpath)
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
