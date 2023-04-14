/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"fmt"

	"github.com/CESSProject/p2p-go/core"
	"github.com/libp2p/go-libp2p/core/connmgr"
	ma "github.com/multiformats/go-multiaddr"
)

// ListenAddrStrings configures libp2p to listen on the given (unparsed)
// addresses.
func ListenAddrStrings(ip string, port int) Option {
	return func(cfg *Config) error {
		var allip = core.AllIpAddress
		if ip == core.LocalAddress {
			allip = ip
		}
		a, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", allip, port))
		if err != nil {
			return err
		}
		cfg.ListenAddrs = a
		return nil
	}
}

// ListenAddrs configures libp2p to listen on the given addresses.
func ListenAddrs(addrs ma.Multiaddr) Option {
	return func(cfg *Config) error {
		cfg.ListenAddrs = addrs
		return nil
	}
}

// ListenAddrs configures libp2p to listen on the given addresses.
func Workspace(workspace string) Option {
	return func(cfg *Config) error {
		cfg.Workspace = workspace
		return nil
	}
}

// ConnectionManager configures libp2p to use the given connection manager.
func ConnectionManager(connman connmgr.ConnManager) Option {
	return func(cfg *Config) error {
		if cfg.ConnManager != nil {
			return fmt.Errorf("cannot specify multiple connection managers")
		}
		cfg.ConnManager = connman
		return nil
	}
}

// BootPeers configures bootstrap nodes
func BootPeers(bootpeers []string) Option {
	return func(cfg *Config) error {
		cfg.BootPeers = bootpeers
		return nil
	}
}
