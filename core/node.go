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

	"github.com/CESSProject/p2p-go/out"
	"github.com/CESSProject/p2p-go/pb"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"

	bsnet "github.com/AstaFrode/boxo/bitswap/network"

	bitswap "github.com/AstaFrode/boxo/bitswap"
	blockstore "github.com/AstaFrode/boxo/blockstore"
	dutil "github.com/AstaFrode/go-libp2p/p2p/discovery/util"
	blocks "github.com/ipfs/go-block-format"
	ds_sync "github.com/ipfs/go-datastore/sync"

	"github.com/AstaFrode/go-libp2p"
	dht "github.com/AstaFrode/go-libp2p-kad-dht"
	"github.com/AstaFrode/go-libp2p/core/connmgr"
	"github.com/AstaFrode/go-libp2p/core/crypto"
	"github.com/AstaFrode/go-libp2p/core/discovery"
	"github.com/AstaFrode/go-libp2p/core/event"
	"github.com/AstaFrode/go-libp2p/core/host"
	"github.com/AstaFrode/go-libp2p/core/network"
	"github.com/AstaFrode/go-libp2p/core/peer"
	"github.com/AstaFrode/go-libp2p/core/peerstore"
	"github.com/AstaFrode/go-libp2p/core/protocol"
	"github.com/AstaFrode/go-libp2p/core/routing"
	drouting "github.com/AstaFrode/go-libp2p/p2p/discovery/routing"
	rcmgr "github.com/AstaFrode/go-libp2p/p2p/host/resource-manager"
	"github.com/ipfs/go-cid"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"

	mh "github.com/multiformats/go-multihash"
	"github.com/pbnjay/memory"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
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

	// GetCtxRoot returns the root context of the host
	GetCtxRoot() context.Context

	// GetCtxCancelFromRoot returns the cancel context from root context
	GetCtxCancelFromRoot() context.Context

	// GetCtxQueryFromCtxCancel returns tne query context from cancel context
	GetCtxQueryFromCtxCancel() context.Context

	// PrivatekeyPath returns the key file location
	PrivatekeyPath() string

	// Workspace returns to the working directory
	Workspace() string

	// AddMultiaddrToPeerstore adds multiaddr to the host's peerstore
	AddMultiaddrToPeerstore(multiaddr string, t time.Duration) (peer.ID, error)

	// GetPeerPublickey returns the host's public key
	GetPeerPublickey() []byte

	// GetProtocolVersion returns the ProtocolVersion of the host
	GetProtocolVersion() string

	// GetDhtProtocolVersion returns the host's DHT ProtocolVersion
	GetDhtProtocolVersion() string

	// GetRendezvousVersion returns the rendezvous protocol
	GetRendezvousVersion() string

	// GetProtocolPrefix returns protocols prefix
	GetProtocolPrefix() string

	// GetDirs returns the data directory structure of the host
	GetDirs() DataDirs

	// GetBootstraps returns a list of host bootstraps
	GetBootstraps() []string

	// SetBootstraps updates the host's bootstrap list
	SetBootstraps(bootstrap []string)

	//
	GetDht() *dht.IpfsDHT

	//
	GetRoutingTable() *drouting.RoutingDiscovery

	// DHTFindPeer searches for a peer with given ID
	DHTFindPeer(peerid string) (peer.AddrInfo, error)

	// PeerID returns your own peerid
	PeerID() peer.ID

	//
	EnableRecv()

	//
	DisableRecv()

	// Close p2p
	Close() error

	//
	GetBlockstore() blockstore.Blockstore

	//
	GetBitSwap() *bitswap.Bitswap

	//
	FidToCid(fid string) (string, error)

	//
	SaveAndNotifyDataBlock(buf []byte) (cid.Cid, error)

	//
	NotifyData(buf []byte) error

	// GetDataFromBlock get data from block
	GetDataFromBlock(ctx context.Context, wantCid string) ([]byte, error)

	//
	GetLocalDataFromBlock(wantCid string) ([]byte, error)

	//
	GetDiscoveredPeers() <-chan *routing.QueryEvent

	// RouteTableFindPeers
	RouteTableFindPeers(limit int) (<-chan peer.AddrInfo, error)

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
type Node struct {
	ctxRoot               context.Context
	ctxCancelFromRoot     context.Context
	ctxQueryFromCtxCancel context.Context
	ctxCancelFuncFromRoot context.CancelFunc
	discoveredPeerCh      <-chan *routing.QueryEvent
	host                  host.Host
	bstore                blockstore.Blockstore
	bswap                 *bitswap.Bitswap
	dir                   DataDirs
	peerPublickey         []byte
	workspace             string
	privatekeyPath        string
	idleTee               atomic.Value
	serviceTee            atomic.Value
	protocolVersion       string
	dhtProtocolVersion    string
	rendezvousVersion     string
	protocolPrefix        string
	enableBswap           bool
	enableRecv            bool
	bootstrap             []string
	*dht.IpfsDHT
	*drouting.RoutingDiscovery
	*protocols
}

