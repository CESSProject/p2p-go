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

const P2PWriteReqRespTime = time.Duration(time.Second * 15)
const P2PReadReqRespTime = time.Duration(time.Second * 15)

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

var (
	FileNameLengthErr   = fmt.Errorf("filename length exceeds %d", MaxFileNameLength)
	FileNameEmptyErr    = fmt.Errorf("filename is empty")
	CustomDataLengthErr = fmt.Errorf("custom data length exceeds %d", MaxCustomDataLength)
)

type Names struct {
	Name []string `json:"name"`
}

type Us struct {
	U []string `json:"u"`
}
