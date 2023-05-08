package main

import (
	"fmt"
	"os"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	myprotocol "github.com/CESSProject/p2p-go/protocol"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	h1, err := p2pgo.New(
		".private1",
		p2pgo.ListenAddrStrings("0.0.0.0", 8080), // regular tcp connections
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	fmt.Println("node1:", h1.Addrs(), h1.ID())
	node1, ok := h1.(*core.Node)
	if !ok {
		panic(err)
	}

	protocol := myprotocol.NewProtocol(node1)
	protocol.WriteFileProtocol = myprotocol.NewWriteFileProtocol(node1)
	protocol.ReadFileProtocol = myprotocol.NewReadFileProtocol(node1)
	protocol.CustomDataTagProtocol = myprotocol.NewCustomDataTagProtocol(node1)
	protocol.IdleDataTagProtocol = myprotocol.NewIdleDataTagProtocol(node1)
	protocol.MusProtocol = myprotocol.NewMusProtocol(node1)
	protocol.FileProtocol = myprotocol.NewFileProtocol(node1)

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
	protocol.Node.AddAddrToPearstore(info.ID, maddr, 0)
	protocol.IdleDataTagProtocol.IdleReq(info.ID, 8*1024*1024, 2, []byte("123456"))
	protocol.CustomDataTagProtocol.TagReq(info.ID, "123456", "", 2)
	err = protocol.FileProtocol.FileReq(info.ID, "123456", int32(pb.FileType_CustomData), "./main.go")
	if err != nil {
		println(err)
	}
	select {}
}
