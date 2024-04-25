/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/out"
	"github.com/CESSProject/p2p-go/pb"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"

	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	tls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// P2P is an object participating in a p2p network, which
// implements protocols or provides services. It handles
// requests like a Server, and issues requests like a Client.
// It is called Host because it is both Server and Client (and Peer
// may be confusing).
// It references libp2p: https://github.com/libp2p/go-libp2p
type P2P interface {
	// Lib-p2p host
	host.Host

	// Message protocol
	Protocol

	// PrivatekeyPath returns the key file location
	PrivatekeyPath() string

	// Workspace returns to the working directory
	Workspace() string

	// GetPeerPublickey returns the host's public key
	GetPeerPublickey() []byte

	// GetProtocolVersion returns the ProtocolVersion of the host
	GetProtocolVersion() string

	//
	GetProtocolPrefix() string

	// GetDhtProtocolVersion returns the host's DHT ProtocolVersion
	GetDhtProtocolVersion() string

	// GetRendezvousVersion returns the rendezvous protocol
	GetRendezvousVersion() string

	// GetDirs returns the data directory structure of the host
	GetDirs() DataDirs

	// GetBootnode returns bootnode
	GetBootnode() string

	// GetNetEnv returns network env
	GetNetEnv() string

	// SetBootnode updates the host's boot node
	SetBootnode(bootnode string)

	//
	GetHost() host.Host

	//
	GetDHTable() *dht.IpfsDHT

	//
	EnableRecv()

	//
	DisableRecv()

	//
	GetRecvFlag() bool

	// Close p2p
	Close() error

	// grpc api

	NewPoisCertifierApiClient(addr string, opts ...grpc.DialOption) (pb.PoisCertifierApiClient, error)

	NewPoisVerifierApiClient(addr string, opts ...grpc.DialOption) (pb.PoisVerifierApiClient, error)

	NewPodr2ApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2ApiClient, error)

	NewPodr2VerifierApiClient(addr string, opts ...grpc.DialOption) (pb.Podr2VerifierApiClient, error)

	NewPubkeyApiClient(addr string, opts ...grpc.DialOption) (pb.CesealPubkeysProviderClient, error)

	RequestMinerGetNewKey(
		addr string,
		accountKey []byte,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseMinerInitParam, error)

	RequestMinerCommitGenChall(
		addr string,
		commitGenChall *pb.RequestMinerCommitGenChall,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.Challenge, error)

	RequestVerifyCommitProof(
		addr string,
		verifyCommitAndAccProof *pb.RequestVerifyCommitAndAccProof,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseVerifyCommitOrDeletionProof, error)

	RequestVerifyDeletionProof(
		addr string,
		requestVerifyDeletionProof *pb.RequestVerifyDeletionProof,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseVerifyCommitOrDeletionProof, error)

	RequestSpaceProofVerifySingleBlock(
		addr string,
		requestSpaceProofVerify *pb.RequestSpaceProofVerify,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseSpaceProofVerify, error)

	RequestVerifySpaceTotal(
		addr string,
		requestSpaceProofVerifyTotal *pb.RequestSpaceProofVerifyTotal,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseSpaceProofVerifyTotal, error)

	RequestGenTag(c pb.Podr2ApiClient) (pb.Podr2Api_RequestGenTagClient, error)

	RequestEcho(
		addr string,
		echoMessage *pb.EchoMessage,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.EchoMessage, error)

	RequestBatchVerify(
		addr string,
		requestBatchVerify *pb.RequestBatchVerify,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.ResponseBatchVerify, error)

	GetIdentityPubkey(
		addr string,
		request *pb.Request,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.IdentityPubkeyResponse, error)

	GetMasterPubkey(
		addr string,
		request *pb.Request,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.MasterPubkeyResponse, error)

	GetPodr2Pubkey(
		addr string,
		request *pb.Request,
		timeout time.Duration,
		dialOpts []grpc.DialOption,
		callOpts []grpc.CallOption,
	) (*pb.Podr2PubkeyResponse, error)
}

// Node type - Implementation of a P2P Host
type PeerNode struct {
	host               host.Host
	dir                DataDirs
	peerPublickey      []byte
	workspace          string
	privatekeyPath     string
	idleTee            atomic.Value
	serviceTee         atomic.Value
	protocolVersion    string
	dhtProtocolVersion string
	rendezvousVersion  string
	protocolPrefix     string
	enableRecv         bool
	bootnode           string
	netenv             string
	dhtable            *dht.IpfsDHT
	*protocols
}

var _ P2P = (*PeerNode)(nil)

