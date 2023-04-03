package protocol

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
const ReadfileProtocol = "readfile"

const readFileRequest = "/file/readreq/v0"
const readFileResponse = "/file/readresp/v0"

func NewReadFileProtocol(node *core.Node) *ReadFileProtocol {
	e := ReadFileProtocol{node: node, requests: make(map[string]*readMsgResp)}
	node.SetStreamHandler(readFileRequest, e.onReadFileRequest)
	node.SetStreamHandler(readFileResponse, e.onReadFileResponse)
	return &e
}

func (e *ReadFileProtocol) ReadFileAction(id peer.ID, roothash, datahash, path string, size int64) error {
	log.Printf("Will Sending readfileAction to: %s....", id)

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
				return err
			}
			if hash == datahash {
				return nil
			}
			return fmt.Errorf("datahash does not match file")
		} else {
			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash, err = CalcSHA256(buf[:size])
			if err != nil {
				return err
			}
			if hash == datahash {
				f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = f.Write(buf[:size])
				if err != nil {
					return err
				}
				err = f.Sync()
				if err != nil {
					return err
				}
				return nil
			}
			os.Remove(path)
		}
	}
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := e.node.NewStream(context.Background(), id, readFileRequest)
	if err != nil {
		return err
	}
	defer s.Close()

	req.Roothash = roothash
	req.Datahash = datahash
	req.MessageData = e.node.NewMessageData(uuid.New().String(), false)

	// store request so response handler has access to it
	var respChan = make(chan bool, 1)
	e.requests[req.MessageData.Id] = &readMsgResp{
		ch: respChan,
	}
	defer delete(e.requests, req.MessageData.Id)
	defer close(respChan)
	timeout := time.NewTicker(P2PReadReqRespTime)
	defer timeout.Stop()

	for {
		req.Offset = offset
		// calc signature
		signature, err := e.node.SignProtoMessage(&req)
		if err != nil {
			log.Println("failed to sign message")
			return err
		}

		// add the signature to the message
		req.MessageData.Sign = signature

		err = e.node.SendProtoMessageToStream(s, &req)
		if err != nil {
			return err
		}

		//
		timeout.Reset(P2PReadReqRespTime)
		select {
		case ok = <-e.requests[req.MessageData.Id].ch:
			if !ok {
				// err, close
				return errors.New("failed")
			}
		case <-timeout.C:
			// timeout
			return errors.New("timeout")
		}
		resp = e.requests[req.MessageData.Id].ReadfileResponse
		num, err = f.Write(resp.Data[:resp.Length])
		if err != nil {
			return err
		}

		err = f.Sync()
		if err != nil {
			return err
		}

		if resp.Code == P2PResponseFinish {
			break
		}

		req.Offset += int64(num)
	}

	return nil
}

// remote peer requests handler
func (e *ReadFileProtocol) onReadFileRequest(s network.Stream) {
	var code = P2PResponseOK
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

	valid := e.node.AuthenticateMessage(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("Sending Readfile response to %s. Message id: %s...", s.Conn().RemotePeer(), data.MessageData.Id)

	f, err := os.OpenFile(filepath.Join(e.node.Workspace(), FileDirectionry, data.Datahash), os.O_RDONLY, 0)
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
	num, err := f.Read(buf)
	if err != nil {
		return
	}

	if num+int(data.Offset) >= int(fstat.Size()) {
		code = P2PResponseFinish
	}

	// send response to the request using the message string he provided
	resp := &pb.ReadfileResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false),
		Code:        code,
		Offset:      data.Offset,
		Length:      uint32(num),
		Data:        readBuf[:num],
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
	valid := e.node.AuthenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate request data and remove it if found
	_, ok := e.requests[data.MessageData.Id]
	if ok {
		if data.Code == P2PResponseOK || data.Code == P2PResponseFinish {
			e.requests[data.MessageData.Id].ReadfileResponse = data
		} else {
			e.requests[data.MessageData.Id].ch <- false
		}
	} else {
		log.Println("Failed to locate request data boject for response")
		return
	}

	log.Printf("Received Readfile response from %s. Message id:%s. Code: %d Offset:%d, Data:%s",
		s.Conn().RemotePeer(), data.MessageData.Id, data.Code, data.Offset, string(data.Data))
}
