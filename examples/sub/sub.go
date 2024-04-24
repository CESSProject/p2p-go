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
	"strings"
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
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

var room string

func main() {
	ctx := context.Background()
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("please enter os.Args[1] as port")
		os.Exit(1)
	}

	peer_node, err := p2pgo.New(
		ctx,
		p2pgo.ListenPort(port),
		p2pgo.Workspace("/home/test/sub"),
		p2pgo.BucketSize(100),
		p2pgo.BootPeers([]string{
			"_dnsaddr.boot-bucket-devnet.cess.cloud",
		}),
	)
	if err != nil {
		panic(err)
	}
	defer peer_node.Close()

	fmt.Println(peer_node.Addrs(), peer_node.ID())

	// create a new PubSub service using the GossipSub router
	gossipSub, err := pubsub.NewGossipSub(ctx, peer_node.GetHost())
	if err != nil {
		panic(err)
	}
	fmt.Println(peer_node.GetRendezvousVersion())
	fmt.Println(peer_node.GetDhtProtocolVersion())
	fmt.Println(peer_node.GetProtocolVersion())
	fmt.Println(peer_node.GetBootnode())

	//go Discover(ctx, peer_node.GetHost(), peer_node.GetDHTable(), peer_node.GetRendezvousVersion())

	// setup local mDNS discovery
	if err := setupDiscovery(peer_node.GetHost()); err != nil {
		panic(err)
	}

	if strings.Contains(peer_node.GetBootnode(), "12D3KooWGDk9JJ5F6UPNuutEKSbHrTXnF5eSn3zKaR27amgU6o9S") {
		room = fmt.Sprintf("%s-12D3KooWGDk9JJ5F6UPNuutEKSbHrTXnF5eSn3zKaR27amgU6o9S", core.NetworkRoom)
	} else if strings.Contains(peer_node.GetBootnode(), "12D3KooWRm2sQg65y2ZgCUksLsjWmKbBtZ4HRRsGLxbN76XTtC8T") {
		room = fmt.Sprintf("%s-12D3KooWRm2sQg65y2ZgCUksLsjWmKbBtZ4HRRsGLxbN76XTtC8T", core.NetworkRoom)
	} else if strings.Contains(peer_node.GetBootnode(), "12D3KooWEGeAp1MvvUrBYQtb31FE1LPg7aHsd1LtTXn6cerZTBBd") {
		room = fmt.Sprintf("%s-12D3KooWEGeAp1MvvUrBYQtb31FE1LPg7aHsd1LtTXn6cerZTBBd", core.NetworkRoom)
	} else {
		panic("Failed to connect to boot node")
	}

	// join the pubsub topic called librum
	topic, err := gossipSub.Join(room)
	if err != nil {
		panic(err)
	}
	fmt.Println("join a room: ", room)
	// subscribe to topic
	subscriber, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	subscribe(subscriber, ctx, peer_node.GetHost().ID())
}

// start subsriber to topic
func subscribe(subscriber *pubsub.Subscription, ctx context.Context, hostID peer.ID) {
	var findpeer peer.AddrInfo
	for {
		msg, err := subscriber.Next(ctx)
		if err != nil {
			panic(err)
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == hostID {
			continue
		}

		err = json.Unmarshal(msg.Data, &findpeer)
		if err != nil {
			continue
		}

		fmt.Printf("got message: %v, from: %s\n", findpeer, msg.ReceivedFrom.String())
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
	fmt.Printf("discovered new peer %s\n", pi.ID.String())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.String(), err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, room, &discoveryNotifee{h: h})
	return s.Start()
}
