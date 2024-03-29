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

func (n *PeerNode) NewPodr2ApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2ApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPodr2ApiClient(conn), nil
}

func (n *PeerNode) NewPodr2VerifierApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2VerifierApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPodr2VerifierApiClient(conn), nil
}

func (n *PeerNode) RequestGenTag(c pb.Podr2ApiClient) (pb.Podr2Api_RequestGenTagClient, error) {
	result, err := c.RequestGenTag(context.Background())
	return result, err
}

func (n *PeerNode) RequestEcho(
	addr string,
	echoMessage *pb.EchoMessage,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.EchoMessage, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.Echo(ctx, echoMessage, callOpts...)
	return result, err
}

func (n *PeerNode) RequestBatchVerify(
	addr string,
	requestBatchVerify *pb.RequestBatchVerify,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.ResponseBatchVerify, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := pb.NewPodr2VerifierApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestBatchVerify(ctx, requestBatchVerify, callOpts...)
	return result, err
}
