/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package config

// Config describes a set of settings for a libp2p node
type Config struct {
	ListenPort     int
	BootPeers      []string
	Workspace      string
	PrivatekeyPath string
	ProtocolPrefix string
	PublicIpv4     string
	RecordCacheLen int
	DialTimeout    int
}

// Option is a libp2p config option that can be given to the libp2p constructor
// (`libp2p.New`).
type Option func(cfg *Config) error

const (
	DevnetProtocolPrefix  = "/devnet"
	TestnetProtocolPrefix = "/testnet"
	MainnetProtocolPrefix = "/mainnet"
)

const (
	ServerMode_Server = "server"
	ServerMode_Client = "client"
	ServerMode_Auto   = "auto"
)

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
