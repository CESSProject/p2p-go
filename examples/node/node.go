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
	"strconv"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/config"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("please enter os.Args[1] as port")
		os.Exit(1)
	}

	peerNode, err := p2pgo.New(
		context.Background(),
		p2pgo.ListenPort(port),
		p2pgo.Workspace("."),
		p2pgo.BootPeers([]string{"_dnsaddr.boot-bucket-devnet.cess.cloud"}),
		p2pgo.BucketSize(100),
		p2pgo.ProtocolPrefix(config.TestnetProtocolPrefix),
	)
	if err != nil {
		panic(err)
	}
	defer peerNode.Close()

	fmt.Println(peerNode.Addrs(), peerNode.ID())
}
