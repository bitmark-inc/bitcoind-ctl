// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

// BitcoindLocalCtl is a bitcoind controller for controlling bitcoind over local machine
type BitcoindLocalCtl struct {
	bitcoindURL  string
	network      string
	did          string
	bitcoindProc *os.Process
}

func NewBitcoindLocalCtl(did, bitcoindNetwork, bitcoindURL string) *BitcoindLocalCtl {
	return &BitcoindLocalCtl{
		bitcoindURL: bitcoindURL,
		network:     bitcoindNetwork,
		did:         did,
	}
}

// GetBitcoindStatus returns the status of a pod
func (ctl *BitcoindLocalCtl) GetBitcoindStatus(ctx context.Context) (BitcoindStatus, error) {
	var status BitcoindStatus

	if ctl.bitcoindProc == nil {
		return status, ErrBitcoindStopped
	}

	blockchainInfo, err := GetBlockchainInfo(ctl.bitcoindURL)
	if err != nil {
		return status, err
	}

	return BitcoindStatus{
		BestBlock:    blockchainInfo.Blocks,
		SyncProgress: roundDigits(blockchainInfo.Progress, 5),
	}, nil
}

// StartBitcoind spawned a bitcoind process using local executable bitcoind
func (ctl *BitcoindLocalCtl) StartBitcoind(ctx context.Context) error {

	cmdArgs := []string{}

	if ctl.network != "mainnet" {
		cmdArgs = append(cmdArgs, "-testnet")
	}

	if configPath := viper.GetString("local.bitcoind_conf_path"); configPath != "" {
		cmdArgs = append(cmdArgs, fmt.Sprint("-conf=", configPath))
	}

	cmd := exec.Command(viper.GetString("local.bitcoind_path"), cmdArgs...)
	if err := cmd.Start(); err != nil {
		return err
	}

	ctl.bitcoindProc = cmd.Process
	return nil
}

// StopBitcoind stops the spawned bitcoind process
func (ctl *BitcoindLocalCtl) StopBitcoind(ctx context.Context, force bool) error {
	if ctl.bitcoindProc == nil {
		return fmt.Errorf("bitcoind process not found")
	}

	if force {
		if err := ctl.bitcoindProc.Kill(); err != nil {
			return err
		}
	} else {
		if err := ctl.bitcoindProc.Signal(os.Interrupt); err != nil {
			return err
		}
	}

	s, err := ctl.bitcoindProc.Wait()
	if err != nil {
		return nil
	}

	if !s.Exited() {
		return fmt.Errorf("bitcoind process is not yet closed")
	}

	ctl.bitcoindProc = nil
	return nil
}
