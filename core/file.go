/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"os"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/pkg/errors"
)

const FILE_PROTOCOL = "/kft/1"

type FileProtocol struct {
	*Node
}

func (n *Node) NewFileProtocol() *FileProtocol {
	e := FileProtocol{Node: n}
	return &e
}

func (e *protocols) FileReq(peerId peer.ID, filehash string, filetype pb.FileType, fpath string) (uint32, error) {
	fstat, err := os.Stat(fpath)
	if err != nil {
		return 0, err
	}
	var req = &pb.Request{}
	var putReq = &pb.PutRequest{
		Type: filetype,
		Hash: filehash,
		Size: uint64(fstat.Size()),
	}

	s, err := e.FileProtocol.NewStream(context.Background(), peerId, protocol.ID(e.ProtocolPrefix+FILE_PROTOCOL))
	if err != nil {
		return 0, errors.Wrapf(err, "[NewStream]")
	}
	defer s.Close()

	putReq.Data, err = os.ReadFile(fpath)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[ReadFile]")
	}

	var reqMsg = &pb.Request_PutRequest{
		PutRequest: putReq,
	}
	reqMsg.PutRequest = putReq
	req.Request = reqMsg

	w := pbio.NewDelimitedWriter(s)
	err = w.WriteMsg(req)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[WriteMsg]")
	}

	respMsg := &pb.PutResponse{}

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[ReadMsg]")
	}

	if respMsg.Code == 0 {
		e.FileProtocol.SetServiceFileTee(peerId.String())
	}

	return respMsg.Code, nil
}
