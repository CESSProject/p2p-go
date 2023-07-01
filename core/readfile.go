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

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// pattern: /protocol-name/request-or-response-message/version
const readFileRequest = "/file/readreq/v0"
const readFileResponse = "/file/readresp/v0"

type readMsgResp struct {
	ch chan bool
	*pb.ReadfileResponse
}

type ReadFileProtocol struct {
	*Node // local host
	*sync.Mutex
	requests map[string]*readMsgResp // determine whether it is your own response
}

func (n *Node) NewReadFileProtocol() *ReadFileProtocol {
	e := ReadFileProtocol{Node: n, Mutex: new(sync.Mutex), requests: make(map[string]*readMsgResp)}
	n.SetStreamHandler(readFileRequest, e.onReadFileRequest)
	n.SetStreamHandler(readFileResponse, e.onReadFileResponse)
	return &e
}

func (e *protocols) ReadFileAction(id peer.ID, roothash, datahash, path string, size int64) error {
	var ok bool
	var err error
	var hash string
	var offset int64
	var num int
	var fstat fs.FileInfo
	var f *os.File
	var req pb.ReadfileRequest
	var resp *pb.ReadfileResponse

	fstat, err = os.Stat(path)
	if err == nil {
		if fstat.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}
		if fstat.Size() < size {
			offset = fstat.Size()
		} else if fstat.Size() == size {
			hash, err = CalcPathSHA256(path)
			if err != nil {
				return errors.Wrapf(err, "[calc file sha256]")
			}
			if hash == datahash {
				return nil
			}
			return fmt.Errorf("fragment hash does not match file")
		} else {
			buf, err := os.ReadFile(path)
			if err != nil {
				return errors.Wrapf(err, "[read file]")
			}
			hash, err = CalcSHA256(buf[:size])
			if err != nil {
				return errors.Wrapf(err, "[calc buf sha256]")
			}
			if hash == datahash {
				f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
				if err != nil {
					return errors.Wrapf(err, "[open file]")
				}
				defer f.Close()
				_, err = f.Write(buf[:size])
				if err != nil {
					return errors.Wrapf(err, "[write file]")
				}
				err = f.Sync()
				if err != nil {
					return errors.Wrapf(err, "[sync file]")
				}
				return nil
			}
			os.Remove(path)
		}
	}
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "[open file]")
	}
	defer f.Close()

	req.Roothash = roothash
	req.Datahash = datahash
	req.MessageData = e.ReadFileProtocol.NewMessageData(uuid.New().String(), false)

	// store request so response handler has access to it
	var respChan = make(chan bool, 1)
	e.ReadFileProtocol.Lock()
	e.ReadFileProtocol.requests[req.MessageData.Id] = &readMsgResp{
		ch: respChan,
	}
	e.ReadFileProtocol.Unlock()
	defer func() {
		e.ReadFileProtocol.Lock()
		delete(e.ReadFileProtocol.requests, req.MessageData.Id)
		close(respChan)
		e.ReadFileProtocol.Unlock()
	}()
	timeout := time.NewTicker(P2PReadReqRespTime)
	defer timeout.Stop()

	for {
		req.Offset = offset

		err = e.ReadFileProtocol.SendProtoMessage(id, readFileRequest, &req)
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

		resp = e.ReadFileProtocol.requests[req.MessageData.Id].ReadfileResponse
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

	return errors.New(ERR_RespInvalidData)
}

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileRequest(s network.Stream) {
	defer s.Close()

	var code = P2PResponseOK
	// get request data
	data := &pb.ReadfileRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		return
	}

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		return
	}

	fpath := filepath.Join(e.ReadFileProtocol.GetDirs().TmpDir, data.Roothash, data.Datahash)

	_, err = os.Stat(fpath)
	if err != nil {
		fpath = filepath.Join(e.ReadFileProtocol.GetDirs().FileDir, data.Roothash, data.Datahash)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer f.Close()

	fstat, err := f.Stat()
	if err != nil {
		return
	}

	_, err = f.Seek(data.Offset, 0)
	if err != nil {
		return
	}

	var readBuf = make([]byte, FileProtocolBufSize)
	num, err := f.Read(readBuf)
	if err != nil {
		return
	}

	if num+int(data.Offset) >= int(fstat.Size()) {
		code = P2PResponseFinish
	}

	// send response to the request using the message string he provided
	resp := &pb.ReadfileResponse{
		MessageData: e.ReadFileProtocol.NewMessageData(data.MessageData.Id, false),
		Code:        code,
		Offset:      data.Offset,
		Length:      uint32(num),
		Data:        readBuf[:num],
	}

	e.ReadFileProtocol.SendProtoMessage(s.Conn().RemotePeer(), readFileResponse, resp)
}

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileResponse(s network.Stream) {
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

	e.ReadFileProtocol.Lock()
	defer e.ReadFileProtocol.Unlock()
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
