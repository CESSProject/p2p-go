package p2pgo

import (
	"os"

	"github.com/multiformats/go-multiaddr"
)

// DefaultListenAddrs configures libp2p to use default listen address.
var DefaultListenAddrs = func(cfg *Config) error {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
	listenAddrs := make([]multiaddr.Multiaddr, 0, len(addrs))
	for _, s := range addrs {
		addr, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			return err
		}
		listenAddrs = append(listenAddrs, addr)
	}
	return cfg.Apply(ListenAddrs(listenAddrs...))
}

// DefaultListenAddrs configures libp2p to use default listen address.
var DefaultWorkspace = func(cfg *Config) error {
	workspace, err := os.Getwd()
	if err != nil {
		return err
	}
	return cfg.Apply(Workspace(workspace))
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
