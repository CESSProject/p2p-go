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
	"google.golang.org/grpc/credentials/insecure"
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
	blockNum uint64,
	filehash string,
	customData string,
	timeout time.Duration,
) (*pb.ResponseGenTag, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(128*1024*1024),
			grpc.MaxCallSendMsgSize(128*1024*1024),
		),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestGenTag(ctx, &pb.RequestGenTag{
		FileData:   fileData,
		BlockNum:   blockNum,
		Name:       filehash,
		CustomData: customData,
	})
	return result, err
}

func (n *Node) PoisServiceRequestBatchVerify(
	addr string,
	names []string,
	us []string,
	mus []string,
	sigma string,
	peerid []byte,
	minerPbk []byte,
	minerPeerIdSign []byte,
	qslices *pb.RequestBatchVerify_Qslice,
	timeout time.Duration,
) (*pb.ResponseBatchVerify, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(128*1024*1024),
			grpc.MaxCallSendMsgSize(128*1024*1024),
		),
	)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestBatchVerify(ctx, &pb.RequestBatchVerify{
		AggProof: &pb.RequestBatchVerify_BatchVerifyParam{
			Names: names,
			Us:    us,
			Mus:   mus,
			Sigma: sigma,
		},
		PeerId:          peerid,
		MinerPbk:        minerPbk,
		MinerPeerIdSign: minerPeerIdSign,
		Qslices:         qslices,
	})
	return result, err
}
