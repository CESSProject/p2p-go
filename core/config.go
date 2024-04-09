/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

const (
	p2pVersion  = "go-p2p-node/0"
	NetworkRoom = "cess-network-room"

	privatekeyFile = ".private"

	AllIpAddress = "0.0.0.0"
	LocalAddress = "127.0.0.1"

	DefaultProtocolPrefix         = "/cess"
	DefaultProtocolOverridePrefix = "/cess/dht/"

	ZeroFileHash_8M = "2daeb1f36095b44b318410b3f4e8b5d989dcc7bb023d1426c492dab0a3053e74"

	DirMode = 0755

	SIZE_1KiB = 1024
	SIZE_1MiB = 1024 * SIZE_1KiB
	SIZE_1GiB = 1024 * SIZE_1MiB

	BufferSize   = 64 * SIZE_1KiB
	SegmentSize  = 32 * SIZE_1MiB
	FragmentSize = 8 * SIZE_1MiB
)
