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

func (n *Node) PoisNewClient(addr string) (pb.PoisApiClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewPoisApiClient(conn), nil
}

func (n *Node) PoisGetMinerInitParam(addr string, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerInitParam, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestMinerGetNewKey(ctx, &pb.RequestMinerInitParam{
		MinerId: accountKey,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisMinerRegister(addr string, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerRegister, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestMinerRegister(ctx, &pb.RequestMinerInitParam{
		MinerId: accountKey,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisMinerCommitGenChall(addr string, accountKey []byte, commit *pb.Commits, timeout time.Duration) (*pb.Challenge, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestMinerCommitGenChall(ctx, &pb.RequestMinerCommitGenChall{
		MinerId: accountKey,
		Commit:  commit,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisVerifyCommitProof(addr string, accountKey []byte, commitProofGroup *pb.CommitProofGroup, accProof *pb.AccProof, timeout time.Duration) (*pb.ResponseVerifyCommitOrDeletionProof, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestVerifyCommitProof(ctx, &pb.RequestVerifyCommitAndAccProof{
		CommitProofGroup: commitProofGroup,
		AccProof:         accProof,
		MinerId:          accountKey,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisSpaceProofVerifySingleBlock(
	addr string,
	accountKey []byte,
	spaceChals []int64,
	keyN []byte,
	keyG []byte,
	acc []byte,
	front int64,
	rear int64,
	proof *pb.SpaceProof,
	spaceProofHashPolkadotSig []byte,
	timeout time.Duration,
) (*pb.ResponseSpaceProofVerify, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestSpaceProofVerifySingleBlock(ctx, &pb.RequestSpaceProofVerify{
		SpaceChals:                     spaceChals,
		MinerId:                        accountKey,
		KeyN:                           keyN,
		KeyG:                           keyG,
		Acc:                            acc,
		Front:                          front,
		Rear:                           rear,
		Proof:                          proof,
		MinerSpaceProofHashPolkadotSig: spaceProofHashPolkadotSig,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
