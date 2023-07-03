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
	"github.com/CESSProject/p2p-go/core"
)

type Nnode struct {
	*core.Node
}

func main() {
	var ok bool
	var nnode = &Nnode{}
	ctx := context.Background()
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("please enter os.Args[1] as port")
		os.Exit(1)
	}

	h1, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private1"),
		p2pgo.ListenPort(port),
		p2pgo.Workspace("."),
		p2pgo.BootPeers([]string{
			"/ip4/221.122.79.2/tcp/10010/p2p/12D3KooWHY6BRu2MtG9SempACgYCcGHRSEai2ZkWY3E4VKDYrqh9",
			"/ip4/45.77.47.184/tcp/10010/p2p/12D3KooWBW5YSqJtABaaTmMZ1ByARcsTtCmmfB6na5HBEuUoKkLM",
			"/ip4/221.122.79.3/tcp/10010/p2p/12D3KooWAdyc4qPWFHsxMtXvSrm7CXNFhUmKPQdoXuKQXki69qBo",
		}),
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	nnode.Node, ok = h1.(*core.Node)
	if !ok {
		panic(err)
	}

	fmt.Println(nnode.Addrs(), nnode.ID())

	for {
		select {
		case peer := <-nnode.DiscoveredPeer():
			fmt.Println("found: ", peer.ID.Pretty())
			fmt.Println("found: ", peer.Addrs)
		}
	}

}
