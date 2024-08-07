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
const writeDataRequest = "/data/writereq/v0"

//const readDataResponse = "/data/readresp/v0"

type writeDataResp struct {
	*pb.WritefileResponse
}

type WriteDataProtocol struct {
	*PeerNode // local host
	*sync.Mutex
	requests map[string]*writeDataResp // determine whether it is your own response
}

func (n *PeerNode) NewWriteDataProtocol() *WriteDataProtocol {
	e := WriteDataProtocol{PeerNode: n, Mutex: new(sync.Mutex), requests: make(map[string]*writeDataResp)}
	n.SetStreamHandler(protocol.ID(n.protocolPrefix+writeDataRequest), e.onWriteDataRequest)
	return &e
}

func (e *protocols) WriteDataAction(id peer.ID, file, fid, fragment string) error {
	fstat, err := os.Stat(file)
	if err != nil {
		return err
	}
	totalsize := fstat.Size()
	if totalsize <= 0 {
		return fmt.Errorf("%s: empty file", file)
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("[os.Open] %v", err)
	}
	defer f.Close()

	req := pb.WriteDataRequest{
		Roothash:    fid,
		Datahash:    fragment,
		MessageData: e.WriteDataProtocol.NewMessageData(uuid.New().String(), false),
		TotalSize:   totalsize,
	}

	e.WriteDataProtocol.Lock()
	for {
		if _, ok := e.WriteDataProtocol.requests[req.MessageData.Id]; ok {
			req.MessageData.Id = uuid.New().String()
			continue
		}
		break
	}
	e.WriteDataProtocol.Unlock()

	defer func() {
		e.WriteDataProtocol.Lock()
		delete(e.WriteDataProtocol.requests, req.MessageData.Id)
		e.WriteDataProtocol.Unlock()
	}()

	stream, err := e.ReadDataProtocol.NewPeerStream(id, protocol.ID(e.ProtocolPrefix+writeDataRequest))
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
	recvbuf := make([]byte, 1024)
	databuf := make([]byte, 3)
	recvdata := &pb.WriteDataResponse{}

	var offset int64 = 0
	req.TotalSize = totalsize
	timeout := time.NewTicker(time.Second * 5)
	defer timeout.Stop()
	for {
		req.Offset = offset
		_, err = f.Seek(req.Offset, 0)
		if err != nil {
			return fmt.Errorf("[f.Seek] %v", err)
		}

		num, err = f.Read(databuf)
		if err != nil {
			return fmt.Errorf("[f.Read] %v", err)
		}

		req.Data = databuf[:num]
		req.DataLength = int32(num)

		buf, err := proto.Marshal(&req)
		if err != nil {
			return fmt.Errorf("[proto.Marshal] %v", err)
		}

		_, err = rw.Write(buf)
		if err != nil {
			return fmt.Errorf("[rw.Write] %v", err)
		}

		err = rw.Flush()
		if err != nil {
			return fmt.Errorf("[rw.Flush] %v", err)
		}
		offset += int64(num)
		timeout.Reset(time.Second * 5)
		select {
		case <-timeout.C:
			return errors.New(ERR_WriteTimeOut)
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
			if recvdata.Code == P2PResponseFinish {
				return nil
			}

			if recvdata.Code != P2PResponseOK {
				return fmt.Errorf("received error code: %d", recvdata.Code)
			}
		}
	}
}

func (e *WriteDataProtocol) onWriteDataRequest(s network.Stream) {
	defer func() {
		s.Close()
	}()
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	var wg sync.WaitGroup
	wg.Add(1)
	go e.readData(rw, &wg)
	wg.Wait()
}

func (e *WriteDataProtocol) readData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
	var err error
	data := &pb.WriteDataRequest{}
	respdata := &pb.WriteDataResponse{}
	recvbuf := make([]byte, 33*1024)
	tmpbuf := make([]byte, 0)
	var f *os.File
	fpath := ""
	defer func() {
		wg.Done()
	}()

	time.Sleep(time.Second)
	if !e.GetRecvFlag() {
		respdata.Code = P2PResponseForbid
		buffer, _ := proto.Marshal(respdata)
		rw.Write(buffer)
		rw.Flush()
		return
	}

	num := 0
	timeout := time.NewTicker(time.Second * 10)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			os.Remove(fpath)
			return
		default:
			num, err = rw.Read(recvbuf)
			if err != nil {
				respdata.Code = P2PResponseFailed
				buffer, _ := proto.Marshal(respdata)
				rw.Write(buffer)
				rw.Flush()
				os.Remove(fpath)
				return
			}
			err = proto.Unmarshal(recvbuf[:num], data)
			if err != nil {
				tmpbuf = append(tmpbuf, recvbuf[:num]...)
				err = proto.Unmarshal(tmpbuf, data)
				if err != nil {
					break
				}
				tmpbuf = make([]byte, 0)
			}
			dir := filepath.Join(e.GetDirs().TmpDir, data.Roothash)
			err = os.MkdirAll(dir, DirMode)
			if err != nil {
				respdata.Code = P2PResponseRemoteFailed
				buffer, _ := proto.Marshal(respdata)
				rw.Write(buffer)
				rw.Flush()
				os.Remove(fpath)
				return
			}

			fpath = filepath.Join(dir, data.Datahash)
			if data.Datahash == ZeroFileHash_8M {
				err = writeZeroToFile(fpath, FragmentSize)
				if err != nil {
					respdata.Code = P2PResponseRemoteFailed
					buffer, _ := proto.Marshal(respdata)
					rw.Write(buffer)
					rw.Flush()
					os.Remove(fpath)
					return
				}
				respdata.Code = P2PResponseFinish
				buffer, _ := proto.Marshal(respdata)
				rw.Write(buffer)
				rw.Flush()
				os.Remove(fpath)
				return
			}

			if f == nil {
				f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
				if err != nil {
					respdata.Code = P2PResponseRemoteFailed
					buffer, _ := proto.Marshal(respdata)
					rw.Write(buffer)
					rw.Flush()
					return
				}
				defer f.Close()
			}

			_, err = f.Write(data.Data[:data.DataLength])
			if err != nil {
				respdata.Code = P2PResponseRemoteFailed
				buffer, _ := proto.Marshal(respdata)
				rw.Write(buffer)
				rw.Flush()
				os.Remove(fpath)
				return
			}

			fstat, err := f.Stat()
			if err != nil {
				respdata.Code = P2PResponseRemoteFailed
				buffer, _ := proto.Marshal(respdata)
				rw.Write(buffer)
				rw.Flush()
				os.Remove(fpath)
				return
			}

			if fstat.Size() >= data.TotalSize {
				respdata.Code = P2PResponseFinish
			} else {
				respdata.Code = P2PResponseOK
			}

			buffer, _ := proto.Marshal(respdata)
			rw.Write(buffer)
			rw.Flush()
			if respdata.Code == P2PResponseFinish {
				return
			}
			timeout.Reset(time.Second * 10)
		}
	}
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
