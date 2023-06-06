/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"

	p2pgo "github.com/CESSProject/p2p-go"
)

func main() {
	ctx := context.Background()
	h1, err := p2pgo.New(
		ctx,
		p2pgo.PrivatekeyFile(".private1"),
		p2pgo.ListenPort(8080),
		p2pgo.Workspace("."),
	)
	if err != nil {
		panic(err)
	}
	defer h1.Close()

	fmt.Println(h1.Addrs(), h1.ID())
}
