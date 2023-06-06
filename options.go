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

// ListenPort configuration listening port
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

// ConnectionManager configuration connection manager
func ConnectionManager(connman connmgr.ConnManager) Option {
	return func(cfg *Config) error {
		if cfg.ConnManager != nil {
			return fmt.Errorf("cannot specify multiple connection managers")
		}
		cfg.ConnManager = connman
		return nil
	}
}

// BootPeers configuration bootstrap nodes
func BootPeers(bootpeers []string) Option {
	return func(cfg *Config) error {
		cfg.BootPeers = bootpeers
		return nil
	}
}

// ProtocolVersion configuration protocol version
func ProtocolVersion(protocolVersion string) Option {
	return func(cfg *Config) error {
		cfg.ProtocolVersion = protocolVersion
		return nil
	}
}

// DhtProtocolVersion configuration DHT protocol version
func DhtProtocolVersion(dhtProtocolVersion string) Option {
	return func(cfg *Config) error {
		cfg.DhtProtocolVersion = dhtProtocolVersion
		return nil
	}
}

// PrivatekeyFile configuration privatekey file
func PrivatekeyFile(privatekey string) Option {
	return func(cfg *Config) error {
		cfg.PrivatekeyPath = privatekey
		return nil
	}
}
