package protocol

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const FILE_PROTOCOL = "/kldr/sft/1"

type FileProtocol struct {
	node *core.Node // local host
	//requests map[string]*pb.FileRequest // used to access request data from response handlers
}

func NewFileProtocol(node *core.Node) *FileProtocol {
	e := FileProtocol{node: node} //requests: make(map[string]*pb.FileRequest)}
	node.SetStreamHandler(FILE_PROTOCOL, e.onFileRequest)
	return &e
}

func (e *FileProtocol) FileReq(peerId peer.ID, fileId string, filetype int32, fpath string) error {
	log.Printf("Sending file req to: %s", peerId)

	fstat, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	reqMsg := &pb.PutRequest{
		Hash: fileId,
		Size: uint64(fstat.Size()),
		Type: pb.FileType(filetype),
	}
	respMsg := &pb.PutResponse{}

	s, err := e.node.NewStream(context.Background(), peerId, FILE_PROTOCOL)
	if err != nil {
		return err
	}
	defer s.Close()

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return err
	}

	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return err
	}

	if respMsg.Code != 0 {
		s.Reset()
		return fmt.Errorf("return failed and code:%d", respMsg.Code)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf = make([]byte, FileProtocolMsgBuf)
	var num int

	for {
		num, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if num == 0 {
			break
		}

		_, err = rw.Write(buf[:num])
		if err != nil {
			return err
		}
	}

	log.Printf("File req suc")
	return nil
}

// remote peer requests handler
func (e *FileProtocol) onFileRequest(s network.Stream) {
	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	reqMsg := &pb.Request{}
	err := r.ReadMsg(reqMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}

	w := pbio.NewDelimitedWriter(s)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	switch reqMsg.GetRequest().(type) {
	case *pb.Request_PutRequest:
		log.Printf("receive put file req")
		respMsg := &pb.PutResponse{
			Code: 0,
		}

		err = w.WriteMsg(respMsg)
		if err != nil {
			s.Reset()
			log.Println(err)
			return
		}

		putReq := reqMsg.GetPutRequest()
		switch putReq.Type {
		case pb.FileType_IdleData:
			err = saveFileStream(rw, filepath.Join(e.node.Workspace(), IdleDirectionry, putReq.Hash), putReq.Size)
			if err != nil {
				log.Println(err)
			}
		case pb.FileType_Tag:
			err = saveFileStream(rw, filepath.Join(e.node.Workspace(), TagDirectionry, putReq.Hash), putReq.Size)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Printf("recv put file req and invalid file type")
		}
	case *pb.Request_GetRequest:
		log.Printf("receive get file req")
		getReq := reqMsg.GetGetRequest()
		switch getReq.Type {
		case pb.FileType_Mu:
			muPath := filepath.Join(e.node.Workspace(), MusDirectionry, getReq.Hash)
			f, err := os.Open(muPath)
			if err != nil {
				log.Println(err)
				s.Reset()
				return
			}
			defer f.Close()

			var buf = make([]byte, FileProtocolMsgBuf)
			var num int

			for {
				num, err = f.Read(buf)
				if err != nil && err != io.EOF {
					log.Println(err)
					s.Reset()
					return
				}

				if num == 0 {
					break
				}

				_, err = rw.Write(buf[:num])
				if err != nil {
					log.Println(err)
					s.Reset()
					return
				}
			}

		default:
			log.Printf("invalid file type")
		}
	default:
		log.Printf("receive invalid file req")
	}

	log.Printf("%s: File response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	s.Reset()
	return
}

func saveFileStream(rw *bufio.ReadWriter, fpath string, size uint64) error {
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	var num int
	var buf = make([]byte, FileProtocolMsgBuf)
	for bytesRead := uint64(0); bytesRead < size; {
		// Receive bytes
		num, err = rw.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if num == 0 {
			break
		}

		bytesRead += uint64(num)

		// Write bytes to the file
		f.Write(buf)
	}

	fstat, err := f.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	if uint64(fstat.Size()) != size {
		log.Printf("recv file size = %d not equal origin size = %d", fstat.Size(), size)
		return err
	}
	return nil
}
