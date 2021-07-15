// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// BitcoindRPCError is the error object of bitcoind RPC response
type BitcoindRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BitcoindResponse is the general response of bitcoind RPC
type BitcoindResponse struct {
	ID     string            `json:"id"`
	Result json.RawMessage   `json:"result"`
	Error  *BitcoindRPCError `json:"error"`
}

// BlockchainInfo is the simplifed result of getblockchaininfo from bitcoind RPC
type BlockchainInfo struct {
	Blocks   int     `json:"blocks"`
	Progress float64 `json:"verificationprogress"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

// GetBlockchainInfo returns the result of `getblockchaininfo` from bitcoind
func GetBlockchainInfo(serverURL string) (BlockchainInfo, error) {
	var blockchainInfo BlockchainInfo

	reqBody := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "bitcoind-ctl",
		"method":  "getblockchaininfo",
	}

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return blockchainInfo, err
	}

	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return blockchainInfo, err
	}
	defer resp.Body.Close()

	var respBody BitcoindResponse

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return blockchainInfo, err
	}

	if resp.StatusCode != http.StatusOK {
		if respBody.Error != nil {
			switch respBody.Error.Code {
			case -28:
				return blockchainInfo, ErrBitcoindRPCNotReady
			default:
				return blockchainInfo, errors.New(respBody.Error.Message)
			}
		}

		return blockchainInfo, ErrBitcoindUnknown
	} else {
		err := json.Unmarshal(respBody.Result, &blockchainInfo)
		return blockchainInfo, err
	}
}
