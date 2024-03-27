/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
const readDataRequest = "/data/readreq/v0"
const readDataResponse = "/data/readresp/v0"

type readDataResp struct {
	ch chan bool
	*pb.ReadfileResponse
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

func (e *protocols) ReadDataAction(id peer.ID, roothash, datahash, path string, size int64) error {
	var ok bool
	var err error
	var offset int64
	var num int
	var fstat fs.FileInfo
	var f *os.File
	var req pb.ReadfileRequest

	if size <= 0 {
		return errors.New("invalid size")
	}

	fstat, err = os.Stat(path)
	if err == nil {
		if fstat.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}
		if fstat.Size() < size {
			offset = fstat.Size()
		} else if fstat.Size() == size {
			return nil
		} else {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				if f != nil {
					f.Close()
				}
			}()
			newpath := filepath.Join(filepath.Dir(path), uuid.New().String())
			new_f, err := os.Create(newpath)
			if err != nil {
				return err
			}
			defer func() {
				if new_f != nil {
					new_f.Close()
				}
			}()

			_, err = io.CopyN(new_f, f, size)
			if err != nil {
				return err
			}
			f.Close()
			f = nil
			new_f.Close()
			new_f = nil
			return nil
		}
	}
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "[open file]")
	}
	defer f.Close()

	req.Roothash = roothash
	req.Datahash = datahash
	req.MessageData = e.ReadDataProtocol.NewMessageData(uuid.New().String(), false)

	// store request so response handler has access to it
	var respChan = make(chan bool, 1)
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
	timeout := time.NewTicker(P2PReadReqRespTime)
	defer timeout.Stop()

	for {
		req.Offset = offset

		err = e.ReadDataProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+readDataRequest), &req)
		if err != nil {
			return errors.Wrapf(err, "[SendProtoMessage]")
		}

		//
		timeout.Reset(P2PReadReqRespTime)
		select {
		case ok = <-respChan:
			if !ok {
				return errors.New(ERR_RespFailure)
			}
		case <-timeout.C:
			return errors.New(ERR_RespTimeOut)
		}

		e.ReadDataProtocol.Lock()
		resp, ok := e.ReadDataProtocol.requests[req.MessageData.Id]
		if !ok {
			e.ReadDataProtocol.Unlock()
			return errors.New(ERR_RespFailure)
		}
		e.ReadDataProtocol.Unlock()

		if resp.ReadfileResponse == nil {
			return errors.New(ERR_RespFailure)
		}

		if len(resp.ReadfileResponse.Data) == 0 || resp.Length == 0 {
			if resp.Code == P2PResponseFinish {
				err = f.Sync()
				if err != nil {
					return errors.Wrapf(err, "[sync file]")
				}
				return nil
			}
			return errors.New(ERR_RespFailure)
		}

		num, err = f.Write(resp.Data[:resp.Length])
		if err != nil {
			return errors.Wrapf(err, "[write file]")
		}

		if resp.Code == P2PResponseFinish {
			err = f.Sync()
			if err != nil {
				return errors.Wrapf(err, "[sync file]")
			}
			return nil
		}

		offset = req.Offset + int64(num)
	}
}

// remote peer requests handler
func (e *ReadDataProtocol) onReadDataRequest(s network.Stream) {
	defer s.Close()

	var code = P2PResponseOK
	// get request data
	data := &pb.ReadfileRequest{}
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

	fpath := FindFile(e.ReadDataProtocol.GetDirs().FileDir, data.Datahash)
	if fpath == "" {
		fpath = filepath.Join(e.ReadDataProtocol.GetDirs().TmpDir, data.Datahash)
	}

	_, err = os.Stat(fpath)
	if err != nil {
		s.Reset()
		return
	}

	f, err := os.Open(fpath)
	if err != nil {
		s.Reset()
		return
	}
	defer f.Close()

	fstat, err := f.Stat()
	if err != nil {
		s.Reset()
		return
	}

	_, err = f.Seek(data.Offset, 0)
	if err != nil {
		s.Reset()
		return
	}

	var readBuf = make([]byte, FileProtocolBufSize)
	num, err := f.Read(readBuf)
	if err != nil {
		s.Reset()
		return
	}

	if num+int(data.Offset) >= int(fstat.Size()) {
		code = P2PResponseFinish
	}

	// send response to the request using the message string he provided
	resp := &pb.ReadfileResponse{
		MessageData: e.ReadDataProtocol.NewMessageData(data.MessageData.Id, false),
		Code:        code,
		Offset:      data.Offset,
		Length:      uint32(num),
		Data:        readBuf[:num],
	}

	e.ReadDataProtocol.SendProtoMessage(s.Conn().RemotePeer(), protocol.ID(e.ProtocolPrefix+readDataResponse), resp)
}

// remote peer requests handler
func (e *ReadDataProtocol) onReadDataResponse(s network.Stream) {
	defer s.Close()
	data := &pb.ReadfileResponse{}
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

	e.ReadDataProtocol.Lock()
	defer e.ReadDataProtocol.Unlock()

	// locate request data and remove it if found
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		if data.Code == P2PResponseOK || data.Code == P2PResponseFinish {
			e.requests[data.MessageData.Id].ch <- true
			e.requests[data.MessageData.Id].ReadfileResponse = data
		} else {
			e.requests[data.MessageData.Id].ch <- false
		}
	}
}
