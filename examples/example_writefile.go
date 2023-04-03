package main

import (
	"flag"
	"fmt"
	"os"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/core"
	myprotocol "github.com/CESSProject/p2p-go/protocol"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	file := "./example_writefile.go"

	sourcePort1 := flag.Int("p1", 15000, "Source port number")
	sourcePort2 := flag.Int("p2", 15001, "Source port number")

	// To construct a simple host with all the default settings, just use `New`
	h1, err := p2pgo.New(
		".private1",
		p2pgo.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort1), // regular tcp connections
		),
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	fmt.Println("node1:", h1.Addrs(), h1.ID())

	// To construct a simple host with all the default settings, just use `New`
	h2, err := p2pgo.New(
		".private2",
		p2pgo.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort2), // regular tcp connections
		),
	)
	if err != nil {
		panic(err)
	}
	defer h2.Close()

	fmt.Println("node2:", h2.Addrs(), h2.ID())

	node1, ok := h1.(*core.Node)
	if !ok {
		panic(err)
	}

	node2, ok := h2.(*core.Node)
	if !ok {
		panic(err)
	}

	protocol := myprotocol.NewProtocol(node1)
	protocol.WriteFileProtocol = myprotocol.NewWriteFileProtocol(node1)
	protocol.ReadFileProtocol = myprotocol.NewReadFileProtocol(node1)

	myprotocol.NewWriteFileProtocol(node2)

	remote := fmt.Sprintf("/ip4/0.0.0.0/tcp/15001/p2p/%v", node2.ID())

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
	node1.AddAddrToPearstore(info.ID, maddr, 0)

	go protocol.WriteFileAction(info.ID, "This is roothash", file)
	select {}
}
