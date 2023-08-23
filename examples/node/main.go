/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/out"
	"golang.org/x/time/rate"
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
			"_dnsaddr.boot-kldr-testnet.cess.cloud",
		}),
		p2pgo.ProtocolPrefix(config.TestnetProtocolPrefix),
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

	nnode.RouteTableFindPeers(0)

	tick := time.NewTicker(time.Second * 30)
	defer tick.Stop()

	var r = rate.Every(time.Second * 3)
	var limit = rate.NewLimiter(r, 1)

	for {
		select {
		case peer, ok := <-nnode.GetDiscoveredPeers():
			if !ok {
				out.Err("-------------------------Closed DiscoveredPeers chan----------------------")
				break
			}
			if limit.Allow() {
				out.Tip("-------------------------Reset----------------------")
				tick.Reset(time.Second * 30)
			}
			if len(peer.Responses) == 0 {
				out.Err("-------------------------response is nil----------------------")
				break
			}
			for _, v := range peer.Responses {
				log.Println("found: ", v.ID.Pretty(), v.Addrs)
			}
		case <-tick.C:
			out.Tip("-------------------------RouteTableFindPeers----------------------")
			_, err = nnode.RouteTableFindPeers(0)
			if err != nil {
				out.Tip("-------------------------RouteTableFindPeers err----------------------")
				out.Err(err.Error())
			}
		}
	}
}
