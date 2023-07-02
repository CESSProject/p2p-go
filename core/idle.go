/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/pkg/errors"
)

const IdleDataTag_Protocol = "/kldr/idtg/1"

type IdleDataTagProtocol struct {
	*Node
}

func (n *Node) NewIdleDataTagProtocol() *IdleDataTagProtocol {
	e := IdleDataTagProtocol{Node: n}
	return &e
}

func (e *protocols) IdleReq(peerId peer.ID, filesize, blocknum uint64, pubkey, sign []byte) (uint32, error) {
	reqMsg := &pb.IdleDataTagRequest{
		FileSize:  filesize,
		BlockNum:  blocknum,
		Publickey: pubkey,
		Sign:      sign,
	}

	s, err := e.IdleDataTagProtocol.NewStream(context.Background(), peerId, IdleDataTag_Protocol)
	if err != nil {
		return 0, errors.Wrapf(err, "[NewStream]")
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[WriteMsg]")
	}

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	var respMsg = &pb.IdleDataTagResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[ReadMsg]")
	}
	if respMsg.Code == 0 {
		e.IdleDataTagProtocol.SetIdleFileTee(peerId.String())
	}
	return respMsg.Code, nil
}
