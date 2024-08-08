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
const writeDataRequest = "/data/writereq/v0"
const writeDataResponse = "/data/writeresp/v0"

type writeDataResp struct {
	ch chan bool
	*pb.WriteDataResponse
}

type WriteDataProtocol struct {
	*PeerNode // local host
	*sync.Mutex
	requests map[string]*writeDataResp // determine whether it is your own response
}

func (n *PeerNode) NewWriteDataProtocol() *WriteDataProtocol {
	e := WriteDataProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*writeDataResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+writeDataRequest), e.onWriteDataRequest)
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+writeDataResponse), e.onWriteDataResponse)
	return &e
}

func (e *protocols) WriteDataAction(id peer.ID, file, fid, fragment string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("[os.ReadFile] %v", err)
	}

	req := pb.WriteDataRequest{
		MessageData: e.WriteDataProtocol.NewMessageData(uuid.New().String(), false),
		Fid:         fid,
		Datahash:    fragment,
		Data:        buf,
		DataLength:  int64(len(buf)),
	}
	respChan := make(chan bool, 1)
	e.WriteDataProtocol.Lock()
	for {
		if _, ok := e.WriteDataProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		e.WriteDataProtocol.requests[req.MessageData.Id] = &writeDataResp{
			ch: respChan,
		}
		break
	}
	e.WriteDataProtocol.Unlock()

	defer func() {
		e.WriteDataProtocol.Lock()
		delete(e.WriteDataProtocol.requests, req.MessageData.Id)
		close(respChan)
		e.WriteDataProtocol.Unlock()
	}()

	err = e.WriteDataProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+writeDataRequest), &req)
	if err != nil {
		return err
	}
	timeout := time.After(time.Second * 15)
	for {
		select {
		case <-timeout:
			return fmt.Errorf(ERR_RecvTimeOut)
		case <-respChan:
			e.WriteDataProtocol.Lock()
			resp, ok := e.WriteDataProtocol.requests[req.MessageData.Id]
			e.WriteDataProtocol.Unlock()
			if !ok {
				return fmt.Errorf("received empty data")
			}
			if resp.Code != P2PResponseOK {
				return fmt.Errorf("received code: %d err: %v", resp.Code, resp.Msg)
			}
			return nil
		}
	}
}

func (e *WriteDataProtocol) onWriteDataRequest(s network.Stream) {
	if !e.enableRecv {
		s.Reset()
		return
	}

	defer s.Close()

	// get request data
	data := &pb.WriteDataRequest{}
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

	resp := &pb.WriteDataResponse{
		MessageData: e.NewMessageData(data.MessageData.Id, false),
	}

	dir := filepath.Join(e.GetDirs().TmpDir, data.Fid)
	err = os.MkdirAll(dir, DirMode)
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("os.MkdirAll: %v", err)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	hash, err := CalcSHA256(data.Data)
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("CalcSHA256(len=%d): %v", len(data.Data), err)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	if hash != data.Datahash {
		resp.Code = P2PResponseFailed
		resp.Msg = fmt.Sprintf("received data hash: %s != recalculated hash: %s", data.Datahash, hash)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	fpath := filepath.Join(dir, data.Datahash)
	fstat, err := os.Stat(fpath)
	if err == nil {
		if fstat.Size() == data.DataLength {
			resp.Code = P2PResponseOK
			resp.Msg = "success"
			e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
			return
		}
		os.Remove(fpath)
	}

	if data.Datahash == ZeroFileHash_8M {
		err = writeZeroToFile(fpath, FragmentSize)
		if err != nil {
			resp.Code = P2PResponseRemoteFailed
			resp.Msg = fmt.Sprintf("writeZeroToFile: %v", err)
			e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
			return
		}
		resp.Code = P2PResponseOK
		resp.Msg = "success"
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	f, err := os.Create(fpath)
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("os.Create: %v", err)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}
	defer f.Close()

	_, err = f.Write(data.Data)
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("f.Write: %v", err)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	err = f.Sync()
	if err != nil {
		resp.Code = P2PResponseRemoteFailed
		resp.Msg = fmt.Sprintf("f.Sync: %v", err)
		e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
		return
	}

	resp.Code = P2PResponseOK
	resp.Msg = "success"
	e.WriteDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeDataResponse), resp)
	return
}

func (e *WriteDataProtocol) onWriteDataResponse(s network.Stream) {
	defer s.Close()
	data := &pb.WriteDataResponse{}
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
	e.WriteDataProtocol.Lock()
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		e.requests[data.MessageData.Id].ch <- true
		e.requests[data.MessageData.Id].WriteDataResponse = data
	}
	e.WriteDataProtocol.Unlock()
}

func writeZeroToFile(file string, size int64) error {
	fstat, err := os.Stat(file)
	if err == nil {
		if fstat.Size() == size {
			return nil
		}
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(make([]byte, size))
	if err != nil {
		return err
	}
	return f.Sync()
}
