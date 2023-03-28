/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"fmt"
	"os"

	"github.com/CESSProject/p2p-go/node"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func WriteFile(multiaddr string, path string) error {
	fmt.Println("will write ", path, " to ", multiaddr)
	maddr, err := ma.NewMultiaddr(multiaddr)
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
	node.AddAddrToPearstore(info.ID, maddr)
	err = node.GetNode().WriteFileAction(info.ID, path)
	return err
}

func ReadFile(multiaddr string, hash string, path string) error {
	fmt.Println("will read ", path, " from ", multiaddr)
	maddr, err := ma.NewMultiaddr(multiaddr)
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
	node.AddAddrToPearstore(info.ID, maddr)
	err = node.GetNode().ReadFileAction(info.ID, hash, path)
	return err
}
