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
	"path/filepath"
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
const readDataRequest = "/data/readreq/v0"
const readDataResponse = "/data/readresp/v0"

type readDataResp struct {
	ch chan bool
	*pb.ReadDataResponse
}

type ReadDataProtocol struct {
	*PeerNode // local host
	*sync.Mutex
	requests map[string]*readDataResp // determine whether it is your own response
}

func (n *PeerNode) NewReadDataProtocol() *ReadDataProtocol {
	e := ReadDataProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*readDataResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readDataRequest), e.onReadDataRequest)
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readDataResponse), e.onReadDataResponse)
	return &e
}

func (e *protocols) ReadDataAction(id peer.ID, name, savepath string) (int64, error) {
	fstat, err := os.Stat(savepath)
	if err == nil {
		if fstat.IsDir() {
			return 0, fmt.Errorf("%s is a directory", savepath)
		}
		if fstat.Size() > 0 {
			return fstat.Size(), nil
		}
		os.Remove(savepath)
	}

	req := pb.ReadDataRequest{
		MessageData: e.ReadDataProtocol.NewMessageData(uuid.New().String(), false),
		Name:        name,
	}
	respChan := make(chan bool, 1)
	e.ReadDataProtocol.Lock()
	for {
		if _, ok := e.ReadDataProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		e.ReadDataProtocol.requests[req.MessageData.Id] = &readDataResp{
			ch: respChan,
		}
		break
	}
	e.ReadDataProtocol.Unlock()

	defer func() {
		e.ReadDataProtocol.Lock()
		delete(e.ReadDataProtocol.requests, req.MessageData.Id)
		close(respChan)
		e.ReadDataProtocol.Unlock()
	}()

	err = e.ReadDataProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+readDataRequest), &req)
	if err != nil {
		return 0, fmt.Errorf("SendProtoMessage: %v", err)
	}
	timeout := time.After(time.Second * 15)
	for {
		select {
		case <-timeout:
			return 0, fmt.Errorf(ERR_RecvTimeOut)
		case <-respChan:
			e.ReadDataProtocol.Lock()
			resp, ok := e.ReadDataProtocol.requests[req.MessageData.Id]
			e.ReadDataProtocol.Unlock()
			if !ok {
				return 0, fmt.Errorf("received empty data")
			}
			if resp.Code != P2PResponseOK {
				return 0, fmt.Errorf("received code: %d err: %v", resp.Code, resp.Msg)
			}
			if len(resp.Data) != int(resp.DataLength) {
				return 0, fmt.Errorf("received code: %d and len(data)=%d != datalength=%d", resp.Code, len(resp.Data), resp.DataLength)
			}
			err = os.WriteFile(savepath, resp.Data, 0755)
			if err != nil {
				return 0, fmt.Errorf("received suc, local WriteFile failed: %v", err)
			}
			return resp.DataLength, nil
		}
	}
}

func (e *ReadDataProtocol) onReadDataRequest(s network.Stream) {
	defer s.Close()

	// get request data
	data := &pb.ReadDataRequest{}
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

	resp := &pb.ReadDataResponse{
		MessageData: e.NewMessageData(data.MessageData.Id, false),
	}

	fpath := FindFile(e.ReadDataProtocol.GetDirs().FileDir, data.Name)
	if fpath == "" {
		fpath = filepath.Join(e.ReadDataProtocol.GetDirs().TmpDir, data.Name)
	}

	fstat, err := os.Stat(fpath)
	if err != nil {
		resp.Code = P2PResponseEmpty
		resp.Msg = "not fount"
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataResponse), resp)
		return
	}

	if fstat.Size() < FragmentSize {
		resp.Code = P2PResponseForbid
		resp.Msg = "i am also receiving the file, please read it later."
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataResponse), resp)
		return
	}

	buffer, err := os.ReadFile(fpath)
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("os.ReadFile: %v", err)
		e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataResponse), resp)
		return
	}

	resp.Code = P2PResponseOK
	resp.Data = buffer
	resp.DataLength = int64(len(buffer))
	resp.Msg = "success"
	e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataResponse), resp)
	return
}

func (e *ReadDataProtocol) onReadDataResponse(s network.Stream) {
	defer s.Close()
	data := &pb.ReadDataResponse{}
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

	// locate request data and remove it if found
	e.ReadDataProtocol.Lock()
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		e.requests[data.MessageData.Id].ch <- true
		e.requests[data.MessageData.Id].ReadDataResponse = data
	}
	e.ReadDataProtocol.Unlock()
}