// NewPeerNode constructs a new *PeerNode
//
//	  workspace: service working directory
//	  privatekeypath: private key file
//		  If it is empty, automatically created in the program working directory
//		  If it is a directory, it will be created in the specified directory
func NewPeerNode(ctx context.Context, cfg *config.Config) (*PeerNode, error) {
	if !FreeLocalPort(uint32(cfg.ListenPort)) {
		return nil, fmt.Errorf("port %d is already in use", cfg.ListenPort)
	}

	err := verifyWorkspace(cfg.Workspace)
	if err != nil {
		return nil, err
	}

	prvKey, err := identification(cfg.Workspace, cfg.PrivatekeyPath)
	if err != nil {
		return nil, err
	}

	if cfg.PublicIpv4 != "" {
		if !IsIPv4(cfg.PublicIpv4) {
			return nil, fmt.Errorf("illegal IPv4 address: %s", cfg.PublicIpv4)
		}
	}

	var opts []libp2p.Option
	var multiaddrs []ma.Multiaddr
	externalIp, err := GetExternalIp()
	if err != nil {
		if cfg.PublicIpv4 != "" {
			extMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.PublicIpv4, cfg.ListenPort))
			multiaddrs = append(multiaddrs, extMultiAddr)
		}
	} else {
		extMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", externalIp, cfg.ListenPort))
		multiaddrs = append(multiaddrs, extMultiAddr)
	}

	if len(multiaddrs) > 0 {
		opts = append(opts, libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
			addrs = append(addrs, multiaddrs...)
			return addrs
		}))
	}

	var boots = make([]string, 0)
	for _, b := range cfg.BootPeers {
		bootnodes, err := ParseMultiaddrs(b)
		if err != nil {
			out.Err(err.Error())
			continue
		}
		boots = append(boots, bootnodes...)
	}
	if len(boots) == 0 {
		rm, err := buildPrimaryResourceManager()
		if err == nil {
			opts = append(opts, libp2p.ResourceManager(rm))
		}
	}

	opts = append(opts,
		libp2p.Identity(prvKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.ListenPort)),
		libp2p.ConnectionManager(cfg.ConnManager),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Security(tls.ID, tls.New),
		libp2p.ProtocolVersion(cfg.ProtocolPrefix+p2pProtocolVer),
		libp2p.DisableMetrics(),
		libp2p.EnableRelay(),
	)

	bhost, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("[libp2p.New] %v", err)
	}

	if !bhost.ID().MatchesPrivateKey(prvKey) {
		return nil, errors.New("invalid private key")
	}

	publickey, err := base58.Decode(bhost.ID().String())
	if err != nil {
		return nil, err
	}

	peer_node := &PeerNode{
		host:               bhost,
		workspace:          cfg.Workspace,
		privatekeyPath:     cfg.PrivatekeyPath,
		peerPublickey:      publickey,
		protocolPrefix:     cfg.ProtocolPrefix,
		protocolVersion:    cfg.ProtocolPrefix + p2pProtocolVer,
		dhtProtocolVersion: cfg.ProtocolPrefix + dhtProtocolVer,
		rendezvousVersion:  cfg.ProtocolPrefix + rendezvous,
		enableRecv:         true,
		protocols:          NewProtocol(),
	}

	peer_node.dhtable, peer_node.bootnode, peer_node.netenv, err = NewDHT(ctx, bhost, cfg.BucketSize, cfg.Version, boots, cfg.ProtocolPrefix, peer_node.dhtProtocolVersion)
	if err != nil {
		return nil, fmt.Errorf("[NewDHT] %v", err)
	}

	if len(boots) > 0 {
		peer_node.dir, err = mkdir(cfg.Workspace)
		if err != nil {
			return nil, err
		}
		peer_node.initProtocol(cfg.ProtocolPrefix)
		bootstrapAddr, _ := ma.NewMultiaddr(peer_node.bootnode)
		peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrapAddr)
		peer_node.OnlineAction(peerinfo.ID)
	} else {
		peer_node.OnlineProtocol = peer_node.NewOnlineProtocol()
	}

	return peer_node, nil
}

// ID returns the (local) peer.ID associated with this Host
func (n *PeerNode) ID() peer.ID {
	return n.host.ID()
}

// Peerstore returns the Host's repository of Peer Addresses and Keys.
func (n *PeerNode) Peerstore() peerstore.Peerstore {
	return n.host.Peerstore()
}

// Returns the listen addresses of the Host
func (n *PeerNode) Addrs() []ma.Multiaddr {
	return n.host.Addrs()
}

// Networks returns the Network interface of the Host
func (n *PeerNode) Network() network.Network {
	return n.host.Network()
}

// Mux returns the Mux multiplexing incoming streams to protocol handlers
func (n *PeerNode) Mux() protocol.Switch {
	return n.host.Mux()
}

