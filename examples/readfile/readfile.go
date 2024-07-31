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

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/core"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var BootPeers = []string{
	"_dnsaddr.boot-miner-devnet.cess.cloud",
}

func main() {
	ctx := context.Background()
	sourcePort1 := flag.Int("p1", 15000, "Source port number")
	sourcePort2 := flag.Int("p2", 15001, "Source port number")

	// peer1
	peer1, err := p2pgo.New(
		ctx,
		p2pgo.Workspace("./peer1"),
		p2pgo.ListenPort(*sourcePort1),
		p2pgo.BootPeers(BootPeers),
	)
	if err != nil {
		panic(err)
	}
	defer peer1.Close()

	fmt.Println("peer1:", peer1.Addrs(), peer1.ID())

	// peer2
	peer2, err := p2pgo.New(
		ctx,
		p2pgo.Workspace("./peer2"),
		p2pgo.ListenPort(*sourcePort2),
		p2pgo.BootPeers(BootPeers),
	)
	if err != nil {
		panic(err)
	}
	defer peer2.Close()

	fmt.Println("peer2:", peer2.Addrs(), peer2.ID())

	remoteAddrs := peer2.Addrs()

	for _, v := range remoteAddrs {
		remoteAddr, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", v, peer2.ID().String()))
		if err != nil {
			fmt.Println("NewMultiaddr err: ", err)
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(remoteAddr)
		if err != nil {
			fmt.Println("AddrInfoFromP2pAddr err: ", err)
			continue
		}

		err = peer1.Connect(context.TODO(), *info)
		if err != nil {
			fmt.Println("Connect err: ", err)
			continue
		}

		err = peer1.ReadFileAction(info.ID, "fid", "fragment_hash", "test_file", core.FragmentSize)
		if err != nil {
			fmt.Println("ReadDataAction err: ", err)
			continue
		}
		fmt.Println("ok")
		return
	}
	fmt.Println("failed")
}
