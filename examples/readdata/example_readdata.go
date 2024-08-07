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
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
)

var BootPeers = []string{
	"_dnsaddr.boot-miner-devnet.cess.cloud",
}

func main() {
	file := "readfile"
	ctx := context.Background()
	sourcePort1 := flag.Int("p1", 15000, "Source port number")
	sourcePort2 := flag.Int("p2", 15001, "Source port number")

	// peer1
	peer1, err := p2pgo.New(
		ctx,
		p2pgo.ListenPort(*sourcePort1),
		p2pgo.Workspace("./peer1"),
		p2pgo.BootPeers(BootPeers),
	)
	if err != nil {
		panic(err)
	}
	defer peer1.Close()

	fmt.Println("node1:", peer1.Addrs(), peer1.ID())

	// peer2
	peer2, err := p2pgo.New(
		ctx,
		p2pgo.ListenPort(*sourcePort2),
		p2pgo.Workspace("./peer2"),
		p2pgo.BootPeers(BootPeers),
	)
	if err != nil {
		panic(err)
	}
	defer peer2.Close()

	fmt.Println("node2:", peer2.Addrs(), peer2.ID())

	peer1.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), time.Second*5)

	// you need to put the test.txt file in ./peer2/file directory
	// the size cannot exceed the size of test.txt
	err = peer1.ReadDataAction(peer2.ID(), "test.txt", file, 40)
	if err != nil {
		fmt.Println("ReadDataAction err: ", err)
		return
	}
	fmt.Println("ok")
}
