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
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CESSProject/p2p-go/pb"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"
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

	// GetRootCtx returns the root context of the host
	GetRootCtx() context.Context

	// GetDiscoverSt returns whether the discovery service is running
	GetDiscoverSt() bool

	// StartDiscover starts the node discovery service, If you have
	// already started the service, calling it again will have no effect.
	StartDiscover()

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

	// GetProtocolPrefix returns protocols prefix
	GetProtocolPrefix() string

	// GetDirs returns the data directory structure of the host
	GetDirs() DataDirs

	// GetBootstraps returns a list of host bootstraps
	GetBootstraps() []string

	// SetBootstraps updates the host's bootstrap list
	SetBootstraps(bootstrap []string)

	// DiscoveredPeer returns the node channel discovered by the host
	DiscoveredPeer() <-chan peer.AddrInfo

	// GetIdleDataCh returns the idle data channel received by the host
	GetIdleDataCh() <-chan string

	// GetIdleTagCh returns the tag channel of the idle data received by the host
	GetIdleTagCh() <-chan string

	// GetServiceTagCh returns the tag channel of the service data received by the host
	GetServiceTagCh() <-chan string
}

// Node type - Implementation of a P2P Host
type Node struct {
	ctx                context.Context
	ctxCancel          context.Context
	ctxReg             context.Context
	cancelFunc         context.CancelFunc
	discoverEvent      <-chan *routing.QueryEvent
	host               host.Host
	dir                DataDirs
	peerPublickey      []byte
	workspace          string
	privatekeyPath     string
	idleTee            atomic.Value
	serviceTee         atomic.Value
	idleDataCh         chan string
	idleTagDataCh      chan string
	serviceTagDataCh   chan string
	protocolVersion    string
	dhtProtocolVersion string
	protocolPrefix     string
	discoverStat       atomic.Uint32
	bootstrap          []string
	discoveredPeerCh   chan peer.AddrInfo
	*protocols
}

var _ P2P = (*Node)(nil)

