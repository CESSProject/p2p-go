/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/CESSProject/p2p-go/pb"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// pattern: /protocol-name/request-or-response-message/version
const writeFileRequest = "/file/writereq/v0"
const writeFileResponse = "/file/writeresp/v0"
const readFileRequest = "/file/readreq/v0"
const readFileResponse = "/file/readresp/v0"

type WriteFileProtocol struct {
	node     *Node                // local host
	requests map[string]chan bool // determine whether it is your own response
}

type ReadFileProtocol struct {
	node     *Node                // local host
	requests map[string]chan bool // determine whether it is your own response
}

func NewWriteFileProtocol(node *Node) *WriteFileProtocol {
	e := WriteFileProtocol{node: node, requests: make(map[string]chan bool)}
	node.SetStreamHandler(writeFileRequest, e.onWriteFileRequest)
	node.SetStreamHandler(writeFileResponse, e.onWriteFileResponse)
	return &e
}

func NewReadFileProtocol(node *Node) *ReadFileProtocol {
	e := ReadFileProtocol{node: node, requests: make(map[string]chan bool)}
	node.SetStreamHandler(readFileRequest, e.onReadFileRequest)
	node.SetStreamHandler(readFileResponse, e.onReadFileResponse)
	return &e
}

func (e *WriteFileProtocol) WriteFileAction(id peer.ID, path string) error {
	log.Printf("Will Sending writefileAction to: %s....", id)

	var err error
	var ok bool
	var num int
	var offset int64
	var respChan = make(chan bool, 1)
	var timeout = time.NewTicker(P2PRespTimeout)

	// create message data
	req := &pb.WritefileRequest{
		MessageData: e.node.NewMessageData(uuid.New().String(), false),
		Roothash:    "This is root hash",
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	req.Datahash, err = CalcFileSHA256(f)
	if err != nil {
		return err
	}

	s, err := e.node.NewStream(context.Background(), id, writeFileRequest)
	if err != nil {
		return err
	}
	defer s.Close()

	// store request so response handler has access to it
	e.requests[req.MessageData.Id] = respChan
	defer delete(e.requests, req.MessageData.Id)
	defer close(respChan)

	buf := make([]byte, FileProtocolBufSize)
	for {
		num, err = f.Read(buf)
		if err != nil && err != io.EOF {
			// err
			return err
		}
		if num == 0 {
			break
		}

		req.Data = buf[:num]
		req.Length = uint32(num)
		req.Offset = offset
		hash := sha256.Sum256(req.Data)

		// calc signature
		signature, err := e.node.signProtoMessage(req)
		if err != nil {
			log.Println("failed to sign message")
			return err
		}

		// add the signature to the message
		req.MessageData.Sign = signature

		err = e.node.sendProtoMessageToStream(s, req)
		if err != nil {
			return err
		}

		//
		timeout.Reset(P2PRespTimeout)
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
		log.Printf("Writefile to: %s was sent. Msg Id: %s, Data hash: %s", id, req.MessageData.Id, hex.EncodeToString(hash[:]))
	}
	return nil
}

func (e *ReadFileProtocol) ReadFileAction(id peer.ID, hash, path string) error {
	log.Printf("Will Sending readfileAction to: %s....", id)
	return nil
}

// remote peer requests handler
func (e *WriteFileProtocol) onWriteFileRequest(s network.Stream) {
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

	valid := e.node.authenticateMessage(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("Sending Writefile response to %s. Message id: %s...", s.Conn().RemotePeer(), data.MessageData.Id)

	// send response to the request using the message string he provided
	resp := &pb.WritefileResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false),
		Code:        P2PResponseOK,
		Offset:      0,
	}

	// sign the data
	signature, err := e.node.signProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	err = e.node.sendProtoMessage(s.Conn().RemotePeer(), writeFileResponse, resp)
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
	valid := e.node.authenticateMessage(data, data.MessageData)

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

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileRequest(s network.Stream) {
	// get request data
	data := &pb.ReadfileRequest{}
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

	log.Printf("Received Readfile from %s. Roothash:%s Datahash:%s offset:%d",
		s.Conn().RemotePeer(), data.Roothash, data.Datahash, data.Offset)

	valid := e.node.authenticateMessage(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("Sending Readfile response to %s. Message id: %s...", s.Conn().RemotePeer(), data.MessageData.Id)

	// send response to the request using the message string he provided
	resp := &pb.ReadfileResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false),
		Code:        P2PResponseOK,
		Offset:      0,
		Length:      4,
		Data:        []byte{'t', 'e', 's', 't'},
	}

	// sign the data
	signature, err := e.node.signProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	err = e.node.sendProtoMessage(s.Conn().RemotePeer(), writeFileResponse, resp)
	if err != nil {
		log.Printf("Writefile response to %s sent failed.", s.Conn().RemotePeer().String())
	}
}

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileResponse(s network.Stream) {
	data := &pb.ReadfileResponse{}
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
	valid := e.node.authenticateMessage(data, data.MessageData)

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

	log.Printf("Received Readfile response from %s. Message id:%s. Code: %d Offset:%d, Data:%s",
		s.Conn().RemotePeer(), data.MessageData.Id, data.Code, data.Offset, string(data.Data))
}
