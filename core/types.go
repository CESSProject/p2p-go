/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

const P2PWriteReqRespTime = time.Duration(time.Second * 30)
const P2PReadReqRespTime = time.Duration(time.Second * 30)

const FileProtocolBufSize = 2 * 1024 * 1024

const P2PResponseOK uint32 = 200
const P2PResponseFinish uint32 = 210
const P2PResponseFailed uint32 = 400

const TagProtocolMsgBuf = 1024 * 1024
const FileProtocolMsgBuf = 16 * 1024 * 1024
const IdleProtocolMsgBuf = 1024
const MusProtocolMsgBuf = 32

const MaxFileNameLength = 255
const MaxCustomDataLength = 255

const (
	FileType_ServiceFile uint32 = iota
	FileType_IdleFile
	FileType_TagFile
	FileType_MusFile
	FileType_UsFile
	FileType_NamesFile
)

const (
	//
	FileDataDirectionry   = "file"
	TmpDataDirectionry    = "tmp"
	IdleDataDirectionry   = "space"
	IdleTagDirectionry    = "itag"
	ServiceTagDirectionry = "stag"
	ProofDirectionry      = "proof"
	//
	IdleProofFile    = "iproof"
	ServiceProofFile = "sproof"
)

const (
	ERR_RespTimeOut     = "peer response timeout"
	ERR_RespFailure     = "peer response failure"
	ERR_RespInvalidData = "peer response invalid data"
)

const (
	p2pProtocolVer = "/1.0"
	dhtProtocolVer = "/kad/1.0"
	rendezvous     = "/rendezvous/1.0.0"
	grpcProtocolID = "/grpc/1.0"
)

var (
	FileNameLengthErr   = fmt.Errorf("filename length exceeds %d", MaxFileNameLength)
	FileNameEmptyErr    = fmt.Errorf("filename is empty")
	CustomDataLengthErr = fmt.Errorf("custom data length exceeds %d", MaxCustomDataLength)
)

type DataDirs struct {
	FileDir       string
	TmpDir        string
	IdleDataDir   string
	IdleTagDir    string
	ServiceTagDir string
	ProofDir      string
	IproofFile    string
	SproofFile    string
}

type ProofFileType struct {
	Names []string `json:"names"`
	Us    []string `json:"us"`
	Mus   []string `json:"mus"`
	Sigma string   `json:"sigma"`
}

type DiscoveredPeer struct {
	PeerID peer.ID
	Addr   ma.Multiaddr
}
