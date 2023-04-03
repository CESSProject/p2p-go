/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/core"
)

// Config describes a set of settings for a p2p node.
type Config = config.Config

// Option is a p2p config option that can be given to the p2p constructor
type Option = config.Option

// New constructs a new libp2p node with the given options, falling back on
// reasonable defaults. The defaults are:
//
// - If no transport and listen addresses are provided, the node listens to
// the multiaddresses "/ip4/0.0.0.0/tcp/0" and "/ip6/::/tcp/0";
//
// To stop/shutdown the returned p2p node, the user needs to cancel the passed
// context and call `Close` on the returned Host.
func New(privatekeyPath string, opts ...Option) (core.P2P, error) {
	return NewWithoutDefaults(privatekeyPath, append(opts, FallbackDefaults)...)
}

// NewWithoutDefaults constructs a new libp2p node with the given options but
// *without* falling back on reasonable defaults.
//
// Warning: This function should not be considered a stable interface. We may
// choose to add required services at any time and, by using this function, you
// opt-out of any defaults we may provide.
func NewWithoutDefaults(privatekeyPath string, opts ...Option) (core.P2P, error) {
	var cfg Config
	if err := cfg.Apply(opts...); err != nil {
		return nil, err
	}
	return cfg.NewNode(privatekeyPath)
}