var _ P2P = (*Node)(nil)

// NewBasicNode constructs a new *Node
//
//	  workspace: service working directory
//	  privatekeypath: private key file
//		  If it is empty, automatically created in the program working directory
//		  If it is a directory, it will be created in the specified directory
func NewBasicNode(
	ctx context.Context,
	port int,
	workspace string,
	privatekeypath string,
	bootstrap []string,
	cmgr connmgr.ConnManager,
	protocolPrefix string,
	publicip string,
	enableBswap bool,
) (P2P, error) {
	if !FreeLocalPort(uint32(port)) {
		return nil, errors.New("port is in use")
	}

	var boots = make([]string, 0)
	for _, b := range bootstrap {
		bootnodes, err := ParseMultiaddrs(b)
		if err != nil {
			out.Err(err.Error())
			continue
		}
		boots = append(boots, bootnodes...)
	}

	if err := verifyWorkspace(workspace); err != nil {
		return nil, err
	}

	prvKey, err := identification(workspace, privatekeypath)
	if err != nil {
		return nil, err
	}

	var opts []libp2p.Option
	var multiaddrs []ma.Multiaddr
	if publicip == "" {
		externalIp, err := GetExternalIp()
		if err != nil {
			return nil, err
		}
		extMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", externalIp, port))
		multiaddrs = append(multiaddrs, extMultiAddr)
	} else {
		if !IsIPv4(publicip) {
			return nil, errors.New("invalid ipv4")
		}
		extMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", publicip, port))
		multiaddrs = append(multiaddrs, extMultiAddr)
	}

	addressFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		addrs = append(addrs, multiaddrs...)
		return addrs
	}

	rm, err := buildResourceManager()
	if err != nil {
		return nil, err
	}

	opts = append(opts,
		libp2p.Identity(prvKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.ConnectionManager(cmgr),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.DefaultPeerstore,
		libp2p.DefaultEnableRelay,
		libp2p.ProtocolVersion(protocolPrefix+p2pProtocolVer),
		libp2p.AddrsFactory(addressFactory),
		libp2p.DisableMetrics(),
		libp2p.ResourceManager(rm),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
		libp2p.EnableRelayService(),
		libp2p.EnableHolePunching(),
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

	dataDir, err := mkdir(workspace)
	if err != nil {
		return nil, err
	}

	ctxcancel, cancel := context.WithCancel(ctx)
	ctxreg, events := routing.RegisterForQueryEvents(ctxcancel)

	n := &Node{
		ctxRoot:               ctx,
		ctxCancelFromRoot:     ctxcancel,
		ctxQueryFromCtxCancel: ctxreg,
		ctxCancelFuncFromRoot: cancel,
		discoveredPeerCh:      events,
		host:                  bhost,
		workspace:             workspace,
		privatekeyPath:        privatekeypath,
		dir:                   dataDir,
		peerPublickey:         publickey,
		protocolVersion:       protocolPrefix + p2pProtocolVer,
		dhtProtocolVersion:    protocolPrefix + dhtProtocolVer,
		rendezvousVersion:     protocolPrefix + rendezvous,
		protocolPrefix:        protocolPrefix,
		bootstrap:             boots,
		enableBswap:           enableBswap,
		enableRecv:            true,
		protocols:             NewProtocol(),
	}

	err = n.initDHT()
	if err != nil {
		return nil, fmt.Errorf("[initDHT] %v", err)
	}

	if n.enableBswap {
		network := bsnet.NewFromIpfsHost(n.host, n.RoutingDiscovery)
		fsdatastore, err := NewDatastore(filepath.Join(n.workspace, FileBlockDir))
		if err != nil {
			return nil, err
		}
		n.bstore = blockstore.NewBlockstore(ds_sync.MutexWrap(fsdatastore))
		n.bswap = bitswap.New(
			n.ctxQueryFromCtxCancel,
			network,
			n.bstore,
		)
	}

	n.initProtocol(protocolPrefix)

	return n, nil
}

func (n *Node) GetBlockstore() blockstore.Blockstore {
	return n.bstore
}

func (n *Node) GetBitSwap() *bitswap.Bitswap {
	return n.bswap
}

// SaveAndNotifyDataBlock
func (n *Node) SaveAndNotifyDataBlock(buf []byte) (cid.Cid, error) {
	if !n.enableBswap {
		return cid.Cid{}, errors.New("The bitswap function is not enabled")
	}
	blockData := blocks.NewBlock(buf)
	err := n.bstore.Put(n.ctxQueryFromCtxCancel, blockData)
	if err != nil {
		return blockData.Cid(), err
	}
	err = n.bswap.NotifyNewBlocks(n.ctxQueryFromCtxCancel, blockData)
	return blockData.Cid(), err
}

// NotifyData notify data
func (n *Node) NotifyData(buf []byte) error {
	if !n.enableBswap {
		return errors.New("The bitswap function is not enabled")
	}
	blockData := blocks.NewBlock(buf)
	return n.bswap.NotifyNewBlocks(n.ctxQueryFromCtxCancel, blockData)
}

// GetDataFromBlock get data from block
func (n *Node) GetDataFromBlock(ctx context.Context, wantCid string) ([]byte, error) {
	if !n.enableBswap {
		return nil, errors.New("The bitswap function is not enabled")
	}
	wantcid, err := cid.Decode(wantCid)
	if err != nil {
		return nil, err
	}
	block, err := n.bswap.GetBlock(ctx, wantcid)
	if err != nil {
		return nil, err
	}
	return block.RawData(), err
}

// GetLocalDataFromBlock get local data from block
func (n *Node) GetLocalDataFromBlock(wantCid string) ([]byte, error) {
	if !n.enableBswap {
		return nil, errors.New("The bitswap function is not enabled")
	}
	wantcid, err := cid.Decode(wantCid)
	if err != nil {
		return nil, err
	}
	block, err := n.bstore.Get(n.ctxQueryFromCtxCancel, wantcid)
	if err != nil {
		return nil, err
	}
	return block.RawData(), err
}

// FidToCid
func (n *Node) FidToCid(fid string) (string, error) {
	mhash, err := mh.FromHexString("1220" + fid)
	if err != nil {
		return "", err
	}
	return cid.NewCidV0(mhash).String(), nil
}

// DHTFindPeer searches for a peer with given ID.
func (n *Node) DHTFindPeer(peerid string) (peer.AddrInfo, error) {
	id, err := peer.Decode(peerid)
	if err != nil {
		return peer.AddrInfo{}, err
	}
	return n.IpfsDHT.FindPeer(n.ctxQueryFromCtxCancel, id)
}

// RouteTableFindPeers
func (n *Node) RouteTableFindPeers(limit int) (<-chan peer.AddrInfo, error) {
	if limit <= 0 {
		return n.RoutingDiscovery.FindPeers(n.ctxQueryFromCtxCancel, n.rendezvousVersion)
	}
	return n.RoutingDiscovery.FindPeers(n.ctxQueryFromCtxCancel, n.rendezvousVersion, discovery.Limit(limit))
}

// PeerID returns for own peerid
func (n *Node) PeerID() peer.ID {
	return n.IpfsDHT.PeerID()
}

// GetDiscoveredPeers
func (n *Node) GetDiscoveredPeers() <-chan *routing.QueryEvent {
	return n.discoveredPeerCh
}

// ID returns the (local) peer.ID associated with this Host
func (n *Node) ID() peer.ID {
	return n.host.ID()
}

// Peerstore returns the Host's repository of Peer Addresses and Keys.
func (n *Node) Peerstore() peerstore.Peerstore {
	return n.host.Peerstore()
}

// Returns the listen addresses of the Host
func (n *Node) Addrs() []ma.Multiaddr {
	return n.host.Addrs()
}

// Networks returns the Network interface of the Host
func (n *Node) Network() network.Network {
	return n.host.Network()
}

// Mux returns the Mux multiplexing incoming streams to protocol handlers
func (n *Node) Mux() protocol.Switch {
	return n.host.Mux()
}

// Connect ensures there is a connection between this host and the peer with
// given peer.ID. Connect will absorb the addresses in pi into its internal
// peerstore. If there is not an active connection, Connect will issue a
// h.Network.Dial, and block until a connection is open, or an error is
// returned. // TODO: Relay + NAT.
func (n *Node) Connect(ctx context.Context, pi peer.AddrInfo) error {
	return n.host.Connect(ctx, pi)
}

// SetStreamHandler sets the protocol handler on the Host's Mux.
// This is equivalent to: host.Mux().SetHandler(proto, handler) (Threadsafe)
func (n *Node) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	n.host.SetStreamHandler(pid, handler)
}

