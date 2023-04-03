/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/connmgr"
	ma "github.com/multiformats/go-multiaddr"
)

// ListenAddrStrings configures libp2p to listen on the given (unparsed)
// addresses.
func ListenAddrStrings(s ...string) Option {
	return func(cfg *Config) error {
		for _, addrstr := range s {
			a, err := ma.NewMultiaddr(addrstr)
			if err != nil {
				return err
			}
			cfg.ListenAddrs = append(cfg.ListenAddrs, a)
		}
		return nil
	}
}

// ListenAddrs configures libp2p to listen on the given addresses.
func ListenAddrs(addrs ...ma.Multiaddr) Option {
	return func(cfg *Config) error {
		cfg.ListenAddrs = append(cfg.ListenAddrs, addrs...)
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
//
// The current "standard" connection manager lives in github.com/libp2p/go-libp2p-connmgr. See
// https://pkg.go.dev/github.com/libp2p/go-libp2p-connmgr?utm_source=godoc#NewConnManager.
func ConnectionManager(connman connmgr.ConnManager) Option {
	return func(cfg *Config) error {
		if cfg.ConnManager != nil {
			return fmt.Errorf("cannot specify multiple connection managers")
		}
		cfg.ConnManager = connman
		return nil
	}
}
