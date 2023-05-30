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
	"os"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	file := "./readfile"
	ctx := context.Background()
	sourcePort1 := flag.Int("p1", 15000, "Source port number")
	sourcePort2 := flag.Int("p2", 15001, "Source port number")

	// To construct a simple host with all the default settings, just use `New`
	h1, err := p2pgo.New(
		ctx,
		".private1",
		p2pgo.ListenPort(*sourcePort1), // regular tcp connections
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	fmt.Println("node1:", h1.Addrs(), h1.ID())

	// To construct a simple host with all the default settings, just use `New`
	h2, err := p2pgo.New(
		ctx,
		".private2",
		p2pgo.ListenPort(*sourcePort2), // regular tcp connections
	)
	if err != nil {
		panic(err)
	}
	defer h2.Close()

	fmt.Println("node2:", h2.Addrs(), h2.ID())

	remote := fmt.Sprintf("/ip4/0.0.0.0/tcp/15001/p2p/%v", h2.ID())

	maddr, err := ma.NewMultiaddr(remote)
	if err != nil {
		fmt.Println("NewMultiaddr err: ", err)
		os.Exit(1)
	}
	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Println("AddrInfoFromP2pAddr err: ", err)
		os.Exit(1)
	}
	h1.Peerstore().AddAddr(info.ID, maddr, 0)

	go h1.ReadFileAction(info.ID, "roothash", "7f0221a07b204d83b743ee58091ddabe9daa2c9b1b7d0900a6dff7c2b0bd418d", file, 8388608)
	select {}
}
