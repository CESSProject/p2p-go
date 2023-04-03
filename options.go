package p2pgo

import (
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
