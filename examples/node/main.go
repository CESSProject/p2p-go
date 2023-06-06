/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"os"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()
	h1, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private1"),
		p2pgo.ListenPort(8080), // regular tcp connections
		p2pgo.Workspace("."),
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	fmt.Println("node1:", h1.Addrs(), h1.ID())

	maddr, err := ma.NewMultiaddr(os.Args[1])
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

	h1.IdleReq(info.ID, 8*1024*1024, 2, nil, []byte("123456"))
	h1.TagReq(info.ID, "123456", "", 2)
	_, err = h1.FileReq(info.ID, "123456", pb.FileType_CustomData, "./main.go")
	if err != nil {
		fmt.Println(err)
	}
	select {}
}
