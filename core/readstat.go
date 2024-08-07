/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"bufio"
	"fmt"
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

type readDataStatResp struct {
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
	return &e
}

func (e *protocols) ReadDataStatAction(id peer.ID, fid string, fragment string) (uint64, error) {
	var err error
	var req = &pb.ReadDataStatRequest{
		Roothash:    fid,
		Datahash:    fragment,
		MessageData: e.ReadDataStatProtocol.NewMessageData(uuid.New().String(), false),
	}

	e.ReadDataStatProtocol.Lock()
	for {
		if _, ok := e.ReadDataStatProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		break
	}
	e.ReadDataStatProtocol.Unlock()

	defer func() {
		e.ReadDataStatProtocol.Lock()
		delete(e.ReadDataStatProtocol.requests, req.MessageData.Id)
		e.ReadDataStatProtocol.Unlock()
	}()

	stream, err := e.ReadDataStatProtocol.NewPeerStream(id, protocol.ID(e.ProtocolPrefix+readDataStatRequest))
	if err != nil {
		return 0, err
	}
	defer stream.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	num := 0
	tmpbuf := make([]byte, 0)
	recvbuf := make([]byte, 1024)
	recvdata := &pb.ReadDataStatResponse{}

	timeout := time.NewTicker(time.Second * 10)
	defer timeout.Stop()
	for {
		buf, err := proto.Marshal(req)
		if err != nil {
			return 0, fmt.Errorf("[proto.Marshal] %v", err)
		}

		_, err = rw.Write(buf)
		if err != nil {
			return 0, fmt.Errorf("[rw.Write] %v", err)
		}

		err = rw.Flush()
		if err != nil {
			return 0, fmt.Errorf("[rw.Flush] %v", err)
		}

		timeout.Reset(time.Second * 10)
		select {
		case <-timeout.C:
			return 0, errors.New(ERR_RecvTimeOut)
		default:
			num, _ = rw.Read(recvbuf)
			err = proto.Unmarshal(recvbuf[:num], recvdata)
			if err != nil {
				tmpbuf = append(tmpbuf, recvbuf[:num]...)
				err = proto.Unmarshal(tmpbuf, recvdata)
				if err != nil {
					break
				}
				tmpbuf = make([]byte, 0)
			}

			if recvdata.Code == P2PResponseFinish || recvdata.Code == P2PResponseOK {
				return uint64(recvdata.DataSize), nil
			}

			return 0, errors.New(ERR_RespFailure)
		}
	}
}

func (e *ReadDataStatProtocol) onReadDataStatRequest(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go e.readData(rw)
}

// remote peer requests handler
func (e *ReadDataStatProtocol) readData(rw *bufio.ReadWriter) {
	var err error
	resp := &pb.ReadDataStatResponse{}
	data := &pb.ReadDataStatRequest{}
	recvbuf := make([]byte, 1024)
	tmpbuf := make([]byte, 0)
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			return
		default:
			num, _ := rw.Read(recvbuf)
			err = proto.Unmarshal(recvbuf[:num], data)
			if err != nil {
				tmpbuf = append(tmpbuf, recvbuf[:num]...)
				err = proto.Unmarshal(tmpbuf, data)
				if err != nil {
					break
				}
				tmpbuf = make([]byte, 0)
			}
			fpath := FindFile(e.ReadDataStatProtocol.GetDirs().FileDir, data.Datahash)
			if fpath == "" {
				fpath = FindFile(e.ReadDataStatProtocol.GetDirs().TmpDir, data.Datahash)
			}
			fstat, err := os.Stat(fpath)
			if err != nil {
				resp.Code = P2PResponseEmpty
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}
			resp.Code = P2PResponseOK
			resp.DataSize = fstat.Size()
			buffer, _ := proto.Marshal(resp)
			rw.Write(buffer)
			rw.Flush()
			return
		}
	}
}
