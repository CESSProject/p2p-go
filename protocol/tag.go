package protocol

import (
	"context"
	"log"
	"time"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const TAG_PROTOCOL = "/tag/req/1"

type TagProtocol struct {
	node     *core.Node                // local host
	requests map[string]*pb.TagRequest // used to access request data from response handlers
}

func NewTagProtocol(node *core.Node) *TagProtocol {
	e := TagProtocol{node: node, requests: make(map[string]*pb.TagRequest)}
	//node.SetStreamHandler(TAG_PROTOCOL, e.onTagRequest)
	return &e
}

func (e *TagProtocol) TagReq(peerId peer.ID, filename, customdata string, blocknum int64) (uint32, error) {
	log.Printf("Sending tag req to: %s", peerId)

	if err := checkFileName(filename); err != nil {
		return 0, err
	}

	if err := checkCustomData(customdata); err != nil {
		return 0, err
	}

	s, err := e.node.NewStream(context.Background(), peerId, TAG_PROTOCOL)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.TagRequest{
		MessageData: &pb.Messagedata{
			Timestamp: time.Now().UnixMilli(),
			Id:        uuid.NewString(),
		},
		FileName:   filename,
		CustomData: customdata,
		BlockNum:   blocknum,
	}

	err = w.WriteMsg(reqMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
	respMsg := &pb.TagResponse{}
	err = r.ReadMsg(respMsg)
	if err != nil {
		s.Reset()
		return 0, err
	}

	log.Printf("Tag req resp code: %d", respMsg.Code)
	return respMsg.Code, nil
}

func checkFileName(filename string) error {
	if len(filename) > MaxFileNameLength {
		return FileNameLengthErr
	}
	if len(filename) == 0 {
		return FileNameEmptyErr
	}
	return nil
}

func checkCustomData(customdata string) error {
	if len(customdata) > MaxCustomDataLength {
		return CustomDataLengthErr
	}
	return nil
}

// remote peer requests handler
// func (e *EchoProtocol) onTagRequest(s network.Stream) {
// 	r := pbio.NewDelimitedReader(s, TagProtocolMsgBuf)
// 	reqMsg := &pb.TagRequest{}
// 	err := r.ReadMsg(reqMsg)
// 	if err != nil {
// 		s.Reset()
// 		log.Println(err)
// 		return
// 	}

// 	log.Printf("receive tag req: %s", string(reqMsg.FileName))

// 	w := pbio.NewDelimitedWriter(s)
// 	respMsg := &pb.TagResponse{
// 		// Message: reqMsg.Message,
// 		// MessageData: &pb.MessageData{
// 		// 	Timestamp: time.Now().UnixMilli(),
// 		// 	Id:        reqMsg.MessageData.Id,
// 		// },
// 	}
// 	w.WriteMsg(respMsg)

// 	log.Printf("%s: Tag response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
// }
