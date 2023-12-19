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

func (n *Node) NewPodr2ApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2ApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPodr2ApiClient(conn), nil
}

func (n *Node) NewPodr2VerifierApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2VerifierApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPodr2VerifierApiClient(conn), nil
}

func (n *Node) RequestGenTag(
	addr string,
	requestGenTag *pb.RequestGenTag,
	timeout time.Duration,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (*pb.ResponseGenTag, error) {
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPodr2ApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestGenTag(ctx, requestGenTag, callOpts...)
	return result, err
}

func (n *Node) RequestBatchVerify(
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
