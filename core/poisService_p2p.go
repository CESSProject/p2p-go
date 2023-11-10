/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"time"

	"github.com/AstaFrode/go-libp2p/core/peer"
	"github.com/CESSProject/p2p-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (n *Node) PoisServiceRequestGenTagP2P(
	peerid peer.ID,
	fileData []byte,
	filehash string,
	customData string,
	timeout time.Duration,
) (*pb.ResponseGenTag, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)
	result, err := c.RequestGenTag(ctx, &pb.RequestGenTag{
		FileData:   fileData,
		Name:       filehash,
		CustomData: customData,
	})
	return result, err
}

func (n *Node) PoisServiceRequestBatchVerifyP2P(
	peerid peer.ID,
	names []string,
	us []string,
	mus []string,
	sigma string,
	minerPeerid []byte,
	minerPbk []byte,
	minerPeerIdSign []byte,
	qslices *pb.RequestBatchVerify_Qslice,
	timeout time.Duration,
) (*pb.ResponseBatchVerify, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	result, err := c.RequestBatchVerify(ctx, &pb.RequestBatchVerify{
		AggProof: &pb.RequestBatchVerify_BatchVerifyParam{
			Names: names,
			Us:    us,
			Mus:   mus,
			Sigma: sigma,
		},
		PeerId:          minerPeerid,
		MinerPbk:        minerPbk,
		MinerPeerIdSign: minerPeerIdSign,
		Qslices:         qslices,
	})
	return result, err
}
