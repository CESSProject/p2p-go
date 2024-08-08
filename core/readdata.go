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

type readDataResp struct {
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
	return &e
}

func (e *protocols) ReadDataAction(id peer.ID, name, savepath string, size int64) error {
	if size <= 0 {
		return errors.New("invalid size")
	}

	offset, err := checkFileSize(savepath)
	if err != nil {
		return err
	}

	var f *os.File
	if offset > 0 {
		if offset > size {
			return fmt.Errorf("the file already exists and the size is greater than %d", size)
		}
		if offset == size {
			return nil
		}
		f, err = os.OpenFile(savepath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "open an existing file: ")
		}
	} else {
		f, err = os.OpenFile(savepath, os.O_WRONLY|os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "open file: ")
		}
	}
	defer f.Close()

	req := pb.ReadDataRequest{
		Roothash:    "",
		Datahash:    name,
		MessageData: e.ReadDataProtocol.NewMessageData(uuid.New().String(), false),
	}

	e.ReadDataProtocol.Lock()
	for {
		if _, ok := e.ReadDataProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		break
	}
	e.ReadDataProtocol.Unlock()

	defer func() {
		e.ReadDataProtocol.Lock()
		delete(e.ReadDataProtocol.requests, req.MessageData.Id)
		e.ReadDataProtocol.Unlock()
	}()

	stream, err := e.ReadDataProtocol.NewPeerStream(id, protocol.ID(e.ProtocolPrefix+readDataRequest))
	if err != nil {
		return err
	}
	defer func() {
		stream.Reset()
		stream.Close()
	}()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	num := 0
	tmpbuf := make([]byte, 0)
	recvbuf := make([]byte, 33*1024)
	recvdata := &pb.ReadDataResponse{}
	timeout := time.NewTicker(time.Second * 15)
	defer timeout.Stop()
	for {
		req.Offset = offset
		buf, err := proto.Marshal(&req)
		if err != nil {
			return errors.Wrapf(err, "[Marshal]")
		}

		_, err = rw.Write(buf)
		if err != nil {
			return errors.Wrapf(err, "[rw.Write]")
		}

		err = rw.Flush()
		if err != nil {
			return errors.Wrapf(err, "[rw.Flush]")
		}

		timeout.Reset(time.Second * 15)
		select {
		case <-timeout.C:
			return fmt.Errorf("%s", ERR_RecvTimeOut)
		default:
			num, err = rw.Read(recvbuf)
			if err != nil {
				return fmt.Errorf("[rw.Read] %v", err)
			}
			err = proto.Unmarshal(recvbuf[:num], recvdata)
			if err != nil {
				tmpbuf = append(tmpbuf, recvbuf[:num]...)
				err = proto.Unmarshal(tmpbuf, recvdata)
				if err != nil {
					break
				}
				tmpbuf = make([]byte, 0)
			}

			if recvdata.Length > 0 {
				num, err = f.Write(recvdata.Data[:recvdata.Length])
				if err != nil {
					return errors.Wrapf(err, "[write file]")
				}
				err = f.Sync()
				if err != nil {
					return errors.Wrapf(err, "[sync file]")
				}
			}

			if recvdata.Code == P2PResponseFinish {
				return nil
			}

			if recvdata.Code != P2PResponseOK {
				return fmt.Errorf("received error code: %d", recvdata.Code)
			}

			offset = req.Offset + int64(num)
		}
	}
}

func (e *ReadDataProtocol) onReadDataRequest(s network.Stream) {
	defer func() {
		s.Close()
	}()
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	var wg sync.WaitGroup
	wg.Add(1)
	go e.readData(rw, &wg)
	wg.Wait()
}

func (e *ReadDataProtocol) readData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	var err error
	num := 0
	data := &pb.ReadDataRequest{}
	recvbuf := make([]byte, 1024)
	tmpbuf := make([]byte, 0)
	readBuf := make([]byte, FileProtocolBufSize)
	tick := time.NewTicker(time.Second * 15)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			return
		default:
			num, err = rw.Read(recvbuf)
			if err != nil {
				return
			}
			tick.Reset(time.Second * 15)
			err = proto.Unmarshal(recvbuf[:num], data)
			if err != nil {
				tmpbuf = append(tmpbuf, recvbuf[:num]...)
				err = proto.Unmarshal(tmpbuf, data)
				if err != nil {
					break
				}
				tmpbuf = make([]byte, 0)
			}

			fpath := FindFile(e.ReadDataProtocol.GetDirs().FileDir, data.Datahash)
			if fpath == "" {
				fpath = filepath.Join(e.ReadDataProtocol.GetDirs().TmpDir, data.Datahash)
			}

			resp := &pb.ReadDataResponse{
				MessageData: e.ReadDataProtocol.NewMessageData(data.MessageData.Id, false),
				Code:        P2PResponseOK,
				Offset:      data.Offset,
				Length:      0,
			}

			_, err = os.Stat(fpath)
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}

			f, err := os.Open(fpath)
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}
			defer f.Close()

			fstat, err := f.Stat()
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}

			_, err = f.Seek(data.Offset, 0)
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}

			num, err = f.Read(readBuf)
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}

			if num+int(data.Offset) >= int(fstat.Size()) {
				resp.Code = P2PResponseFinish
			}

			resp.Length = uint32(num)
			resp.Data = readBuf[:num]

			buffer, err := proto.Marshal(resp)
			if err != nil {
				resp.Code = P2PResponseRemoteFailed
				resp.Length = 0
				resp.Data = make([]byte, 0)
				buffer, _ := proto.Marshal(resp)
				rw.Write(buffer)
				rw.Flush()
				return
			}
			rw.Write(buffer)
			rw.Flush()
			if resp.Code == P2PResponseFinish {
				return
			}
		}
	}
}

func checkFileSize(file string) (int64, error) {
	fstat, err := os.Stat(file)
	if err != nil {
		return 0, nil
	}

	if fstat.IsDir() {
		return 0, fmt.Errorf("%s is a directory", file)
	}

	return fstat.Size(), nil
}
