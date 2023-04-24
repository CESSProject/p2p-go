/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package protocol

import (
	"fmt"
	"time"
)

const FileDirectionry = "file"
const TmpDirectionry = "tmp"

const P2PWriteReqRespTime = time.Duration(time.Second * 15)
const P2PReadReqRespTime = time.Duration(time.Second * 15)

const FileProtocolBufSize = 2 * 1024 * 1024

const P2PResponseOK uint32 = 200
const P2PResponseFinish uint32 = 210
const P2PResponseFailed uint32 = 400

const TagProtocolMsgBuf = 1024

const MaxFileNameLength = 255
const MaxCustomDataLength = 255

var (
	FileNameLengthErr = fmt.Errorf("filename length exceeds %d", MaxFileNameLength)
	FileNameEmptyErr  = fmt.Errorf("filename is empty")

	CustomDataLengthErr = fmt.Errorf("custom data length exceeds %d", MaxCustomDataLength)
)
