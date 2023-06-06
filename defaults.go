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

// DefaultListenAddrs configures libp2p to use default listen address.
var DefaultListenPort = func(cfg *Config) error {
	port := 4001
	return cfg.Apply(ListenPort(port))
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
	defaultBotPeers := []string{"_dnsaddr.bootstrap-kldr.cess.cloud"}
	return cfg.Apply(BootPeers(defaultBotPeers))
}

// DefaultConnectionManager creates a default connection manager
var DefaultProtocolVersion = func(cfg *Config) error {
	defaultProtocolVersion := "/kldr/1.0"
	return cfg.Apply(ProtocolVersion(defaultProtocolVersion))
}

// DefaultConnectionManager creates a default connection manager
var DefaultDhtProtocolVersion = func(cfg *Config) error {
	defaultDhtProtocolVersion := "/kldr/kad/1.0"
	return cfg.Apply(DhtProtocolVersion(defaultDhtProtocolVersion))
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
		fallback: func(cfg *Config) bool { return cfg.ListenPort == 0 },
		opt:      DefaultListenPort,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.ConnManager == nil },
		opt:      DefaultConnectionManager,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.BootPeers == nil },
		opt:      DefaultBootPeers,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.ProtocolVersion == "" },
		opt:      DefaultProtocolVersion,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.DhtProtocolVersion == "" },
		opt:      DefaultDhtProtocolVersion,
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