// NewBasicNode constructs a new *Node
//
//	  multiaddr: listen addresses of p2p host
//	  workspace: service working directory
//	  privatekeypath: private key file
//		  If it's empty, automatically created in the program working directory
//		  If it's a directory, it will be created in the specified directory
func NewBasicNode(
	ctx context.Context,
	port int,
	workspace string,
	privatekeypath string,
	bootstrap []string,
	cmgr connmgr.ConnManager,
	protocolPrefix string,
) (P2P, error) {

	if protocolPrefix == "" {
		protocolPrefix = "/kldr"
	}

	var boots = make([]string, 0)
	for _, b := range bootstrap {
		bootnodes, err := ParseMultiaddrs(b)
		if err != nil {
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

	externalIp, err := GetExternalIp()
	if err == nil {
		extMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", externalIp, port))
		multiaddrs = append(multiaddrs, extMultiAddr)
	}

	addressFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		addrs = append(addrs, multiaddrs...)
		return addrs
	}

	opts = append(opts,
		libp2p.Identity(prvKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.ConnectionManager(cmgr),
		libp2p.DefaultTransports,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.ProtocolVersion(protocolPrefix+p2pProtocolVer),
		libp2p.DefaultMuxers,
		libp2p.AddrsFactory(addressFactory),
	)

	bhost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	if !bhost.ID().MatchesPrivateKey(prvKey) {
		return nil, errors.New("")
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
		ctx:                ctx,
		ctxCancel:          ctxcancel,
		ctxReg:             ctxreg,
		cancelFunc:         cancel,
		discoverEvent:      events,
		host:               bhost,
		workspace:          workspace,
		privatekeyPath:     privatekeypath,
		dir:                dataDir,
		peerPublickey:      publickey,
		idleDataCh:         make(chan string, 1),
		idleTagDataCh:      make(chan string, 1),
		serviceTagDataCh:   make(chan string, 1),
		protocolVersion:    protocolPrefix + p2pProtocolVer,
		dhtProtocolVersion: protocolPrefix + dhtProtocolVer,
		protocolPrefix:     protocolPrefix,
		discoverStat:       atomic.Uint32{},
		bootstrap:          boots,
		discoveredPeerCh:   make(chan peer.AddrInfo, 600),
		protocols:          NewProtocol(),
	}

	n.initProtocol(protocolPrefix)

	go n.discoverPeers(n.ctxReg, n.host, n.protocolPrefix, n.dhtProtocolVersion, n.bootstrap)

	return n, nil
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
	n.cancelFunc()
	err := n.host.Close()
	if err != nil {
		return err
	}
	close(n.idleDataCh)
	close(n.idleTagDataCh)
	close(n.serviceTagDataCh)
	close(n.discoveredPeerCh)
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

func (n *Node) GetDiscoverSt() bool {
	return n.discoverStat.Load() > 0
}

func (n *Node) GetProtocolPrefix() string {
	return n.protocolPrefix
}

func (n *Node) StartDiscover() {
	go n.discoverPeers(n.ctxReg, n.host, n.protocolPrefix, n.dhtProtocolVersion, n.bootstrap)
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

func (n *Node) DiscoveredPeer() <-chan peer.AddrInfo {
	return n.discoveredPeerCh
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

func (n *Node) GetBootstraps() []string {
	return n.bootstrap
}

func (n *Node) SetBootstraps(bootstrap []string) {
	n.bootstrap = bootstrap
}

func (n *Node) GetRootCtx() context.Context {
	return n.ctx
}

func (n *Node) GetDirs() DataDirs {
	return n.dir
}

func (n *Node) putIdleDataCh(path string) {
	if len(n.idleDataCh) > 0 {
		_ = <-n.idleDataCh
	}
	n.idleDataCh <- path
}

func (n *Node) GetIdleDataCh() <-chan string {
	return n.idleDataCh
}

func (n *Node) putIdleTagCh(path string) {
	if len(n.idleTagDataCh) > 0 {
		_ = <-n.idleTagDataCh
	}
	n.idleTagDataCh <- path
}

func (n *Node) GetIdleTagCh() <-chan string {
	return n.idleTagDataCh
}

func (n *Node) putServiceTagCh(path string) {
	if len(n.serviceTagDataCh) > 0 {
		_ = <-n.serviceTagDataCh
	}
	n.serviceTagDataCh <- path
}

func (n *Node) GetServiceTagCh() <-chan string {
	return n.serviceTagDataCh
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
		FileDir:       filepath.Join(workspace, FileDataDirectionry),
		TmpDir:        filepath.Join(workspace, TmpDataDirectionry),
		IdleDataDir:   filepath.Join(workspace, IdleDataDirectionry),
		IdleTagDir:    filepath.Join(workspace, IdleTagDirectionry),
		ServiceTagDir: filepath.Join(workspace, ServiceTagDirectionry),
		ProofDir:      filepath.Join(workspace, ProofDirectionry),
		IproofFile:    filepath.Join(workspace, ProofDirectionry, IdleProofFile),
		SproofFile:    filepath.Join(workspace, ProofDirectionry, ServiceProofFile),
	}
	if err := os.MkdirAll(dataDir.FileDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.TmpDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.IdleDataDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.IdleTagDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.ServiceTagDir, DirMode); err != nil {
		return dataDir, err
	}
	if err := os.MkdirAll(dataDir.ProofDir, DirMode); err != nil {
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

func (n *Node) discoverPeers(ctx context.Context, h host.Host, protocolPrefix, dhtProtocolVersion string, bootstrap []string) {
	if n.discoverStat.Load() > 0 {
		return
	}
	ctxCancel, cancel := context.WithCancel(ctx)
	defer func() {
		recover()
		cancel()
		n.discoverStat.Store(0)
	}()

	n.discoverStat.Add(1)
	time.Sleep(time.Second)
	if n.discoverStat.Load() != 1 {
		return
	}

	kademliaDHT, err := initDHT(ctx, h, dhtProtocolVersion, bootstrap)
	if err != nil {
		return
	}

	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, protocolPrefix+Rendezvous)

	go findPeers(ctxCancel, routingDiscovery, protocolPrefix+Rendezvous)

	for {
		select {
		case peer := <-n.discoverEvent:
			for _, v := range peer.Responses {
				n.discoveredPeerCh <- *v
			}
		}
	}
}

func findPeers(ctx context.Context, routingDiscovery *drouting.RoutingDiscovery, rendezvous string) {
	var ok bool

	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			peerChan, err := routingDiscovery.FindPeers(ctx, rendezvous, discovery.Limit(300))
			if err == nil {
				for {
					_, ok = <-peerChan
					if !ok {
						break
					}
				}
			}
		}
	}
}

func (n *Node) initProtocol(protocolPrefix string) {
	n.SetProtocolPrefix(protocolPrefix)
	n.WriteFileProtocol = n.NewWriteFileProtocol(protocolPrefix)
	n.ReadFileProtocol = n.NewReadFileProtocol(protocolPrefix)
	n.CustomDataTagProtocol = n.NewCustomDataTagProtocol()
	n.IdleDataTagProtocol = n.NewIdleDataTagProtocol()
	n.FileProtocol = n.NewFileProtocol(protocolPrefix)
	n.AggrProofProtocol = n.NewAggrProofProtocol()
	n.PushTagProtocol = n.NewPushTagProtocol(protocolPrefix)
}

func initDHT(ctx context.Context, h host.Host, dhtProtocolVersion string, bootstrap []string) (*dht.IpfsDHT, error) {
	var options []dht.Option
	options = append(options, dht.V1ProtocolOverride(protocol.ID(dhtProtocolVersion)))

	if len(bootstrap) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h, options...)
	if err != nil {
		return nil, err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrap {
		bootstrapAddr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			continue
		}

		peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrapAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := h.Connect(ctx, *peerinfo)
			if err != nil {
				log.Println("Failed to connect to the bootstrap node: ", peerinfo.ID.Pretty())
			} else {
				log.Println("Connected to the bootstrap node: ", peerinfo.ID.Pretty())
			}
		}()
	}
	wg.Wait()

	return kademliaDHT, nil
}
