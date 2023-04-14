/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"os"
	"time"

	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

// DefaultListenAddrs configures libp2p to use default listen address.
var DefaultListenAddrs = func(cfg *Config) error {
	ip := "0.0.0.0"
	port := 80
	return cfg.Apply(ListenAddrStrings(ip, port))
}

// DefaultListenAddrs configures libp2p to use default listen address.
var DefaultWorkspace = func(cfg *Config) error {
	workspace, err := os.Getwd()
	if err != nil {
		return err
	}
	return cfg.Apply(Workspace(workspace))
}

// DefaultConnectionManager creates a default connection manager
var DefaultConnectionManager = func(cfg *Config) error {
	mgr, err := connmgr.NewConnManager(100, 1000, connmgr.WithGracePeriod(time.Minute))
	if err != nil {
		return err
	}

	return cfg.Apply(ConnectionManager(mgr))
}

// DefaultConnectionManager creates a default connection manager
var DefaultBootPeers = func(cfg *Config) error {
	defaultBotPeers := []string{
		"",
	}

	return cfg.Apply(BootPeers(defaultBotPeers))
}

// Complete list of default options and when to fallback on them.
//
// Please *DON'T* specify default options any other way. Putting this all here
// makes tracking defaults *much* easier.
var defaults = []struct {
	fallback func(cfg *Config) bool
	opt      Option
}{
	{
		fallback: func(cfg *Config) bool { return cfg.ListenAddrs == nil },
		opt:      DefaultListenAddrs,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.Workspace == "" },
		opt:      DefaultWorkspace,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.ConnManager == nil },
		opt:      DefaultConnectionManager,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.BootPeers == nil },
		opt:      DefaultBootPeers,
	},
}

// FallbackDefaults applies default options to the libp2p node if and only if no
// other relevant options have been applied. will be appended to the options
// passed into New.
var FallbackDefaults Option = func(cfg *Config) error {
	for _, def := range defaults {
		if !def.fallback(cfg) {
			continue
		}
		if err := cfg.Apply(def.opt); err != nil {
			return err
		}
	}
	return nil
}
