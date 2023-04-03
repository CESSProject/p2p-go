/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// pattern: /protocol-name/request-or-response-message/version
const writeFileRequest = "/file/writereq/v0"
const writeFileResponse = "/file/writeresp/v0"

type readMsgResp struct {
	ch chan bool
	*pb.ReadfileResponse
}

type WriteFileProtocol struct {
	node     *core.Node           // local host
	requests map[string]chan bool // determine whether it is your own response
}

func NewWriteFileProtocol(node *core.Node) *WriteFileProtocol {
	e := WriteFileProtocol{node: node, requests: make(map[string]chan bool)}
	node.SetStreamHandler(writeFileRequest, e.onWriteFileRequest)
	node.SetStreamHandler(writeFileResponse, e.onWriteFileResponse)
	return &e
}

func (e *WriteFileProtocol) WriteFileAction(id peer.ID, roothash, path string) error {
	log.Printf("Will Sending writefileAction to: %s", id)
	var err error
	var ok bool
	var num int
	var offset int64
	var f *os.File

	// create message data
	req := &pb.WritefileRequest{
		MessageData: e.node.NewMessageData(uuid.New().String(), false),
		Roothash:    roothash,
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
	e.requests[req.MessageData.Id] = respChan
	defer delete(e.requests, req.MessageData.Id)
	defer close(respChan)

	timeout := time.NewTicker(P2PWriteReqRespTime)
	buf := make([]byte, FileProtocolBufSize)
	for {
		f.Seek(offset, 0)
		num, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if num == 0 {
			fmt.Println("num: ", num)
			break
		}

		req.Data = buf[:num]
		req.Length = uint32(num)
		req.Offset = offset
		hash := sha256.Sum256(req.Data)
		req.MessageData.Timestamp = time.Now().Unix()
		// calc signature
		req.MessageData.Sign = nil
		signature, err := e.node.SignProtoMessage(req)
		if err != nil {
			log.Println("failed to sign message")
			return err
		}

		// add the signature to the message
		req.MessageData.Sign = signature

		err = e.node.SendProtoMessage(id, writeFileRequest, req)
		if err != nil {
			return err
		}

		log.Printf("Writefile to: %s was sent. Msg Id: %s, Data hash: %s", id, req.MessageData.Id, hex.EncodeToString(hash[:]))

		// wait response
		timeout.Reset(P2PWriteReqRespTime)
		select {
		case ok = <-e.requests[req.MessageData.Id]:
			if !ok {
				// err, close
				return errors.New("failed")
			}
		case <-timeout.C:
			// timeout
			return errors.New("timeout")
		}
		offset += int64(num)
	}
	return nil
}

// remote peer requests handler
func (e *WriteFileProtocol) onWriteFileRequest(s network.Stream) {
	log.Printf("Recv writefileAction from: %s", s.ID())
	// get request data
	data := &pb.WritefileRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	hash := sha256.Sum256(data.Data)

	log.Printf("Received Writefile from %s. Roothash:%s Datahash:%s length:%d offset:%d Message sign: %s",
		s.Conn().RemotePeer(), data.Roothash, data.Datahash, data.Length, data.Offset, hex.EncodeToString(hash[:]))

	valid := e.node.AuthenticateMessage(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("Sending Writefile response to %s. Message id: %s", s.Conn().RemotePeer(), data.MessageData.Id)

	f, err := os.OpenFile(filepath.Join(e.node.Workspace(), FileDirectionry, data.Datahash), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println("OpenFile err:", err)
		return
	}
	defer f.Close()

	_, err = f.Write(data.Data[:data.Length])
	if err != nil {
		log.Println("Write err:", err)
		return
	}

	// send response to the request using the message string he provided
	resp := &pb.WritefileResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false),
		Code:        P2PResponseOK,
		Offset:      0,
	}

	// sign the data
	signature, err := e.node.SignProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	err = e.node.SendProtoMessage(s.Conn().RemotePeer(), writeFileResponse, resp)
	if err != nil {
		log.Printf("Writefile response to %s sent failed.", s.Conn().RemotePeer().String())
	}
}

// remote peer response handler
func (e *WriteFileProtocol) onWriteFileResponse(s network.Stream) {
	data := &pb.WritefileResponse{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	// authenticate message content
	valid := e.node.AuthenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate request data and remove it if found
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		if data.Code == P2PResponseOK {
			e.requests[data.MessageData.Id] <- true
		} else {
			e.requests[data.MessageData.Id] <- false
		}
	} else {
		log.Println("Failed to locate request data boject for response")
		return
	}

	log.Printf("Received Writefile response from %s. Message id:%s. Code: %d Offset:%d.", s.Conn().RemotePeer(), data.MessageData.Id, data.Code, data.Offset)
}
