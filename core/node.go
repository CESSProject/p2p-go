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
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	ma "github.com/multiformats/go-multiaddr"
)

// P2P is an object participating in a p2p network, which
// implements protocols or provides services. It handles
// requests like a Server, and issues requests like a Client.
// It is called Host because it is both Server and Client (and Peer
// may be confusing).
// It references libp2p: https://github.com/libp2p/go-libp2p
type P2P interface {
	host.Host // lib-p2p host
}

// Node type - Implementation of a P2P Host
type Node struct {
	ctx            context.Context
	ctxCancel      context.CancelFunc
	host           host.Host // lib-p2p host
	workspace      string    // data
	privatekeyPath string
	multiaddr      string
	idleFileTee    string
	serviceFileTee string
	FileDir        string
	TmpDir         string
	IdleDir        string
	TagDir         string
	ProofDir       string
	IproofFile     string
	IproofMuFile   string
	SproofFile     string
	SproofMuFile   string
	idleDataCh     chan string
	tagDataCh      chan string
}

// NewBasicNode constructs a new *Node
//
//	  multiaddr: listen addresses of p2p host
//	  workspace: service working directory
//	  privatekeypath: private key file
//		  If it's empty, automatically created in the program working directory
//		  If it's a directory, it will be created in the specified directory
func NewBasicNode(multiaddr ma.Multiaddr, workspace string, privatekeypath string, bootpeers []string, cmgr connmgr.ConnManager) (*Node, error) {
	ctx := context.Background()
	if multiaddr == nil || workspace == "" {
		return nil, errors.New("invalid parameter")
	}

	prvKey, err := identify(workspace, privatekeypath)
	if err != nil {
		return nil, err
	}

	ip, port, err := ExtractIp4FromMultiaddr(multiaddr.String())
	if err != nil {
		return nil, err
	}

	publicip := ip
	if ip == AllIpAddress || ip == "" {
		publicip, err = GetExternalIp()
		if err != nil {
			return nil, err
		}
	}

	host, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(prvKey),
		// Multiple listen addresses
		libp2p.ListenAddrs(multiaddr),
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(cmgr),
		// Support TLS connections
		//libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		libp2p.DefaultMuxers,
	)
	if err != nil {
		return nil, err
	}

	if !host.ID().MatchesPrivateKey(prvKey) {
		return nil, errors.New("")
	}

	if len(bootpeers) > 0 {
		// Construct a datastore (needed by the DHT)
		dstore := dsync.MutexWrap(ds.NewMapDatastore())

		// Make the DHT
		dhtTable := dht.NewDHT(ctx, host, dstore)

		// Make the routed host
		routedHost := rhost.Wrap(host, dhtTable)

		err = bootstrapConnect(ctx, routedHost, convertPeers(bootpeers))
		if err != nil {
			return nil, err
		}

		// Bootstrap the host
		err = dhtTable.Bootstrap(ctx)
		if err != nil {
			return nil, err
		}
	}

	basicCtx, cancel := context.WithCancel(context.Background())
	n := &Node{
		ctx:            basicCtx,
		ctxCancel:      cancel,
		host:           host,
		workspace:      workspace,
		privatekeyPath: privatekeypath,
		multiaddr:      fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", publicip, port, host.ID()),
		FileDir:        filepath.Join(workspace, FileDirectionry),
		TmpDir:         filepath.Join(workspace, TmpDirectionry),
		IdleDir:        filepath.Join(workspace, IdleDirectionry),
		TagDir:         filepath.Join(workspace, TagDirectionry),
		ProofDir:       filepath.Join(workspace, ProofDirectionry),
		IproofFile:     filepath.Join(workspace, ProofDirectionry, IdleProofFile),
		IproofMuFile:   filepath.Join(workspace, ProofDirectionry, IdleMuFile),
		SproofFile:     filepath.Join(workspace, ProofDirectionry, ServiceProofFile),
		SproofMuFile:   filepath.Join(workspace, ProofDirectionry, ServiceMuFile),
		idleDataCh:     make(chan string, 1),
		tagDataCh:      make(chan string, 1),
	}

	if err := mkdir(workspace); err != nil {
		return nil, err
	}

	//n.StarFileTransferProtocol()
	return n, nil
}

func (n *Node) AddAddrToPearstore(id peer.ID, addr ma.Multiaddr, t time.Duration) {
	time := peerstore.RecentlyConnectedAddrTTL
	if t.Seconds() > 0 {
		time = t
	}
	n.Peerstore().AddAddr(id, addr, time)
}

func (n *Node) AddMultiaddrToPearstore(multiaddr string, t time.Duration) (peer.ID, error) {
	time := peerstore.RecentlyConnectedAddrTTL
	if t.Seconds() > 0 {
		time = t
	}

	maddr, err := ma.NewMultiaddr(multiaddr)
	if err != nil {
		return "", err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return "", err
	}

	n.Peerstore().AddAddr(info.ID, maddr, time)
	return info.ID, nil
}

func (n *Node) PrivatekeyPath() string {
	return n.privatekeyPath
}

func (n *Node) Workspace() string {
	return n.workspace
}

func (n *Node) Multiaddr() string {
	return n.multiaddr
}

func (n *Node) ID() peer.ID {
	return n.host.ID()
}

func (n *Node) Peerstore() peerstore.Peerstore {
	return n.host.Peerstore()
}

