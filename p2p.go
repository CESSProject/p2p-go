/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

import (
	"context"

	"github.com/CESSProject/p2p-go/config"
	"github.com/CESSProject/p2p-go/core"
)

// Config describes a set of settings for a p2p peer.
type Config = config.Config

// Option is a p2p config option that can be given to the p2p constructor.
type Option = config.Option

// New constructs a new libp2p node with the given options, falling back on
// reasonable defaults. The defaults are:
//
// - If no listening port is provided, the host listens on the default
// port: 4001
//
// - If no bootstrap nodes is provided, the host will run as a server
//
// - If no protocol version is provided, The host uses the default
// protocol version: /kldr/1.0
//
// - If no DHT protocol version is provided, The host uses the default
// DHT protocol version: /kldr/kad/1.0
//
// To stop/shutdown the returned p2p node, the user needs to cancel the passed
// context and call `Close` on the returned Host.
func New(ctx context.Context, opts ...Option) (core.P2P, error) {
	return NewWithoutDefaults(ctx, append(opts, FallbackDefaults)...)
}

// NewWithoutDefaults constructs a new libp2p node with the given options but
// *without* falling back on reasonable defaults.
//
// Warning: This function should not be considered a stable interface. We may
// choose to add required services at any time and, by using this function, you
// opt-out of any defaults we may provide.
func NewWithoutDefaults(ctx context.Context, opts ...Option) (core.P2P, error) {
	var cfg Config
	if err := cfg.Apply(opts...); err != nil {
		return nil, err
	}
	return cfg.NewNode(ctx)
}
