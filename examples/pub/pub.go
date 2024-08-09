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

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/core"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Nnode struct {
	*core.PeerNode
}

func main() {
	//var ok bool
	ctx := context.Background()
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("please enter os.Args[1] as port")
		os.Exit(1)
	}

	peer_node, err := p2pgo.New(
		ctx,
		p2pgo.ListenPort(port),
		p2pgo.Workspace("/home/test/boot"),
		p2pgo.ProtocolPrefix(config.DevnetProtocolPrefix),
	)
	if err != nil {
		panic(err)
	}
	defer peer_node.Close()

	log.Println(peer_node.Addrs(), peer_node.ID().String())
	log.Println(peer_node.GetProtocolVersion())
	log.Println(peer_node.GetRendezvousVersion())
	log.Println(peer_node.GetDhtProtocolVersion())

	// create a new PubSub service using the GossipSub router
	gossipSub, err := pubsub.NewGossipSub(ctx, peer_node.GetHost())
	if err != nil {
		panic(err)
	}

	// join the pubsub topic called librum
	topic, err := gossipSub.Join(core.NetworkRoom)
	if err != nil {
		panic(err)
	}

	// create publisher
	publish(ctx, topic, peer_node)
}

// start publisher to topic
func publish(ctx context.Context, topic *pubsub.Topic, pnode *core.PeerNode) {
	ch := pnode.GetRecord()
	for p := range ch {
		log.Printf("Published peer: %s", p)
		topic.Publish(ctx, []byte(p))
	}
}