// SetStreamHandlerMatch sets the protocol handler on the Host's Mux
// using a matching function for protocol selection.
func (n *Node) SetStreamHandlerMatch(pid protocol.ID, m func(protocol.ID) bool, handler network.StreamHandler) {
	n.host.SetStreamHandlerMatch(pid, m, handler)
}

// RemoveStreamHandler removes a handler on the mux that was set by
// SetStreamHandler
func (n *Node) RemoveStreamHandler(pid protocol.ID) {
	n.host.RemoveStreamHandler(pid)
}

// NewStream opens a new stream to given peer p, and writes a p2p/protocol
// header with given ProtocolID. If there is no connection to p, attempts
// to create one. If ProtocolID is "", writes no header.
// (Threadsafe)
func (n *Node) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, p, pids...)
}

// Close shuts down the host, its Network, and services.
func (n *Node) Close() error {
	err := n.host.Close()
	if err != nil {
		return err
	}
	n.ctxCancelFuncFromRoot()
	return nil
}

// ConnManager returns this hosts connection manager
func (n *Node) ConnManager() connmgr.ConnManager {
	return n.host.ConnManager()
}

// EventBus returns the hosts eventbus
func (n *Node) EventBus() event.Bus {
	return n.host.EventBus()
}

func (n *Node) GetProtocolPrefix() string {
	return n.protocolPrefix
}

