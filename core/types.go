/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"fmt"
	"time"
)

const P2PWriteReqRespTime = time.Duration(time.Second * 30)
const P2PReadReqRespTime = time.Duration(time.Second * 30)

const FileProtocolBufSize = 2 * 1024 * 1024

const P2PResponseOK uint32 = 200
const P2PResponseFinish uint32 = 210
const P2PResponseFailed uint32 = 400
const P2PResponseRemoteFailed uint32 = 500
const P2PResponseEmpty uint32 = 404

const MaxFileNameLength = 255
const MaxCustomDataLength = 255

const (
	FileDataDirectionry = "file"
	TmpDataDirectionry  = "tmp"
)

const (
	ERR_RespTimeOut     = "peer response timeout"
	ERR_RespFailure     = "peer response failure"
	ERR_RespInvalidData = "peer response invalid data"
)

const (
	p2pProtocolVer = "/protoco/1.0"
	dhtProtocolVer = "/dht/1.0"
	rendezvous     = "/rendezvous/1.0"
)

var (
	FileNameLengthErr   = fmt.Errorf("filename length exceeds %d", MaxFileNameLength)
	FileNameEmptyErr    = fmt.Errorf("filename is empty")
	CustomDataLengthErr = fmt.Errorf("custom data length exceeds %d", MaxCustomDataLength)
)

type DataDirs struct {
	FileDir string
	TmpDir  string
}
