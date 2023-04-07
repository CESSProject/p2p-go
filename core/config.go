/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
)

const p2pVersion = "go-p2p-node/0"
const privatekeyFile = ".private"

const AllIpAddress = "0.0.0.0"
const LocalAddress = "127.0.0.1"

var yamuxOpt = libp2p.Muxer("/yamux", yamux.DefaultTransport)
var mplexOpt = libp2p.Muxer("/mplex", mplex.DefaultTransport)