func (n *Node) AddMultiaddrToPeerstore(multiaddr string, t time.Duration) (peer.ID, error) {
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

func (n *Node) GetPeerPublickey() []byte {
	return n.peerPublickey
}

func (n *Node) GetProtocolVersion() string {
	return n.protocolVersion
}

func (n *Node) GetDhtProtocolVersion() string {
	return n.dhtProtocolVersion
}

func (n *Node) GetRendezvousVersion() string {
	return n.rendezvousVersion
}

func (n *Node) GetBootstraps() []string {
	return n.bootstrap
}

func (n *Node) SetBootstraps(bootstrap []string) {
	n.bootstrap = bootstrap
}

func (n *Node) GetDht() *dht.IpfsDHT {
	return n.IpfsDHT
}

func (n *Node) GetRoutingTable() *drouting.RoutingDiscovery {
	return n.RoutingDiscovery
}

func (n *Node) GetCtxRoot() context.Context {
	return n.ctxRoot
}

func (n *Node) GetCtxCancelFromRoot() context.Context {
	return n.ctxCancelFromRoot
}

func (n *Node) GetCtxQueryFromCtxCancel() context.Context {
	return n.ctxQueryFromCtxCancel
}

func (n *Node) GetDirs() DataDirs {
	return n.dir
}

func (n *Node) EnableRecv() {
	n.enableRecv = true
}

func (n *Node) DisableRecv() {
	n.enableRecv = false
}

func (n *Node) SetIdleFileTee(peerid string) {
	n.idleTee.Store(peerid)
}

func (n *Node) GetIdleFileTee() string {
	value, ok := n.idleTee.Load().(string)
	if !ok {
		return ""
	}
	return value
}

func (n *Node) SetServiceFileTee(peerid string) {
	n.serviceTee.Store(peerid)
}

func (n *Node) GetServiceFileTee() string {
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

func (n *Node) initProtocol(protocolPrefix string) {
	n.SetProtocolPrefix(protocolPrefix)
	n.WriteFileProtocol = n.NewWriteFileProtocol()
	n.ReadFileProtocol = n.NewReadFileProtocol()
	n.ReadDataProtocol = n.NewReadDataProtocol()
	n.ReadDataStatProtocol = n.NewReadDataStatProtocol()
}

func (n *Node) initDHT() error {
	var options []dht.Option
	options = append(options,
		dht.V1ProtocolOverride(protocol.ID(n.dhtProtocolVersion)),
		dht.Resiliency(10),
		dht.DisableAutoRefresh(),
		dht.Mode(dht.ModeAutoServer),
	)

	bootstrap := n.bootstrap
	// var bootaddrs []peer.AddrInfo
	// for _, v := range bootstrap {
	// 	muladdr, err := ma.NewMultiaddr(v)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	addrinfo, err := peer.AddrInfoFromP2pAddr(muladdr)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	bootaddrs = append(bootaddrs, *addrinfo)
	// }

	// if len(bootaddrs) > 0 {
	// 	options = append(options, dht.BootstrapPeers(bootaddrs...))
	// }

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(n.ctxQueryFromCtxCancel, n.host, options...)
	if err != nil {
		return err
	}

	if err = kademliaDHT.Bootstrap(n.ctxQueryFromCtxCancel); err != nil {
		return err
	}

	for _, peerAddr := range bootstrap {
		bootstrapAddr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			continue
		}
		peerinfo, err := peer.AddrInfoFromP2pAddr(bootstrapAddr)
		if err != nil {
			continue
		}
		n.host.Connect(n.ctxQueryFromCtxCancel, *peerinfo)
		// if err != nil {
		// 	out.Err(fmt.Sprintf("Connection to boot node failed: %s", peerinfo.ID.Pretty()))
		// } else {
		// 	out.Ok(fmt.Sprintf("Connection to boot node successful: %s", peerinfo.ID.Pretty()))
		// }
		kademliaDHT.RoutingTable().PeerAdded(peerinfo.ID)
		n.AddMultiaddrToPeerstore(bootstrapAddr.String(), peerstore.PermanentAddrTTL)
	}
	n.RoutingDiscovery = drouting.NewRoutingDiscovery(kademliaDHT)
	n.IpfsDHT = kademliaDHT
	dutil.Advertise(n.ctxQueryFromCtxCancel, n.RoutingDiscovery, n.rendezvousVersion)
	return nil
}

func getNumFDs() int {
	var l unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &l); err != nil {
		out.Warn(fmt.Sprintf("failed to get fd limit: %v", err))
		return DefaultFDCount
	}
	return int(l.Cur)
}

