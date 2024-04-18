/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/CESSProject/p2p-go/pb"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"google.golang.org/protobuf/proto"
)

// pattern: /protocol-name/request-or-response-message/version
const OnlineRequest = "/online/req/v0"
const OnlineResponse = "/online/resp/v0"

type onlineResp struct {
	ch chan bool
	*pb.MessageData
}

type OnlineProtocol struct { // local host
	*PeerNode
	*sync.Mutex
	requests map[string]*onlineResp // determine whether it is your own response
}

func (n *PeerNode) NewOnlineProtocol() *OnlineProtocol {
	e := OnlineProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*onlineResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+OnlineResponse), e.onOnlineResponse)
	return &e
}

func (e *protocols) OnlineAction(id peer.ID) error {
	var err error
	var ok bool
	// create message data
	req := &pb.MessageData{
		Id:     uuid.New().String(),
		NodeId: e.OnlineProtocol.ID().String(),
	}

	// store request so response handler has access to it
	respChan := make(chan bool, 1)

	e.OnlineProtocol.Lock()
	for {
		if _, ok := e.OnlineProtocol.requests[req.Id]; ok {
			req.Id = uuid.New().String()
			continue
		}
		e.OnlineProtocol.requests[req.Id] = &onlineResp{
			ch:          respChan,
			MessageData: &pb.MessageData{},
		}
		break
	}
	e.OnlineProtocol.Unlock()

	defer func() {
		e.OnlineProtocol.Lock()
		delete(e.OnlineProtocol.requests, req.Id)
		close(respChan)
		e.OnlineProtocol.Unlock()
	}()

	timeout := time.NewTicker(P2PWriteReqRespTime)
	defer timeout.Stop()

	err = e.OnlineProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+OnlineRequest), req)
	if err != nil {
		return err
	}

	// wait response
	timeout.Reset(P2PWriteReqRespTime)
	select {
	case ok = <-respChan:
		if !ok {
			return errors.New(ERR_RespFailure)
		}
	case <-timeout.C:
		return errors.New(ERR_RespTimeOut)
	}

	e.OnlineProtocol.Lock()
	resp, ok := e.OnlineProtocol.requests[req.Id]
	if !ok {
		e.OnlineProtocol.Unlock()
		return errors.New(ERR_RespFailure)
	}
	e.OnlineProtocol.Unlock()
	fmt.Println("received onlline response msg from: ", resp.NodeId)
	if resp.NodeId == id.String() {
		return nil
	}

	return errors.New(ERR_RespInvalidData)
}

// remote peer requests handler
func (e *OnlineProtocol) onOnlineRequest(s network.Stream) {
	defer s.Close()

	if !e.enableRecv {
		s.Reset()
		return
	}

	// get request data
	data := &pb.MessageData{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		return
	}

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		s.Reset()
		return
	}

	fmt.Println("i already receiced onlline msg from ", s.Conn().RemotePeer())
	fmt.Println("msg multiaddr is: ", data.NodeId)
	resp := &pb.MessageData{
		Id:     data.Id,
		NodeId: e.OnlineProtocol.ID().String(),
	}
	fmt.Println("RemoteMultiaddr: ", s.Conn().RemoteMultiaddr())
	fmt.Println("RemoteMultiaddr.string: ", s.Conn().RemoteMultiaddr().String())
	fmt.Println("response online: ", e.OnlineProtocol.ID().String())
	e.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+OnlineResponse), resp)
}

// remote peer response handler
func (e *OnlineProtocol) onOnlineResponse(s network.Stream) {
	defer s.Close()

	data := &pb.MessageData{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		return
	}

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		s.Reset()
		return
	}

	if data.NodeId == "" {
		s.Reset()
		return
	}

	// locate request data and remove it if found
	e.OnlineProtocol.Lock()
	defer e.OnlineProtocol.Unlock()

	_, ok := e.requests[data.Id]
	if ok {
		e.requests[data.Id].ch <- true
		e.requests[data.Id].NodeId = data.NodeId
	}
}
