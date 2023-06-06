/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/connmgr"
)

// ListenAddrStrings configures libp2p to listen on the given (unparsed)
// addresses.
func ListenPort(port int) Option {
	return func(cfg *Config) error {
		cfg.ListenPort = port
		return nil
	}
}

// Workspace configuration working directory
func Workspace(workspace string) Option {
	return func(cfg *Config) error {
		cfg.Workspace = workspace
		return nil
	}
}

// ConnectionManager configures libp2p to use the given connection manager
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

// BootPeers configures bootstrap nodes
func ProtocolVersion(protocolVersion string) Option {
	return func(cfg *Config) error {
		cfg.ProtocolVersion = protocolVersion
		return nil
	}
}

// BootPeers configures bootstrap nodes
func DhtProtocolVersion(dhtProtocolVersion string) Option {
	return func(cfg *Config) error {
		cfg.DhtProtocolVersion = dhtProtocolVersion
		return nil
	}
}

// BootPeers configures bootstrap nodes
func PrivatekeyFile(privatekey string) Option {
	return func(cfg *Config) error {
		cfg.PrivatekeyPath = privatekey
		return nil
	}
}
