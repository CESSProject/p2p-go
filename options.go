/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"time"

	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
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

// MaxConnection configuration max connection
func MaxConnection(low, high int) Option {
	return func(cfg *Config) error {
		if low < 0 {
			low = 1
		}
		if high < 0 {
			high = low
		}
		mgr, err := connmgr.NewConnManager(low, high, connmgr.WithGracePeriod(time.Hour), connmgr.WithSilencePeriod(time.Minute))
		if err != nil {
			return nil
		}
		cfg.ConnManager = mgr
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

// PrivatekeyFile configuration privatekey file
func PrivatekeyFile(privatekey string) Option {
	return func(cfg *Config) error {
		cfg.PrivatekeyPath = privatekey
		return nil
	}
}

// PrivatekeyFile configuration privatekey file
func ProtocolPrefix(protocolPrefix string) Option {
	return func(cfg *Config) error {
		cfg.ProtocolPrefix = protocolPrefix
		return nil
	}
}

// PrivatekeyFile configuration privatekey file
func PublicIpv4(ip string) Option {
	return func(cfg *Config) error {
		cfg.PublicIpv4 = ip
		return nil
	}
}

// BucketSize configuration bucket size
func BucketSize(size int) Option {
	return func(cfg *Config) error {
		cfg.BucketSize = size
		return nil
	}
}

// Version configuration version
func Version(ver string) Option {
	return func(cfg *Config) error {
		cfg.Version = ver
		return nil
	}
}
