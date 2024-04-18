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
	// To construct a simple host with all the default settings, just use `New`
	h1, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private1"),
		p2pgo.ListenPort(*sourcePort1), // regular tcp connections
		p2pgo.Workspace("."),
		p2pgo.BootPeers([]string{"_dnsaddr.boot-bucket-devnet.cess.cloud"}),
	)
	if err != nil {
		log.Println("[p2pgo.New]", err)
		return
	}
	defer h1.Close()

	log.Println("node1:", h1.Addrs(), h1.ID())

	remote := "/ip4/127.0.0.1/tcp/8001/p2p/12D3KooWGDk9JJ5F6UPNuutEKSbHrTXnF5eSn3zKaR27amgU6o9S"

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

	err = h1.OnlineAction(info.ID)
	fmt.Println(err)
}
