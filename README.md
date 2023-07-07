<div align="center">

# Go p2p library for cess distributed storage system

[![GitHub license](https://img.shields.io/badge/license-Apache2-blue)](#LICENSE) <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.19-blue.svg" /></a> [![Go Reference](https://pkg.go.dev/badge/github.com/CESSProject/p2p-go.svg)](https://pkg.go.dev/github.com/CESSProject/p2p-go) [![build&test](https://github.com/CESSProject/p2p-go/actions/workflows/build&test.yml/badge.svg)](https://github.com/CESSProject/p2p-go/actions/workflows/build&test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/CESSProject/p2p-go)](https://goreportcard.com/report/github.com/CESSProject/p2p-go)

</div>

The go p2p library implementation of the CESS distributed storage network, which is customized based on [go-libp2p](https://github.com/libp2p/go-libp2p), has many advantages, and at the same time abandons the bandwidth and space redundancy caused by multiple file backups, making nodes more focused on storage allocation Give yourself data to make it more in line with the needs of the CESS network.

## 📝 Reporting a Vulnerability
If you find out any system bugs or you have a better suggestions, please send an email to frode@cess.one or join [CESS discord](https://discord.gg/mYHTMfBwNS) to communicate with us.

## 🏗 Usage
To get the package use the standard:
```
go get -u "github.com/CESSProject/p2p-go"
```

## 📖 Documentation 
Please refer to https://pkg.go.dev/github.com/CESSProject/p2p-go

## 💡 Examples

The following code demonstrates how to create a p2p node and perform node discovery:
```
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	p2pgo "github.com/CESSProject/p2p-go"
	"github.com/CESSProject/p2p-go/core"
	"golang.org/x/time/rate"
)

func main() {
	var ok bool
	var node = &core.Node{}
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
		p2pgo.ProtocolPrefix("/kldr-testnet"),
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	node, ok = h1.(*core.Node)
	if !ok {
		panic(err)
	}

	fmt.Println(node.Addrs(), node.ID())

	node.RouteTableFindPeers(0)

	tick := time.NewTicker(time.Second * 30)
	defer tick.Stop()

	var r = rate.Every(time.Second * 3)
	var limit = rate.NewLimiter(r, 1)

	for {
		select {
		case peer, ok := <-node.GetDiscoveredPeers():
			if !ok {
				break
			}
			if limit.Allow() {
				tick.Reset(time.Second * 30)
			}
			if len(peer.Responses) == 0 {
				break
			}
			for _, v := range peer.Responses {
				log.Println("found: ", v.ID.Pretty(), v.Addrs)
			}
		case <-tick.C:
			node.RouteTableFindPeers(0)
		}
	}
}
```

## License
Licensed under [Apache 2.0](https://github.com/CESSProject/p2p-go/blob/main/LICENSE)
