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

// DefaultListenPort configures libp2p to use default port.
var DefaultListenPort = func(cfg *Config) error {
	port := 4001
	return cfg.Apply(ListenPort(port))
}

// DefaultConnectionManager creates a default connection manager.
var DefaultConnectionManager = func(cfg *Config) error {
	mgr, err := connmgr.NewConnManager(10, 100, connmgr.WithGracePeriod(time.Minute), connmgr.WithSilencePeriod(time.Minute))
	if err != nil {
		return err
	}
	cfg.ConnManager = mgr
	return nil
}

// DefaultWorkSpace configures libp2p to use default work space.
var DefaultWorkSpace = func(cfg *Config) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	return cfg.Apply(Workspace(dir))
}

// DefaultWorkSpace configures libp2p to use default work space.
var DefaultBucketSize = func(cfg *Config) error {
	return cfg.Apply(BucketSize(100))
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
		fallback: func(cfg *Config) bool { return cfg.Workspace == "" },
		opt:      DefaultWorkSpace,
	},
	{
		fallback: func(cfg *Config) bool { return cfg.BucketSize == 0 },
		opt:      DefaultBucketSize,
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
