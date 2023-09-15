package core

import (
	"context"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (n *Node) PoisGetMinerInitParamP2P(peerid peer.ID, accountKey []byte, timeout time.Duration) (*pb.ResponseMinerInitParam, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	result, err := c.RequestMinerGetNewKey(ctx, &pb.RequestMinerInitParam{
		MinerId: accountKey,
	})
	return result, err
}

func (n *Node) PoisMinerCommitGenChallP2P(
	peerid peer.ID,
	commitGenChall *pb.RequestMinerCommitGenChall,
	timeout time.Duration,
) (*pb.Challenge, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	result, err := c.RequestMinerCommitGenChall(ctx, commitGenChall)
	return result, err
}

func (n *Node) PoisVerifyCommitProofP2P(
	peerid peer.ID,
	verifyCommitAndAccProof *pb.RequestVerifyCommitAndAccProof,
	timeout time.Duration,
) (*pb.ResponseVerifyCommitOrDeletionProof, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	result, err := c.RequestVerifyCommitProof(ctx, verifyCommitAndAccProof)
	return result, err
}

func (n *Node) PoisSpaceProofVerifySingleBlockP2P(
	peerid peer.ID,
	accountKey []byte,
	spaceChals []int64,
	minerPoisInfo *pb.MinerPoisInfo,
	proof *pb.SpaceProof,
	spaceProofHashPolkadotSig []byte,
	timeout time.Duration,
) (*pb.ResponseSpaceProofVerify, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	result, err := c.RequestSpaceProofVerifySingleBlock(ctx, &pb.RequestSpaceProofVerify{
		SpaceChals:                     spaceChals,
		MinerId:                        accountKey,
		PoisInfo:                       minerPoisInfo,
		Proof:                          proof,
		MinerSpaceProofHashPolkadotSig: spaceProofHashPolkadotSig,
	})
	return result, err
}

func (n *Node) PoisRequestVerifySpaceTotalP2P(
	peerid peer.ID,
	accountKey []byte,
	proofList []*pb.BlocksProof,
	front int64,
	rear int64,
	acc []byte,
	spaceChals []int64,
	timeout time.Duration,
) (*pb.ResponseSpaceProofVerifyTotal, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

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

func (n *Node) PoisRequestVerifyDeletionProofP2P(
	peerid peer.ID,
	requestVerifyDeletionProof *pb.RequestVerifyDeletionProof,
	timeout time.Duration,
) (*pb.ResponseVerifyCommitOrDeletionProof, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := n.libp2pgrpcCli.Dial(ctx, peerid, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewPoisApiClient(conn)

	result, err := c.RequestVerifyDeletionProof(ctx, requestVerifyDeletionProof)
	return result, err
}
