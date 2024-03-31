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

	ZeroFileHash_16M = "080acf35a507ac9849cfcba47dc2ad83e01b75663a516279c8b9d243b719643e"

	DirMode = 0755

	SIZE_1KiB = 1024
	SIZE_1MiB = 1024 * SIZE_1KiB
	SIZE_1GiB = 1024 * SIZE_1MiB

	BufferSize   = 64 * SIZE_1KiB
	SegmentSize  = 64 * SIZE_1MiB
	FragmentSize = 16 * SIZE_1MiB
)
