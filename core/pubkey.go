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

func (n *PeerNode) NewPubkeyApiClient(addr string, opts ...grpc.DialOption) (pb.CesealPubkeysProviderClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewCesealPubkeysProviderClient(conn), nil
}

func (n *PeerNode) GetIdentityPubkey(
	addr string,
	request *pb.Request,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.IdentityPubkeyResponse, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewCesealPubkeysProviderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.GetIdentityPubkey(ctx, request, callOpts...)
	return result, err
}

func (n *PeerNode) GetMasterPubkey(
	addr string,
	request *pb.Request,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.MasterPubkeyResponse, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewCesealPubkeysProviderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.GetMasterPubkey(ctx, request, callOpts...)
	return result, err
}

func (n *PeerNode) GetPodr2Pubkey(
	addr string,
	request *pb.Request,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.Podr2PubkeyResponse, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewCesealPubkeysProviderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.GetPodr2Pubkey(ctx, request, callOpts...)
	return result, err
}
