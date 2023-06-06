/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"context"

	"github.com/CESSProject/p2p-go/core"
	"github.com/libp2p/go-libp2p/core/connmgr"
)

// Config describes a set of settings for a libp2p node
//
// This is *not* a stable interface. Use the options defined in the root
// package.
type Config struct {
	ListenPort         int
	ConnManager        connmgr.ConnManager
	BootPeers          []string
	Workspace          string
	ProtocolVersion    string
	DhtProtocolVersion string
	PrivatekeyPath     string
}

// Option is a libp2p config option that can be given to the libp2p constructor
// (`libp2p.New`).
type Option func(cfg *Config) error

// NewNode constructs a new libp2p Host from the Config.
//
// This function consumes the config. Do not reuse it (really!).
func (cfg *Config) NewNode(ctx context.Context) (core.P2P, error) {
	return core.NewBasicNode(ctx, cfg.ListenPort, cfg.Workspace, cfg.ProtocolVersion, cfg.DhtProtocolVersion, cfg.PrivatekeyPath, cfg.BootPeers, cfg.ConnManager)
}

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}
