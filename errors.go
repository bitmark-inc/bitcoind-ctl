// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import "fmt"

var (
	ErrReadBitcoindStatus   = fmt.Errorf("fail to read bitcoind status")
	ErrBitcoindStopped      = fmt.Errorf("bitcoind is stopped")
	ErrBitcoindProcNotReady = fmt.Errorf("bitcoind process is not ready")
	ErrBitcoindRPCNotReady  = fmt.Errorf("bitcoind RPC is not ready")
	ErrBitcoindUnknown      = fmt.Errorf("unknown error from bitcoind RPC")
)
