/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

const p2pVersion = "go-p2p-node/0"
const privatekeyFile = ".private"

const AllIpAddress = "0.0.0.0"
const LocalAddress = "127.0.0.1"

// byte size
const (
	SIZE_1KiB = 1024
	SIZE_1MiB = 1024 * SIZE_1KiB
	SIZE_1GiB = 1024 * SIZE_1MiB
)

const (
	//
	FileDirectionry  = "file"
	TmpDirectionry   = "tmp"
	IdleDirectionry  = "space"
	TagDirectionry   = "tag"
	ProofDirectionry = "proof"
	//
	IdleNamesFile    = "inames"
	IdleUsFile       = "ius"
	IdleMuFile       = "imu"
	ServiceNamesFile = "snames"
	ServiceUsFile    = "sus"
	ServiceMuFile    = "smu"
)

const BufferSize = 64 * SIZE_1KiB

const DirMode = 0644

const SegmentSize = 16 * SIZE_1MiB
const FragmentSize = 8 * SIZE_1MiB

const DataShards = 2
const ParShards = 1
