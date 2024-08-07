/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"fmt"
	"time"

	"github.com/CESSProject/p2p-go/pb"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// pattern: /protocol-name/request-or-response-message/version
const OnlineRequest = "/online/req/v0"

type onlineResp struct {
	*pb.MessageData
}

type OnlineProtocol struct { // local host
	*PeerNode
	record chan string
}

func (n *PeerNode) NewOnlineProtocol(cacheLen int) *OnlineProtocol {
	if cacheLen <= 0 {
		cacheLen = 1000
	}
	e := OnlineProtocol{PeerNode: n, record: make(chan string, cacheLen)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+OnlineRequest), e.onOnlineResponse)
	return &e
}

func (e *protocols) OnlineAction(id peer.ID) error {
	err := e.OnlineProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+OnlineRequest), &pb.MessageData{})
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return nil
}

// remote peer response handler
func (e *OnlineProtocol) onOnlineResponse(s network.Stream) {
	record := fmt.Sprintf("%v/p2p/%s", s.Conn().RemoteMultiaddr(), s.Conn().RemotePeer().String())
	s.Close()
	if len(e.record) < 100 {
		e.record <- record
	}
}

func (e *OnlineProtocol) GetRecord() <-chan string {
	return e.record
}
