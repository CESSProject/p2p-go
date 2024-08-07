/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package p2pgo

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

// set protocol prefix
func ProtocolPrefix(protocolPrefix string) Option {
	return func(cfg *Config) error {
		cfg.ProtocolPrefix = protocolPrefix
		return nil
	}
}

// set timeout for dialing
func DialTimeout(timeout int) Option {
	return func(cfg *Config) error {
		cfg.DialTimeout = timeout
		return nil
	}
}

// set the boot node to record the cache channel length of other nodes
func RecordCacheLen(length int) Option {
	return func(cfg *Config) error {
		cfg.RecordCacheLen = length
		return nil
	}
}