func (n *Node) Addrs() []ma.Multiaddr {
	return n.host.Addrs()
}

func (n *Node) Mux() protocol.Switch {
	return n.host.Mux()
}

func (n *Node) Close() error {
	return n.host.Close()
}

func (n *Node) Network() network.Network {
	return n.host.Network()
}

func (n *Node) ConnManager() connmgr.ConnManager {
	return n.host.ConnManager()
}

func (n *Node) Connect(ctx context.Context, pi peer.AddrInfo) error {
	return n.host.Connect(ctx, pi)
}

func (n *Node) EventBus() event.Bus {
	return n.host.EventBus()
}

func (n *Node) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	n.host.SetStreamHandler(pid, handler)
}

func (n *Node) SetStreamHandlerMatch(pid protocol.ID, m func(protocol.ID) bool, handler network.StreamHandler) {
	n.host.SetStreamHandlerMatch(pid, m, handler)
}

func (n *Node) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, p, pids...)
}

func (n *Node) RemoveStreamHandler(pid protocol.ID) {
	n.host.RemoveStreamHandler(pid)
}

func (n *Node) PutIdleDataEventCh(path string) {
	go func() {
		if len(n.idleDataCh) > 0 {
			_ = <-n.idleDataCh
		}
		n.idleDataCh <- path
	}()
}

func (n *Node) GetIdleDataEvent() chan string {
	return n.idleDataCh
}

func (n *Node) PutTagEventCh(path string) {
	go func() {
		if len(n.tagDataCh) > 0 {
			_ = <-n.tagDataCh
		}
		n.tagDataCh <- path
	}()
}

func (n *Node) GetTagEvent() chan string {
	return n.tagDataCh
}

func (n *Node) SetIdleFileTee(peerid string) {
	n.idleFileTee = peerid
}

func (n *Node) GetIdleFileTee() string {
	return n.idleFileTee
}

func (n *Node) SetServiceFileTee(peerid string) {
	n.serviceFileTee = peerid
}

func (n *Node) GetServiceFileTee() string {
	return n.serviceFileTee
}

// identify reads or creates the private key file specified by fpath
func identify(workspace, fpath string) (crypto.PrivKey, error) {
	if fpath == "" {
		fpath = filepath.Join(workspace, privatekeyFile)
	}

	fstat, err := os.Stat(fpath)
	if err == nil {
		if fstat.IsDir() {
			fpath = filepath.Join(fpath, privatekeyFile)
		} else {
			content, err := os.ReadFile(fpath)
			if err != nil {
				return nil, err
			}
			return crypto.UnmarshalEd25519PrivateKey(content)
		}
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, rand.Reader)
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

func mkdir(workspace string) error {
	if err := os.MkdirAll(filepath.Join(workspace, FileDirectionry), DirMode); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(workspace, TmpDirectionry), DirMode); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(workspace, IdleDirectionry), DirMode); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(workspace, TagDirectionry), DirMode); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(workspace, ProofDirectionry), DirMode); err != nil {
		return err
	}
	return nil
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
	return n.verifyData(bin, sign, peerId, data.NodePubKey)
}

// sign an outgoing p2p message payload
func (n *Node) SignProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return n.signData(data)
}

// sign binary data using the local node's private key
func (n *Node) signData(data []byte) ([]byte, error) {
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

// NewPeerStream
func (n *Node) NewPeerStream(id peer.ID, p protocol.ID) (network.Stream, error) {
	return n.host.NewStream(context.Background(), id, p)
}

// SendMsgToStream
func (n *Node) SendMsgToStream(s network.Stream, msg []byte) error {
	_, err := s.Write(msg)
	return err
}

func ExtractIp4FromMultiaddr(multiaddr string) (string, uint64, error) {
	var ip string
	temp := strings.TrimPrefix(multiaddr, "/ip4/")
	temps := strings.Split(temp, "/")
	if !isIPv4(temps[0]) {
		return "", 0, fmt.Errorf("Invalid ip")
	}
	ip = temps[0]
	temp = strings.TrimPrefix(multiaddr, fmt.Sprintf("/ip4/%s/tcp/", temps[0]))
	temps = strings.Split(temp, "/")
	port, err := strconv.ParseUint(temps[0], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("Invalid port")
	}
	return ip, port, nil
}

// IsIPv4 is used to determine whether ipAddr is an ipv4 address
func isIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

// IsIPv6 is used to determine whether ipAddr is an ipv6 address
func isIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

// Get external network ip
func GetExternalIp() (string, error) {
	var (
		err        error
		externalIp string
	)

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	resp, err := client.Get("http://myexternalip.com/raw")
	if err == nil {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		externalIp = fmt.Sprintf("%s", string(b))
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx1, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err := exec.CommandContext(ctx1, "bash", "-c", "curl ifconfig.co").Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err = exec.CommandContext(ctx2, "bash", "-c", "curl cip.cc | grep  IP | awk '{print $3;}'").Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx3, _ := context.WithTimeout(context.Background(), 10*time.Second)
	output, err = exec.CommandContext(ctx3, "bash", "-c", `curl ipinfo.io | grep \"ip\" | awk '{print $2;}'`).Output()
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\"", "")
		externalIp = strings.ReplaceAll(externalIp, ",", "")
		externalIp = strings.ReplaceAll(externalIp, "\n", "")
		if isIPv4(externalIp) {
			return externalIp, nil
		}
	}
	return "", errors.New("Please check your network status")
}
