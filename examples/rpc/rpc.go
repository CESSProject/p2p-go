/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"time"

	"github.com/CESSProject/p2p-go/core"
	"github.com/CESSProject/p2p-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var n = &core.Node{}
	result, err := n.RequestSpaceProofVerifySingleBlock(
		"127.0.0.1:10010",
		&pb.RequestSpaceProofVerify{},
		time.Minute,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		nil,
	)
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("result: ", result)
	}
}
