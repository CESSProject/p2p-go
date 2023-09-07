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

func (n *Node) PoisNewClient(addr string, opts ...grpc.DialOption) (pb.PoisApiClient, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewPoisApiClient(conn), nil
}

func (n *Node) PoisGetMinerInitParam(addr string, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerInitParam, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
	return result, err
}

func (n *Node) PoisMinerCommitGenChall(
	addr string,
	accountKey []byte,
	commit *pb.Commits,
	minerPoisInfo *pb.MinerPoisInfo,
	minerSign []byte,
	timeout time.Duration,
) (*pb.Challenge, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestMinerCommitGenChall(ctx, &pb.RequestMinerCommitGenChall{
		MinerId:   accountKey,
		Commit:    commit,
		PoisInfo:  minerPoisInfo,
		MinerSign: minerSign,
	})
	return result, err
}

func (n *Node) PoisVerifyCommitProof(
	addr string,
	accountKey []byte,
	commitProofGroup *pb.CommitProofGroup,
	accProof *pb.AccProof,
	minerSign []byte,
	timeout time.Duration,
) (*pb.ResponseVerifyCommitOrDeletionProof, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
		MinerSign:        minerSign,
	})
	return result, err
}

func (n *Node) PoisSpaceProofVerifySingleBlock(
	addr string,
	accountKey []byte,
	spaceChals []int64,
	minerPoisInfo *pb.MinerPoisInfo,
	proof *pb.SpaceProof,
	spaceProofHashPolkadotSig []byte,
	timeout time.Duration,
) (*pb.ResponseSpaceProofVerify, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
		PoisInfo:                       minerPoisInfo,
		Proof:                          proof,
		MinerSpaceProofHashPolkadotSig: spaceProofHashPolkadotSig,
	})
	return result, err
}

func (n *Node) PoisRequestVerifySpaceTotal(
	addr string,
	accountKey []byte,
	proofList []*pb.BlocksProof,
	front int64,
	rear int64,
	acc []byte,
	spaceChals []int64,
	timeout time.Duration,
) (*pb.ResponseSpaceProofVerifyTotal, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestVerifySpaceTotal(ctx, &pb.RequestSpaceProofVerifyTotal{
		MinerId:    accountKey,
		ProofList:  proofList,
		Front:      front,
		Rear:       rear,
		Acc:        acc,
		SpaceChals: spaceChals,
	})
	return result, err
}

func (n *Node) PoisRequestVerifyDeletionProof(
	addr string,
	roots [][]byte,
	witChain *pb.AccWitnessNode,
	accPath [][]byte,
	minerId []byte,
	minerPoisInfo *pb.MinerPoisInfo,
	minerSign []byte,
	timeout time.Duration,
) (*pb.ResponseVerifyCommitOrDeletionProof, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := c.RequestVerifyDeletionProof(ctx, &pb.RequestVerifyDeletionProof{
		Roots:     roots,
		WitChain:  witChain,
		AccPath:   accPath,
		MinerId:   minerId,
		PoisInfo:  minerPoisInfo,
		MinerSign: minerSign,
	})
	return result, err
}
