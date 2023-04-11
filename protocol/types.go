/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package protocol

import "time"

const FileDirectionry = "file"
const TmpDirectionry = "tmp"

const P2PWriteReqRespTime = time.Duration(time.Second * 15)
const P2PReadReqRespTime = time.Duration(time.Second * 15)

const FileProtocolBufSize = 2 * 1024 * 1024

const P2PResponseOK uint32 = 200
const P2PResponseFinish uint32 = 210
const P2PResponseFailed uint32 = 400
