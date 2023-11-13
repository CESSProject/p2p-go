/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CESSProject/p2p-go/pb"

	"github.com/AstaFrode/go-libp2p/core/network"
	"github.com/AstaFrode/go-libp2p/core/peer"
	"github.com/AstaFrode/go-libp2p/core/protocol"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// pattern: /protocol-name/request-or-response-message/version
const writeFileRequest = "/file/writereq/v0"
const writeFileResponse = "/file/writeresp/v0"

type writeMsgResp struct {
	ch chan bool
	*pb.WritefileResponse
}

type WriteFileProtocol struct { // local host
	*Node
	*sync.Mutex
	requests map[string]*writeMsgResp // determine whether it is your own response
}

func (n *Node) NewWriteFileProtocol() *WriteFileProtocol {
	e := WriteFileProtocol{Node: n, Mutex: new(sync.Mutex), requests: make(map[string]*writeMsgResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+writeFileRequest), e.onWriteFileRequest)
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+writeFileResponse), e.onWriteFileResponse)
	return &e
}

func (e *protocols) WriteFileAction(id peer.ID, roothash, path string) error {
	var err error
	var ok bool
	var num int
	var offset int64
	var f *os.File
	// create message data
	req := &pb.WritefileRequest{
		MessageData: e.WriteFileProtocol.NewMessageData(uuid.New().String(), false),
		Roothash:    roothash,
	}

	fstat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fstat.Size() != FragmentSize {
		return errors.New("invalid fragment size")
	}

	req.Datahash, err = CalcPathSHA256(path)
	if err != nil {
		return err
	}

	f, err = os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// store request so response handler has access to it
	respChan := make(chan bool, 1)

	e.WriteFileProtocol.Lock()
	for {
		if _, ok := e.WriteFileProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		e.WriteFileProtocol.requests[req.MessageData.Id] = &writeMsgResp{
			ch: respChan,
		}
		break
	}
	e.WriteFileProtocol.Unlock()

	defer func() {
		e.WriteFileProtocol.Lock()
		delete(e.WriteFileProtocol.requests, req.MessageData.Id)
		close(respChan)
		e.WriteFileProtocol.Unlock()
	}()

	timeout := time.NewTicker(P2PWriteReqRespTime)
	defer timeout.Stop()
	buf := make([]byte, FileProtocolBufSize)
	for {
		_, err = f.Seek(offset, 0)
		if err != nil {
			return err
		}
		num, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if num == 0 {
			break
		}

		req.Data = buf[:num]
		req.Length = uint32(num)
		req.Offset = offset
		req.MessageData.Timestamp = time.Now().Unix()

		err = e.WriteFileProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+writeFileRequest), req)
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

		e.WriteFileProtocol.Lock()
		resp, ok := e.WriteFileProtocol.requests[req.MessageData.Id]
		if !ok {
			e.WriteFileProtocol.Unlock()
			return errors.New(ERR_RespFailure)
		}
		e.WriteFileProtocol.Unlock()

		if resp.WritefileResponse == nil {
			return errors.New(ERR_RespFailure)
		}

		if resp.WritefileResponse.Code == P2PResponseFinish {
			return nil
		}

		if resp.WritefileResponse.Code == P2PResponseOK {
			offset = resp.Offset
		} else {
			return errors.New(ERR_RespFailure)
		}
	}
	return errors.New(ERR_RespInvalidData)
}

// remote peer requests handler
func (e *WriteFileProtocol) onWriteFileRequest(s network.Stream) {
	defer s.Close()
	// get request data
	data := &pb.WritefileRequest{}
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

	resp := &pb.WritefileResponse{
		MessageData: e.NewMessageData(data.MessageData.Id, false),
		Code:        P2PResponseOK,
		Offset:      0,
	}

	dir := filepath.Join(e.GetDirs().TmpDir, data.Roothash)
	fstat, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, DirMode)
		if err != nil {
			s.Reset()
			return
		}
	} else {
		if !fstat.IsDir() {
			os.Remove(dir)
			err = os.MkdirAll(dir, DirMode)
			if err != nil {
				s.Reset()
				return
			}
		}
	}
	var size int64

	fpath := filepath.Join(dir, data.Datahash)

	fstat, err = os.Stat(fpath)
	if err == nil {
		size = fstat.Size()
		if size >= FragmentSize {
			time.Sleep(time.Second * 5)
			if size > FragmentSize {
				os.Remove(fpath)
			} else {
				hash, err := CalcPathSHA256(fpath)
				if err != nil || hash != data.Datahash {
					os.Remove(fpath)
				} else {
					resp.Code = P2PResponseFinish
				}
			}

			e.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeFileResponse), resp)
			return
		}
	}

	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		s.Reset()
		return
	}
	defer f.Close()
	fstat, err = f.Stat()
	if err != nil {
		s.Reset()
		return
	}
	size = fstat.Size()

	_, err = f.Write(data.Data[:data.Length])
	if err != nil {
		s.Reset()
		return
	}

	if int(int(size)+int(data.Length)) == FragmentSize {
		hash, err := CalcPathSHA256(fpath)
		if err != nil || hash != data.Datahash {
			os.Remove(fpath)
		} else {
			resp.Code = P2PResponseFinish
		}
	} else {
		resp.Offset = size + int64(data.Length)
	}
	// send response to the request using the message string he provided
	e.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+writeFileResponse), resp)
}

// remote peer response handler
func (e *WriteFileProtocol) onWriteFileResponse(s network.Stream) {
	defer s.Close()

	data := &pb.WritefileResponse{}
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

	// locate request data and remove it if found
	e.WriteFileProtocol.Lock()
	defer e.WriteFileProtocol.Unlock()

	_, ok := e.requests[data.MessageData.Id]
	if ok {
		if data.Code == P2PResponseOK || data.Code == P2PResponseFinish {
			e.requests[data.MessageData.Id].ch <- true
			e.requests[data.MessageData.Id].WritefileResponse = data
		} else {
			e.requests[data.MessageData.Id].ch <- false
		}
	}
}
