/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/pkg/errors"
)

const FILE_PROTOCOL = "/kft/1"

type FileProtocol struct {
	*Node
}

func (n *Node) NewFileProtocol() *FileProtocol {
	e := FileProtocol{Node: n}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+FILE_PROTOCOL), e.onFileRequest)
	return &e
}

func (e *protocols) FileReq(peerId peer.ID, filehash string, filetype pb.FileType, fpath string) (uint32, error) {
	fstat, err := os.Stat(fpath)
	if err != nil {
		return 0, err
	}
	var req = &pb.Request{}
	var putReq = &pb.PutRequest{
		Type: filetype,
		Hash: filehash,
		Size: uint64(fstat.Size()),
	}

	s, err := e.FileProtocol.NewStream(context.Background(), peerId, protocol.ID(e.ProtocolPrefix+FILE_PROTOCOL))
	if err != nil {
		return 0, errors.Wrapf(err, "[NewStream]")
	}
	defer s.Close()

	putReq.Data, err = os.ReadFile(fpath)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[ReadFile]")
	}

	var reqMsg = &pb.Request_PutRequest{
		PutRequest: putReq,
	}
	reqMsg.PutRequest = putReq
	req.Request = reqMsg

	w := pbio.NewDelimitedWriter(s)
	err = w.WriteMsg(req)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[WriteMsg]")
	}

	respMsg := &pb.PutResponse{}

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, errors.Wrapf(err, "[ReadMsg]")
	}

	if respMsg.Code == 0 {
		e.FileProtocol.SetServiceFileTee(peerId.String())
	}

	return respMsg.Code, nil
}

// remote peer requests handler
func (e *FileProtocol) onFileRequest(s network.Stream) {
	defer s.Close()

	var resp = &pb.Response{}
	var putResp = &pb.PutResponse{
		Code: 1,
	}
	var respMsg = &pb.Response_PutResponse{
		PutResponse: putResp,
	}
	resp.Response = respMsg

	var reqMsg = &pb.Request{}

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)
	err := r.ReadMsg(reqMsg)
	if err != nil {
		w.WriteMsg(resp)
		return
	}

	switch reqMsg.GetRequest().(type) {
	case *pb.Request_PutRequest:
		putReq := reqMsg.GetPutRequest()
		switch putReq.Type {
		case pb.FileType_IdleData:
			fpath := filepath.Join(e.FileProtocol.GetDirs().IdleDataDir, putReq.Hash)
			err = saveFileStream(r, w, reqMsg, resp, fpath, putReq.Size, putReq.Data)
			if err == nil {
				putResp.Code = 0
			}
			e.FileProtocol.putIdleDataCh(fpath)
		default:
		}
	default:
	}
	w.WriteMsg(resp)
	return
}

func saveFileStream(r pbio.ReadCloser, w pbio.WriteCloser, reqMsg *pb.Request, resp *pb.Response, fpath string, size uint64, data []byte) error {
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "[OpenFile]")
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrapf(err, "[f.Write]")
	}

	putReq := &pb.PutRequest{}
	for bytesRead := uint64(0); bytesRead < size; {
		// Receive bytes
		err := r.ReadMsg(reqMsg)
		if err != nil {
			if err != io.EOF {
				return errors.Wrapf(err, "[ReadMsg]")
			}
			break
		}

		err = w.WriteMsg(resp)
		if err != nil {
			return errors.Wrapf(err, "[WriteMsg]")
		}

		putReq = reqMsg.GetPutRequest()

		// Write bytes to the file
		_, err = f.Write(putReq.Data)
		if err != nil {
			return errors.Wrapf(err, "[f.Write]")
		}
	}

	err = f.Sync()
	if err != nil {
		return errors.Wrapf(err, "[f.Sync]")
	}

	fstat, err := f.Stat()
	if err != nil {
		return errors.Wrapf(err, "[f.Stat]")
	}

	if uint64(fstat.Size()) != size {
		return errors.New("invalid size")
	}
	return nil
}
