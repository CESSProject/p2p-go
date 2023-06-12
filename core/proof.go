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
)

const AggrProof_PROTOCOL = "/kldr/apv/1"

type AggrProofProtocol struct {
	*Node
}

func (n *Node) NewAggrProofProtocol() *AggrProofProtocol {
	e := AggrProofProtocol{Node: n}
	return &e
}

func (e *protocols) AggrProofReq(peerId peer.ID, ihash, shash []byte, qslice []*pb.Qslice, puk, sign []byte) (uint32, error) {
	s, err := e.AggrProofProtocol.NewStream(context.Background(), peerId, AggrProof_PROTOCOL)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.AggrProofRequest{
		IdleProofFileHash:    ihash,
		ServiceProofFileHash: shash,
		Publickey:            puk,
		Sign:                 sign,
		Qslice:               qslice,
	}

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
	return respMsg.Code, nil
}
