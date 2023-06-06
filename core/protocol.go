/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Protocol interface {
	WriteFileAction(id peer.ID, roothash, path string) error
	ReadFileAction(id peer.ID, roothash, datahash, path string, size int64) error
	TagPushReq(peerid peer.ID) (uint32, error)
	IdleReq(peerId peer.ID, filesize, blocknum uint64, pubkey, sign []byte) (uint32, error)
	TagReq(peerId peer.ID, filename, customdata string, blocknum uint64) (uint32, error)
	FileReq(peerId peer.ID, filehash string, filetype pb.FileType, fpath string) (uint32, error)
	AggrProofReq(peerId peer.ID, ihash, shash []byte, qslice []*pb.Qslice, puk, sign []byte) (uint32, error)
	// add other protocols here...
}

type protocols struct {
	*WriteFileProtocol
	*ReadFileProtocol
	*CustomDataTagProtocol
	*IdleDataTagProtocol
	*FileProtocol
	*AggrProofProtocol
	*PushTagProtocol
	// add other protocols here...
}

func NewProtocol() *protocols {
	return &protocols{}
}
