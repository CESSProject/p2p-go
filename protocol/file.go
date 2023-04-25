package protocol

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const FILE_PROTOCOL = "/file/req/1"

type FileProtocol struct {
	node     *core.Node                 // local host
	requests map[string]*pb.FileRequest // used to access request data from response handlers
}

func NewFileProtocol(node *core.Node) *FileProtocol {
	e := FileProtocol{node: node, requests: make(map[string]*pb.FileRequest)}
	node.SetStreamHandler(FILE_PROTOCOL, e.onFileRequest)
	return &e
}

func (e *FileProtocol) FileReq(peerId peer.ID, filetype uint32, fpath string) error {
	log.Printf("Sending file req to: %s", peerId)

	fstat, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	reqMsg := &pb.FileRequest{
		MessageData: &pb.Messagedata{
			Id: uuid.NewString(),
		},
		FileType: filetype,
		FileSize: uint64(fstat.Size()),
	}

	s, err := e.node.NewStream(context.Background(), peerId, FILE_PROTOCOL)
	if err != nil {
		return err
	}
	defer s.Close()

	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	w := pbio.NewDelimitedWriter(s)

	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf = make([]byte, FileProtocolMsgBuf)
	var num int
	var respMsg = &pb.FileResponse{}
	for {
		num, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if num == 0 {
			break
		}

		reqMsg.Data = buf[:num]
		reqMsg.DataLength = uint32(num)
		reqMsg.MessageData.Timestamp = time.Now().UnixMilli()

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
		if respMsg.Code != P2PResponseOK {
			return fmt.Errorf("File req returns failed: %d", respMsg.Code)
		}
	}
	log.Printf("File req suc")
	return nil
}

// remote peer requests handler
func (e *FileProtocol) onFileRequest(s network.Stream) {
	r := pbio.NewDelimitedReader(s, FileProtocolMsgBuf)
	reqMsg := &pb.FileRequest{}
	err := r.ReadMsg(reqMsg)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}

	log.Printf("receive file req: %d", reqMsg.FileType)

	w := pbio.NewDelimitedWriter(s)
	respMsg := &pb.TagResponse{
		Code: P2PResponseOK,
		MessageData: &pb.Messagedata{
			Id: reqMsg.MessageData.Id,
		},
	}

	defer func() {
		respMsg.MessageData.Timestamp = time.Now().UnixMilli()
		err = w.WriteMsg(respMsg)
		if err != nil {
			s.Reset()
			log.Println(err)
		}
	}()

	var f *os.File
	switch reqMsg.FileType {
	case FileType_ServiceFile:
	case FileType_IdleFile:
		f, err = os.OpenFile(filepath.Join(e.node.Workspace(), IdleDirectionry), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			respMsg.Code = P2PResponseFailed
			return
		}
		defer f.Close()
	case FileType_TagFile:
		f, err = os.OpenFile(filepath.Join(e.node.Workspace(), TagDirectionry), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			respMsg.Code = P2PResponseFailed
			return
		}
		defer f.Close()
	case FileType_MusFile:
		f, err = os.OpenFile(filepath.Join(e.node.Workspace(), MusDirectionry), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			respMsg.Code = P2PResponseFailed
			return
		}
		defer f.Close()
	case FileType_UsFile:
		f, err = os.OpenFile(filepath.Join(e.node.Workspace(), UsDirectionry), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			respMsg.Code = P2PResponseFailed
			return
		}
		defer f.Close()
	case FileType_NamesFile:
		f, err = os.OpenFile(filepath.Join(e.node.Workspace(), NamesDirectionry), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			respMsg.Code = P2PResponseFailed
			return
		}
		defer f.Close()
	default:
		respMsg.Code = P2PResponseFailed
		log.Println("Unknown file type")
		return
	}

	f.Write(reqMsg.Data[:reqMsg.DataLength])
	err = f.Sync()
	if err != nil {
		log.Println(err)
		respMsg.Code = P2PResponseFailed
		return
	}

	log.Printf("%s: File response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
}
