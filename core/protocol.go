/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Protocol interface {
	WriteDataAction(ctx context.Context, id peer.ID, file, fid, fragment string) error
	ReadDataAction(ctx context.Context, id peer.ID, name, savepath string) (int64, error)
	ReadDataStatAction(ctx context.Context, id peer.ID, name string) (uint64, error)
	OnlineAction(id peer.ID) error
	GetRecord() <-chan string
}

type protocols struct {
	ProtocolPrefix string
	*WriteDataProtocol
	*ReadDataProtocol
	*ReadDataStatProtocol
	*OnlineProtocol
}

func NewProtocol() *protocols {
	return &protocols{}
}

func (p *protocols) SetProtocolPrefix(protocolPrefix string) {
	p.ProtocolPrefix = protocolPrefix
}
