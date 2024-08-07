/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"time"
)

const P2PWriteReqRespTime = time.Duration(time.Second * 30)
const P2PReadReqRespTime = time.Duration(time.Second * 30)

const FileProtocolBufSize = 32 * 1024

const P2PResponseOK uint32 = 200
const P2PResponseFinish uint32 = 210
const P2PResponseFailed uint32 = 400
const P2PResponseRemoteFailed uint32 = 500
const P2PResponseEmpty uint32 = 404
const P2PResponseForbid uint32 = 403

const MaxFileNameLength = 255
const MaxCustomDataLength = 255

const (
	FileDataDirectionry = "file"
	TmpDataDirectionry  = "tmp"
)

const (
	ERR_RecvTimeOut     = "receiving data timeout"
	ERR_WriteTimeOut    = "writing data timeout"
	ERR_RespFailure     = "peer response failure"
	ERR_RespInvalidData = "peer response invalid data"
	ERR_RespEmptyData   = "received empty data"
	ERR_InvalidData     = "received invalid data"
)

const (
	p2pProtocolVer = "/protoco/1.0"
	dhtProtocolVer = "/dht/1.0"
	rendezvous     = "/rendezvous/1.0"
)

type DataDirs struct {
	FileDir string
	TmpDir  string
}
