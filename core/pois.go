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

func (n *Node) PoisNewClient(addr string, timeout time.Duration) (pb.PoisApiClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewPoisApiClient(conn), nil
}

func (n *Node) PoisGetMinerInitParam(cli pb.PoisApiClient, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerInitParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := cli.RequestMinerGetNewKey(ctx, &pb.RequestMinerInitParam{
		MinerId: accountKey,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisMinerRegister(cli pb.PoisApiClient, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerRegister, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := cli.RequestMinerRegister(ctx, &pb.RequestMinerInitParam{
		MinerId: accountKey,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisMinerCommitGenChall(cli pb.PoisApiClient, accountKey []byte, commit *pb.Commits, timeout time.Duration) (*pb.Challenge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := cli.RequestMinerCommitGenChall(ctx, &pb.RequestMinerCommitGenChall{
		MinerId: accountKey,
		Commit:  commit,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *Node) PoisVerifyCommitProof(cli pb.PoisApiClient, accountKey []byte, commitProofGroup *pb.CommitProofGroup, accProof *pb.AccProof, timeout time.Duration) (*pb.ResponseVerifyCommitAndAccProof, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := cli.RequestVerifyCommitProof(ctx, &pb.RequestVerifyCommitAndAccProof{
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
	cli pb.PoisApiClient,
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := cli.RequestSpaceProofVerifySingleBlock(ctx, &pb.RequestSpaceProofVerify{
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
