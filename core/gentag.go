/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"log"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
)

const CustomDataTag_Protocol = "/kldr/cdtg/1"

type CustomDataTagProtocol struct {
	*Node
}

func (n *Node) NewCustomDataTagProtocol() *CustomDataTagProtocol {
	e := CustomDataTagProtocol{Node: n}
	return &e
}

func (e *protocols) TagReq(peerId peer.ID, filename, customdata string, blocknum uint64) (uint32, error) {
	log.Printf("Sending tag req to: %s", peerId)

	if err := checkFileName(filename); err != nil {
		return 0, err
	}

	if err := checkCustomData(customdata); err != nil {
		return 0, err
	}

	s, err := e.CustomDataTagProtocol.NewStream(context.Background(), peerId, CustomDataTag_Protocol)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	w := pbio.NewDelimitedWriter(s)
	reqMsg := &pb.CustomDataTagRequest{
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
	respMsg := &pb.CustomDataTagResponse{}
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
