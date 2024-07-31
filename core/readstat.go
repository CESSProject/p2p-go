/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/pkg/errors"

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

func (e *protocols) ReadDataStatAction(id peer.ID, roothash string, datahash string) (uint64, error) {
	var ok bool
	var err error
	var req pb.ReadDataStatRequest
	req.Roothash = roothash
	req.Datahash = datahash
	req.MessageData = e.ReadDataStatProtocol.NewMessageData(uuid.New().String(), false)

	// store request so response handler has access to it
	var respChan = make(chan bool, 1)
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
	timeout := time.NewTicker(P2PReadReqRespTime)
	defer timeout.Stop()

	err = e.ReadDataStatProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+readDataStatRequest), &req)
	if err != nil {
		return 0, errors.Wrapf(err, "[SendProtoMessage]")
	}

	//
	timeout.Reset(P2PReadReqRespTime)
	select {
	case ok = <-respChan:
		if !ok {
			return 0, errors.New(ERR_RespFailure)
		}
	case <-timeout.C:
		return 0, errors.New(ERR_TimeOut)
	}

	e.ReadDataStatProtocol.Lock()
	resp, ok := e.ReadDataStatProtocol.requests[req.MessageData.Id]
	if !ok {
		e.ReadDataStatProtocol.Unlock()
		return 0, errors.New(ERR_RespFailure)
	}
	e.ReadDataStatProtocol.Unlock()

	if resp.ReadDataStatResponse == nil {
		return 0, errors.New(ERR_RespFailure)
	}

	if resp.Code != P2PResponseOK {
		return 0, errors.New(ERR_RespFailure)
	}
	return uint64(resp.DataSize), nil
}

// remote peer requests handler
func (e *ReadDataProtocol) onReadDataStatRequest(s network.Stream) {
	defer s.Close()
	var resp = &pb.ReadDataStatResponse{
		Code: P2PResponseFailed,
	}
	// get request data
	var data = &pb.ReadDataStatRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
		return
	}
	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
		return
	}

	fpath := FindFile(e.ReadDataProtocol.GetDirs().FileDir, data.Datahash)
	if fpath == "" {
		fpath = FindFile(e.ReadDataProtocol.GetDirs().TmpDir, data.Datahash)
	}

	fstat, err := os.Stat(fpath)
	if err != nil {
		resp.Code = P2PResponseEmpty
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
		return
	}

	// send response to the request using the message string he provided
	resp.Code = P2PResponseOK
	resp.DataHash = data.Datahash
	resp.DataSize = fstat.Size()
	resp.MessageData = e.ReadDataStatProtocol.NewMessageData(data.MessageData.Id, false)
	e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataStatResponse), resp)
}

// remote peer requests handler
func (e *ReadDataStatProtocol) onReadDataStatResponse(s network.Stream) {
	defer s.Close()
	var data = &pb.ReadDataStatResponse{}
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

	if data.MessageData == nil {
		s.Reset()
		return
	}

	e.ReadDataStatProtocol.Lock()
	defer e.ReadDataStatProtocol.Unlock()

	// locate request data and remove it if found
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		if data.Code == P2PResponseOK {
			e.requests[data.MessageData.Id].ch <- true
			e.requests[data.MessageData.Id].ReadDataStatResponse = data
		} else {
			e.requests[data.MessageData.Id].ch <- false
		}
	}
}