// Connect ensures there is a connection between this host and the peer with
// given peer.ID. Connect will absorb the addresses in pi into its internal
// peerstore. If there is not an active connection, Connect will issue a
// h.Network.Dial, and block until a connection is open, or an error is
// returned. // TODO: Relay + NAT.
func (n *PeerNode) Connect(ctx context.Context, pi peer.AddrInfo) error {
	return n.host.Connect(ctx, pi)
}

// SetStreamHandler sets the protocol handler on the Host's Mux.
// This is equivalent to: host.Mux().SetHandler(proto, handler) (Threadsafe)
func (n *PeerNode) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	n.host.SetStreamHandler(pid, handler)
}

// SetStreamHandlerMatch sets the protocol handler on the Host's Mux
// using a matching function for protocol selection.
func (n *PeerNode) SetStreamHandlerMatch(pid protocol.ID, m func(protocol.ID) bool, handler network.StreamHandler) {
	n.host.SetStreamHandlerMatch(pid, m, handler)
}

// RemoveStreamHandler removes a handler on the mux that was set by
// SetStreamHandler
func (n *PeerNode) RemoveStreamHandler(pid protocol.ID) {
	n.host.RemoveStreamHandler(pid)
}

// NewStream opens a new stream to given peer p, and writes a p2p/protocol
// header with given ProtocolID. If there is no connection to p, attempts
// to create one. If ProtocolID is "", writes no header.
// (Threadsafe)
func (n *PeerNode) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, p, pids...)
}

// Close shuts down the host, its Network, and services.
func (n *PeerNode) Close() error {
	err := n.host.Close()
	if err != nil {
		return err
	}
	return nil
}

// ConnManager returns this hosts connection manager
func (n *PeerNode) ConnManager() connmgr.ConnManager {
	return n.host.ConnManager()
}

// EventBus returns the hosts eventbus
func (n *PeerNode) EventBus() event.Bus {
	return n.host.EventBus()
}

func (n *PeerNode) PrivatekeyPath() string {
	return n.privatekeyPath
}

func (n *PeerNode) Workspace() string {
	return n.workspace
}

func (n *PeerNode) GetPeerPublickey() []byte {
	return n.peerPublickey
}

func (n *PeerNode) GetProtocolVersion() string {
	return n.protocolVersion
}

func (n *PeerNode) GetProtocolPrefix() string {
	return n.protocolPrefix
}

func (n *PeerNode) GetDhtProtocolVersion() string {
	return n.dhtProtocolVersion
}

func (n *PeerNode) GetRendezvousVersion() string {
	return n.rendezvousVersion
}

func (n *PeerNode) GetBootnode() string {
	return n.bootnode
}

func (n *PeerNode) SetBootnode(bootnode string) {
	n.bootnode = bootnode
}

func (n *PeerNode) GetHost() host.Host {
	return n.host
}

func (n *PeerNode) GetDHTable() *dht.IpfsDHT {
	return n.dhtable
}

func (n *PeerNode) GetDirs() DataDirs {
	return n.dir
}

func (n *PeerNode) GetNetEnv() string {
	return n.netenv
}

func (n *PeerNode) EnableRecv() {
	n.enableRecv = true
}

func (n *PeerNode) DisableRecv() {
	n.enableRecv = false
}

func (n *PeerNode) GetRecvFlag() bool {
	return n.enableRecv
}

func (n *PeerNode) SetIdleFileTee(peerid string) {
	n.idleTee.Store(peerid)
}

func (n *PeerNode) GetIdleFileTee() string {
	value, ok := n.idleTee.Load().(string)
	if !ok {
		return ""
	}
	return value
}

func (n *PeerNode) SetServiceFileTee(peerid string) {
	n.serviceTee.Store(peerid)
}

func (n *PeerNode) GetServiceFileTee() string {
	value, ok := n.serviceTee.Load().(string)
	if !ok {
		return ""
	}
	return value
}

