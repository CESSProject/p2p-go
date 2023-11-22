/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	"google.golang.org/grpc"
)

func (n *Node) PoisServiceNewClient(addr string, opts ...grpc.DialOption) (pb.Podr2ApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPodr2ApiClient(conn), nil
}

func (n *Node) PoisServiceRequestGenTag(
	addr string,
	fileData []byte,
	roothash string,
	filehash string,
	customData string,
	timeout time.Duration,
	opts ...grpc.DialOption,
) (*pb.ResponseGenTag, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestGenTag(ctx, &pb.RequestGenTag{
		FragmentData: fileData,
		FragmentName: filehash,
		CustomData:   customData,
		FileName:     roothash,
	})
	return result, err
}

func (n *Node) PoisServiceRequestBatchVerify(
	addr string,
	peerid []byte,
	minerPbk []byte,
	minerPeerIdSign []byte,
	batchVerifyParam *pb.RequestBatchVerify_BatchVerifyParam,
	qslices *pb.RequestBatchVerify_Qslice,
	timeout time.Duration,
	opts ...grpc.DialOption,
) (*pb.ResponseBatchVerify, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestBatchVerify(ctx, &pb.RequestBatchVerify{
		AggProof:        batchVerifyParam,
		PeerId:          peerid,
		MinerPbk:        minerPbk,
		MinerPeerIdSign: minerPeerIdSign,
		Qslices:         qslices,
	})
	return result, err
}