func buildResourceManager() (network.ResourceManager, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits

	// Add limits around included libp2p protocols
	libp2p.SetDefaultServiceLimits(&scalingLimits)

	// Turn the scaling limits into a concrete set of limits using `.AutoScale`. This
	// scales the limits proportional to your system memory.
	//scaledDefaultLimits := scalingLimits.AutoScale()

	scaledDefaultLimits := scalingLimits.Scale(int64(memory.TotalMemory()/10*8), int(getNumFDs()/10*8))

	// Tweak certain settings
	cfg := rcmgr.PartialLimitConfig{
		System: rcmgr.ResourceLimits{
			// Allow unlimited outbound streams
			StreamsOutbound: rcmgr.Unlimited,
		},
		// Everything else is default. The exact values will come from `scaledDefaultLimits` above.
	}

	// Create our limits by using our cfg and replacing the default values with values from `scaledDefaultLimits`
	limits := cfg.Build(scaledDefaultLimits)

	// The resource manager expects a limiter, se we create one from our limits.
	limiter := rcmgr.NewFixedLimiter(limits)

	// Metrics are enabled by default. If you want to disable metrics, use the
	// WithMetricsDisabled option
	// Initialize the resource manager
	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, errors.Wrapf(err, "[NewResourceManager]")
	}
	return rm, nil
}
