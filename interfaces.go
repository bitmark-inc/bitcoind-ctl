// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
)

// BitcoindStatus is the response of bitcoind status
type BitcoindStatus struct {
	BestBlock    int     `json:"best_block"`
	SyncProgress float64 `json:"sync_progress"`
}

// BitcoindCtl is a controller to control and report the lifecycle of a bitcoind
type BitcoindCtl interface {
	GetBitcoindStatus(ctx context.Context) (BitcoindStatus, error)
	StartBitcoind(ctx context.Context) error
	StopBitcoind(ctx context.Context, force bool) error
}
