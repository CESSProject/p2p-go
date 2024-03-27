/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()
	sourcePort1 := flag.Int("p1", 15000, "Source port number")
	sourcePort2 := flag.Int("p2", 15001, "Source port number")
	// To construct a simple host with all the default settings, just use `New`
	h1, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private1"),
		p2pgo.ListenPort(*sourcePort1), // regular tcp connections
		p2pgo.Workspace("."),
	)
	if err != nil {
		log.Println("[p2pgo.New]", err)
		return
	}
	defer h1.Close()

	log.Println("node1:", h1.Addrs(), h1.ID())

	// To construct a simple host with all the default settings, just use `New`
	h2, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private2"),
		p2pgo.ListenPort(*sourcePort2), // regular tcp connections
		p2pgo.Workspace("."),
	)
	if err != nil {
		log.Println("[p2pgo.New]", err)
		return
	}
	defer h2.Close()

	log.Println("node2:", h2.Addrs(), h2.ID())

	remote := fmt.Sprintf("/ip4/0.0.0.0/tcp/15001/p2p/%v", h2.ID())

	maddr, err := ma.NewMultiaddr(remote)
	if err != nil {
		log.Println("[ma.NewMultiaddr]", err)
		return
	}
	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println("[peer.AddrInfoFromP2pAddr]", err)
		return
	}
	h1.Peerstore().AddAddr(info.ID, maddr, time.Hour)

	datasize, err := h1.ReadDataStatAction(info.ID, "00e5b0a7c4c1d2f89b6347b2c0ea90720e5eb01aa1cf0e9052794383c6afd305", "00e5b0a7c4c1d2f89b6347b2c0ea90720e5eb01aa1cf0e9052794383c6afd305")
	log.Println("datasize: ", datasize, " err: ", err)
}
