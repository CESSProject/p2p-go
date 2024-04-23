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
const readFileRequest = "/file/readreq/v0"
const readFileResponse = "/file/readresp/v0"

type ReadFileServerHandle func(req *pb.ReadfileRequest) (*pb.ReadfileResponse, error)

type readMsgResp struct {
	ch chan bool
	*pb.ReadfileResponse
}

type ReadFileProtocol struct {
	*PeerNode // local host
	*sync.Mutex
	handle   ReadFileServerHandle
	requests map[string]*readMsgResp // determine whether it is your own response
}

func (n *PeerNode) NewReadFileProtocol() *ReadFileProtocol {
	e := ReadFileProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*readMsgResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readFileRequest), e.onReadFileRequest)
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+readFileResponse), e.onReadFileResponse)
	return &e
}

func (e *protocols) ReadFileAction(id peer.ID, roothash, datahash, path string, size int64) error {
	var err error
	var hash string
	var offset int64
	var fstat fs.FileInfo
	var f *os.File

	if size != FragmentSize {
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
	defer f.Sync()
	req := &pb.ReadfileRequest{
		Roothash:    roothash,
		Datahash:    datahash,
		Offset:      offset,
		MessageData: e.ReadFileProtocol.NewMessageData(uuid.New().String(), false),
	}
	return e.readFileAction(id, req, func(req *pb.ReadfileRequest, resp *readMsgResp) error {
		num, err := f.Write(resp.Data[:resp.Length])
		if err != nil {
			return errors.Wrapf(err, "[write file]")
		}
		req.Offset += int64(num)
		return nil
	})
}

func (e *protocols) ReadFileActionWithExtension(id peer.ID, roothash, datahash string, writer io.Writer, extData []byte) error {
	req := &pb.ReadfileRequest{
		Roothash:    roothash,
		Datahash:    datahash,
		ExtendData:  extData,
		MessageData: e.ReadFileProtocol.NewMessageData(uuid.New().String(), false),
	}
	return e.readFileAction(id, req, func(req *pb.ReadfileRequest, resp *readMsgResp) error {
		num, err := writer.Write(resp.Data[:resp.Length])
		if err != nil {
			return errors.Wrapf(err, "[write file]")
		}
		req.Offset += int64(num)
		return nil
	})
}

func (e *ReadFileProtocol) readFileAction(
	id peer.ID, req *pb.ReadfileRequest, callback func(req *pb.ReadfileRequest, resp *readMsgResp,
	) error) error {
	// store request so response handler has access to it
	var respChan = make(chan bool, 1)
	e.ReadFileProtocol.Lock()
	for {
		if _, ok := e.ReadFileProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		e.ReadFileProtocol.requests[req.MessageData.Id] = &readMsgResp{
			ch: respChan,
		}
		break
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
		err := e.ReadFileProtocol.SendProtoMessage(id, protocol.ID(e.ProtocolPrefix+readFileRequest), req)
		if err != nil {
			return errors.Wrapf(err, "[SendProtoMessage]")
		}

		//
		timeout.Reset(P2PReadReqRespTime)
		select {
		case ok := <-respChan:
			if !ok {
				return errors.New(ERR_RespFailure)
			}
		case <-timeout.C:
			return errors.New(ERR_RespTimeOut)
		}

		e.ReadFileProtocol.Lock()
		resp, ok := e.ReadFileProtocol.requests[req.MessageData.Id]
		if !ok {
			e.ReadFileProtocol.Unlock()
			return errors.New(ERR_RespFailure)
		}
		e.ReadFileProtocol.Unlock()

		if resp.ReadfileResponse == nil {
			return errors.New(ERR_RespFailure)
		}
		if len(resp.ReadfileResponse.Data) == 0 || resp.Length == 0 {
			if resp.Code == P2PResponseFinish {
				return nil
			}
			return errors.New(ERR_RespFailure)
		}
		err = callback(req, resp)
		if err != nil {
			return err
		}
		if resp.Code == P2PResponseFinish {
			return nil
		}
	}
}

func (e *ReadFileProtocol) SetReadFileServiceHandle(handle ReadFileServerHandle) {
	e.handle = handle
}

func (e *ReadFileProtocol) DefaultReadFileServerHandle(req *pb.ReadfileRequest) (*pb.ReadfileResponse, error) {
	var code = P2PResponseOK
	fpath := filepath.Join(e.ReadFileProtocol.GetDirs().TmpDir, req.Roothash, req.Datahash)

	if _, err := os.Stat(fpath); err != nil {
		fpath = filepath.Join(e.ReadFileProtocol.GetDirs().FileDir, req.Roothash, req.Datahash)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return &pb.ReadfileResponse{Code: P2PResponseFailed}, err
	}
	defer f.Close()

	fstat, err := f.Stat()
	if err != nil {
		return &pb.ReadfileResponse{Code: P2PResponseRemoteFailed}, err
	}

	if _, err = f.Seek(req.Offset, 0); err != nil {
		return &pb.ReadfileResponse{Code: P2PResponseRemoteFailed}, err
	}

	var readBuf = make([]byte, FileProtocolBufSize)
	num, err := f.Read(readBuf)
	if err != nil {
		return &pb.ReadfileResponse{Code: P2PResponseRemoteFailed}, err
	}

	if num+int(req.Offset) >= int(fstat.Size()) {
		code = P2PResponseFinish
	}

	// send response to the request using the message string he provided
	resp := &pb.ReadfileResponse{
		Code:   code,
		Offset: req.Offset,
		Length: uint32(num),
		Data:   readBuf[:num],
	}
	return resp, nil
}

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileRequest(s network.Stream) {
	defer s.Close()

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
	handle := e.handle
	if handle == nil {
		handle = e.DefaultReadFileServerHandle
	}
	resp, err := handle(data)
	if err != nil && resp == nil {
		s.Reset()
		return
	}
	resp.MessageData = e.ReadFileProtocol.NewMessageData(data.MessageData.Id, false)
	e.ReadFileProtocol.SendProtoMessage(
		s.Conn().RemotePeer(),
		protocol.ID(e.ProtocolPrefix+readFileResponse),
		resp,
	)
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
