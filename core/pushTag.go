/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/pkg/errors"
)

const PushTag_Protocol = "/tagpush/1"

type PushTagProtocol struct {
	*Node
}

func (n *Node) NewPushTagProtocol() *PushTagProtocol {
	e := PushTagProtocol{Node: n}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+PushTag_Protocol), e.onPushTagRequest)
	return &e
}

// remote peer requests handler
func (e *protocols) TagPushReq(peerid peer.ID) (uint32, error) {
	s, err := e.PushTagProtocol.NewStream(context.Background(), peerid, protocol.ID(e.ProtocolPrefix+PushTag_Protocol))
	if err != nil {
		return 0, errors.Wrapf(err, "[NewStream]")
	}
	defer func() {
		s.Reset()
		time.Sleep(time.Millisecond * 10)
		s.Close()
	}()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.IdleTagGenResult{}

	err = w.WriteMsg(reqMsg)
	if err != nil {
		return 0, errors.Wrapf(err, "[WriteMsg]")
	}

	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	respMsg := &pb.AggrProofResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		return 0, errors.Wrapf(err, "[ReadMsg]")
	}
	return respMsg.Code, nil
}

// remote peer requests handler
func (e *PushTagProtocol) onPushTagRequest(s network.Stream) {
	defer s.Close()

	respMsg := &pb.TagPushResponse{
		Code: 1,
	}
	reqMsg := &pb.TagPushRequest{}
	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)

	err := r.ReadMsg(reqMsg)
	if err != nil {
		w.WriteMsg(respMsg)
		return
	}
	remotePeer := s.Conn().RemotePeer().String()

	if e.PushTagProtocol.GetIdleFileTee() != string(remotePeer) &&
		e.PushTagProtocol.GetServiceFileTee() != string(remotePeer) {
		w.WriteMsg(respMsg)
		return
	}

	switch reqMsg.GetResult().(type) {
	case *pb.TagPushRequest_Ctgr:
		customTag := reqMsg.GetCtgr()
		tagpath := filepath.Join(e.PushTagProtocol.GetDirs().ServiceTagDir, customTag.Tag.T.Name+".tag")
		err = saveTagFile(tagpath, customTag.Tag)
		if err != nil {
			os.Remove(tagpath)
		} else {
			respMsg.Code = 0
			e.PushTagProtocol.putServiceTagCh(tagpath)
		}
	case *pb.TagPushRequest_Error:
	default:
	}
	w.WriteMsg(respMsg)
	return
}

func saveTagFile(tagpath string, tag *pb.Tag2) error {
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
