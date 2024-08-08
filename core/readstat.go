/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"fmt"
	"io"
	"os"
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
const readDataStatRequest = "/read/stat/req/v0"
const readDataStatResponse = "/read/stat/resp/v0"

type readDataStatResp struct {
	ch chan bool
	*pb.ReadDataStatResponse
}

type ReadDataStatProtocol struct {
	*PeerNode // local host
	*sync.Mutex
	requests map[string]*readDataStatResp // determine whether it is your own response
}

func (n *PeerNode) NewReadDataStatProtocol() *ReadDataStatProtocol {
	e := ReadDataStatProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*readDataStatResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readDataStatRequest), e.onReadDataStatRequest)
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readDataStatResponse), e.onReadDataStatResponse)
	return &e
}

func (e *protocols) ReadDataStatAction(id peer.ID, name string) (uint64, error) {
	req := pb.ReadDataStatRequest{
		MessageData: e.ReadDataStatProtocol.NewMessageData(uuid.New().String(), false),
		Name:        name,
	}
	respChan := make(chan bool, 1)
	e.ReadDataStatProtocol.Lock()
	for {
		if _, ok := e.ReadDataStatProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		e.ReadDataStatProtocol.requests[req.MessageData.Id] = &readDataStatResp{
			ch: respChan,
		}
		break
	}
	e.ReadDataStatProtocol.Unlock()

	defer func() {
		e.ReadDataStatProtocol.Lock()
		delete(e.ReadDataStatProtocol.requests, req.MessageData.Id)
		close(respChan)
		e.ReadDataStatProtocol.Unlock()
	}()

	err := e.ReadDataStatProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+readDataStatRequest), &req)
	if err != nil {
		return 0, fmt.Errorf("SendProtoMessage: %v", err)
	}
	timeout := time.After(time.Second * 15)
	for {
		select {
		case <-timeout:
			return 0, fmt.Errorf(ERR_RecvTimeOut)
		case <-respChan:
			e.ReadDataStatProtocol.Lock()
			resp, ok := e.ReadDataStatProtocol.requests[req.MessageData.Id]
			e.ReadDataStatProtocol.Unlock()
			if !ok {
				return 0, fmt.Errorf("received empty data")
			}
			if resp.Code != P2PResponseOK {
				return 0, fmt.Errorf("received code: %d err: %v", resp.Code, resp.Msg)
			}
			return uint64(resp.Size), nil
		}
	}
}

// remote peer requests handler
func (e *ReadDataProtocol) onReadDataStatRequest(s network.Stream) {
	defer s.Close()

	// get request data
	data := &pb.ReadDataStatRequest{}
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

	resp := &pb.ReadDataStatResponse{
		MessageData: e.NewMessageData(data.MessageData.Id, false),
	}

	fpath := FindFile(e.ReadDataProtocol.GetDirs().FileDir, data.Name)
	if fpath == "" {
		fpath = FindFile(e.ReadDataProtocol.GetDirs().TmpDir, data.Name)
	}

	fstat, err := os.Stat(fpath)
	if err != nil {
		resp.Code = P2PResponseEmpty
		resp.Msg = fmt.Sprintf("not found")
		e.ReadDataStatProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
		return
	}

	// send response to the request using the message string he provided
	resp.Code = P2PResponseOK
	resp.Size = fstat.Size()
	resp.Msg = "success"
	e.ReadDataStatProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
}

// remote peer requests handler
func (e *ReadDataStatProtocol) onReadDataStatResponse(s network.Stream) {
	defer s.Close()

	data := &pb.ReadDataStatResponse{}
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

	e.ReadDataStatProtocol.Lock()
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		e.requests[data.MessageData.Id].ch <- true
		e.requests[data.MessageData.Id].ReadDataStatResponse = data
	}
	e.ReadDataStatProtocol.Unlock()
}
