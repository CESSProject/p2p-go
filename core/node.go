/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

type P2P interface {
	host.Host // lib-p2p host
	PrivatekeyPath() string
	WriteFile(multiaddr string, path string) error
	ReadFile(multiaddr string, rootHash, datahash, path string, size int64) error
}

// Node type - a p2p host implementing one or more p2p protocols
type Node struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	// ensures we shutdown ONLY once
	closeSync sync.Once

	host           host.Host // lib-p2p host
	workspace      string    // data
	privatekeyPath string
	protocol       map[string]interface{}
	// *myprotocol.WriteFileProtocol // writefile protocol impl
	// *myprotocol.ReadFileProtocol  // readfile protocol impl
	// add other protocols here...
}

// node - p2p node instance
// var node *Node

// func GetNode() *Node {
// 	return node
// }

func (n *Node) AddAddrToPearstore(id peer.ID, addr ma.Multiaddr) {
	n.host.Peerstore().AddAddr(id, addr, peerstore.AddressTTL)
}

func (n *Node) PrivatekeyPath() string {
	return n.privatekeyPath
}

func (n *Node) Workspace() string {
	return n.workspace
}

func (n *Node) AddProtocol(pid protocol.ID, handler network.StreamHandler) {
	n.host.SetStreamHandler(pid, handler)
}

// StartPeer configures and starts the p2p node service
// ip: listening IP address
// port: listening IP port
// datadir: where to save the key
//
// This function must be called before use
func StartPeer(ip string, port int, datadir string) error {
	// if node == nil {
	// 	newHost, err := newHost(ip, port, datadir)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	node = &Node{Host: newHost}
	// 	node.string = datadir
	// 	node.WriteFileProtocol = NewWriteFileProtocol(node)
	// 	node.ReadFileProtocol = NewReadFileProtocol(node)
	// 	return err
	// }
	return nil
}

func newHost(ip string, port int, datadir string) (host.Host, error) {
	privfileName := filepath.Join(datadir, privatekeyFile)
	prvKey, err := identify(privfileName)
	if err != nil {
		return nil, err
	}
	listen, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, port))
	if err != nil {
		return nil, err
	}
	host, err := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Identity(prvKey),
		yamuxOpt,
		mplexOpt,
	)
	if err != nil {
		return nil, err
	}
	if !host.ID().MatchesPrivateKey(prvKey) {
		return nil, errors.New("")
	}
	return host, err
}

func identify(fpath string) (crypto.PrivKey, error) {
	_, err := os.Stat(fpath)
	if err == nil {
		content, err := os.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		return crypto.UnmarshalSecp256k1PrivateKey(content)
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Secp256k1, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := prvKey.Raw()
	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}

	err = f.Sync()
	if err != nil {
		return nil, err
	}

	return prvKey, nil
}

// Authenticate incoming p2p message
// message: a protobufs go data object
// data: common p2p message data
func (n *Node) AuthenticateMessage(message proto.Message, data *pb.MessageData) bool {
	// store a temp ref to signature and remove it from message data
	// sign is a string to allow easy reset to zero-value (empty string)
	sign := data.Sign
	data.Sign = nil

	// marshall data without the signature to protobufs3 binary format
	bin, err := proto.Marshal(message)
	if err != nil {
		log.Println(err, "failed to marshal pb message")
		return false
	}

	// restore sig in message data (for possible future use)
	data.Sign = sign

	// restore peer id binary format from base58 encoded node id data
	peerId, err := peer.Decode(data.NodeId)
	if err != nil {
		log.Println(err, "Failed to decode node id from base58")
		return false
	}

	// verify the data was authored by the signing peer identified by the public key
	// and signature included in the message
	return n.verifyData(bin, []byte(sign), peerId, data.NodePubKey)
}

// sign an outgoing p2p message payload
func (n *Node) SignProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return n.SignData(data)
}

// sign binary data using the local node's private key
func (n *Node) SignData(data []byte) ([]byte, error) {
	key := n.host.Peerstore().PrivKey(n.host.ID())
	res, err := key.Sign(data)
	return res, err
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerId: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (n *Node) verifyData(data []byte, signature []byte, peerId peer.ID, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Println(err, "Failed to extract key from message key data")
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		log.Println(err, "Failed to extract peer id from public key")
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerId {
		log.Println(err, "Node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Println(err, "Error authenticating data")
		return false
	}

	return res
}

// helper method - generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func (n *Node) NewMessageData(messageId string, gossip bool) *pb.MessageData {
	// Add protobufs bin data for message author public key
	// this is useful for authenticating  messages forwarded by a node authored by another node
	nodePubKey, err := crypto.MarshalPublicKey(n.host.Peerstore().PubKey(n.host.ID()))

	if err != nil {
		panic("Failed to get public key for sender from local peer store.")
	}

	return &pb.MessageData{
		ClientVersion: p2pVersion,
		NodeId:        n.host.ID().String(),
		NodePubKey:    nodePubKey,
		Timestamp:     time.Now().Unix(),
		Id:            messageId,
		Gossip:        gossip,
	}
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *Node) SendProtoMessage(id peer.ID, p protocol.ID, data proto.Message) error {
	s, err := n.host.NewStream(context.Background(), id, p)
	if err != nil {
		return err
	}
	defer s.Close()

	writer := ggio.NewFullWriter(s)
	err = writer.WriteMsg(data)
	if err != nil {
		s.Reset()
		return err
	}
	return nil
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *Node) SendProtoMessageToStream(s network.Stream, data proto.Message) error {
	writer := ggio.NewFullWriter(s)
	err := writer.WriteMsg(data)
	if err != nil {
		s.Reset()
		return err
	}
	return nil
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *Node) NewStream(ctx context.Context, id peer.ID, p protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, id, p)
}