// identify reads or creates the private key file specified by fpath
func identification(workspace, fpath string) (crypto.PrivKey, error) {
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

func mkdir(workspace string) (DataDirs, error) {
	dataDir := DataDirs{
		FileDir: filepath.Join(workspace, FileDataDirectionry),
		TmpDir:  filepath.Join(workspace, TmpDataDirectionry),
	}
	if err := os.MkdirAll(dataDir.FileDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.TmpDir, DirMode); err != nil {
		return dataDir, err
	}
	return dataDir, nil
}

// helper method - generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func (n *PeerNode) NewMessageData(messageId string, gossip bool) *pb.MessageData {
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
func (n *PeerNode) SendProtoMessage(id peer.ID, p protocol.ID, data proto.Message) error {
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
func (n *PeerNode) NewPeerStream(id peer.ID, p protocol.ID) (network.Stream, error) {
	return n.host.NewStream(context.Background(), id, p)
}

// SendMsgToStream
func (n *PeerNode) SendMsgToStream(s network.Stream, msg []byte) error {
	_, err := s.Write(msg)
	return err
}

func verifyWorkspace(ws string) error {
	if ws == "" {
		return fmt.Errorf("empty workspace")
	}
	fstat, err := os.Stat(ws)
	if err != nil {
		return os.MkdirAll(ws, DirMode)
	} else {
		if !fstat.IsDir() {
			return fmt.Errorf("workspace is not a directory")
		}
	}
	return nil
}

func (n *PeerNode) initProtocol(protocolPrefix string) {
	n.SetProtocolPrefix(protocolPrefix)
	n.WriteFileProtocol = n.NewWriteFileProtocol()
	n.ReadFileProtocol = n.NewReadFileProtocol()
	n.ReadDataProtocol = n.NewReadDataProtocol()
	n.ReadDataStatProtocol = n.NewReadDataStatProtocol()
	n.OnlineProtocol = n.NewOnlineProtocol()
}

func NewDHT(ctx context.Context, h host.Host, bucketsize int, version string, boot_nodes []string, protocolPrefix, dhtProtocol string) (*dht.IpfsDHT, string, string, error) {
	var options []dht.Option
	options = append(options,
		dht.ProtocolPrefix(protocol.ID(protocolPrefix)),
		dht.V1ProtocolOverride(protocol.ID(dhtProtocol)),
		dht.Resiliency(10),
		dht.BucketSize(bucketsize),
	)

	if len(boot_nodes) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h, options...)
	if err != nil {
		return nil, "", "", err
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, "", "", err
	}

	if len(boot_nodes) > 0 {
		netenv := ""
		for _, peerAddr := range boot_nodes {
			bootstrapAddr, err := ma.NewMultiaddr(peerAddr)
			if err != nil {
				continue
			}
			peerinfo, err := peer.AddrInfoFromP2pAddr(bootstrapAddr)
			if err != nil {
				continue
			}
			err = h.Connect(ctx, *peerinfo)
			if err == nil {
				out.Ok(fmt.Sprintf("Connect to the boot node: %s", peerinfo.ID.String()))
				switch peerinfo.ID.String() {
				case "12D3KooWS8a18xoBzwkmUsgGBctNo6QCr6XCpUDR946mTBBUTe83",
					"12D3KooWDWeiiqbpNGAqA5QbDTdKgTtwX8LCShWkTpcyxpRf2jA9",
					"12D3KooWNcTWWuUWKhjTVDF1xZ38yCoHXoF4aDjnbjsNpeVwj33U":
					netenv = "testnet"
				case "12D3KooWGDk9JJ5F6UPNuutEKSbHrTXnF5eSn3zKaR27amgU6o9S",
					"12D3KooWEGeAp1MvvUrBYQtb31FE1LPg7aHsd1LtTXn6cerZTBBd",
					"12D3KooWRm2sQg65y2ZgCUksLsjWmKbBtZ4HRRsGLxbN76XTtC8T":
					netenv = "devnet"
				default:
					netenv = "mainnet"
				}
				return kademliaDHT, bootstrapAddr.String(), netenv, nil
			}
		}
	} else {
		return kademliaDHT, "", "", nil
	}
	return kademliaDHT, "", "", fmt.Errorf("failed to connect to all boot nodes")
}

func buildPrimaryResourceManager() (network.ResourceManager, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.InfiniteLimits

	// Add limits around included libp2p protocols
	// libp2p.SetDefaultServiceLimits(&scalingLimits)

	// Turn the scaling limits into a concrete set of limits using `.AutoScale`. This
	// scales the limits proportional to your system memory.
	// scaledDefaultLimits := scalingLimits.AutoScale()

	// scaledDefaultLimits := scalingLimits.Scale(int64(memory.TotalMemory()/10*8), DefaultPrimaryFDCount)
	// // Tweak certain settings
	// cfg := rcmgr.PartialLimitConfig{
	// 	System: rcmgr.ResourceLimits{
	// 		// Allow unlimited outbound streams
	// 		Conns:           2,
	// 		StreamsOutbound: rcmgr.Unlimited,
	// 	},
	// 	// Everything else is default. The exact values will come from `scaledDefaultLimits` above.
	// }

	// // Create our limits by using our cfg and replacing the default values with values from `scaledDefaultLimits`
	// limits := cfg.Build(scaledDefaultLimits)

	// The resource manager expects a limiter, se we create one from our limits.
	limiter := rcmgr.NewFixedLimiter(scalingLimits)

	// Metrics are enabled by default. If you want to disable metrics, use the
	// WithMetricsDisabled option
	// Initialize the resource manager
	rm, err := rcmgr.NewResourceManager(limiter, rcmgr.WithMetricsDisabled())
	if err != nil {
		return nil, errors.Wrapf(err, "[NewResourceManager]")
	}
	return rm, nil
}
