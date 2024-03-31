/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/core"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

type Nnode struct {
	*core.PeerNode
}

var findpeers chan peer.AddrInfo

func init() {
	findpeers = make(chan peer.AddrInfo, 1000)
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
		p2pgo.BucketSize(1000),
		p2pgo.ProtocolPrefix(config.TestnetProtocolPrefix),
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

	go Discover(ctx, peer_node.GetHost(), peer_node.GetDHTable(), peer_node.GetRendezvousVersion())

	// setup local mDNS discovery
	if err := setupDiscovery(peer_node.GetHost()); err != nil {
		panic(err)
	}

	// join the pubsub topic called librum
	topic, err := gossipSub.Join(core.NetworkRoom)
	if err != nil {
		panic(err)
	}

	// create publisher
	publish(ctx, topic)
}

// start publisher to topic
func publish(ctx context.Context, topic *pubsub.Topic) {
	for p := range findpeers {
		data, err := json.Marshal(p)
		if err != nil {
			continue
		}
		log.Printf("Published peer: %v", p)
		topic.Publish(ctx, data)
	}
}

func Discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	var routingDiscovery = drouting.NewRoutingDiscovery(dht)

	dutil.Advertise(ctx, routingDiscovery, rendezvous)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers, err := routingDiscovery.FindPeers(ctx, rendezvous)
			if err != nil {
				log.Fatal(err)
			}

			for p := range peers {
				if p.ID == h.ID() {
					continue
				}
				findpeers <- p
				fmt.Printf("Find peer: %s\n", p.ID.String())
				if h.Network().Connectedness(p.ID) != network.Connected {
					_, err = h.Network().DialPeer(ctx, p.ID)
					if err != nil {
						continue
					}
					fmt.Printf("Connected to peer %s\n", p.ID.String())
				}
			}
		}
	}
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID != n.h.ID() {
		fmt.Printf("discovered new peer %s\n", pi.ID.String())
		err := n.h.Connect(context.Background(), pi)
		if err != nil {
			fmt.Printf("error connecting to peer %s: %s\n", pi.ID.String(), err)
		}
		findpeers <- pi
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, core.NetworkRoom, &discoveryNotifee{h: h})
	return s.Start()
}
