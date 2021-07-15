// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// podID returns a unique id for an Autonomy pod by account number
func podID(accountNumber string) string {
	first := sha256.Sum256([]byte(accountNumber))
	second := sha256.Sum256(first[:])

	return fmt.Sprintf("%x", second[:16])
}

// roundDigits returns a number rounded to digits precision after the decimal point.
func roundDigits(number float64, digits int) float64 {
	rate := math.Pow10(digits)
	return math.Round(number*rate) / rate
}

// BitcoindK8SCtl is a bitcoind controller for controlling bitcoind over kubernetes cluster
type BitcoindK8SCtl struct {
	k8sClient    *kubernetes.Clientset
	bitcoindURL  string
	network      string
	namespace    string
	did          string
	networkInfix string
	podID        string
}

func NewBitcoindK8SCtl(kubeConfig *rest.Config, namespace, did, bitcoindNetwork, bitcoindURL string) *BitcoindK8SCtl {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	networkInfix := ""
	if bitcoindNetwork == "mainnet" {
		networkInfix = "-mainnet"
	}

	return &BitcoindK8SCtl{
		k8sClient:    clientset,
		namespace:    namespace,
		bitcoindURL:  bitcoindURL,
		network:      bitcoindNetwork,
		did:          did,
		networkInfix: networkInfix,
		podID:        podID(did),
	}
}

// updateBitcoindReplica updates the replicas of the statefulset of a autonomy pod
func (ctl *BitcoindK8SCtl) updateBitcoindReplica(ctx context.Context, count int) error {
	payload := []map[string]interface{}{{
		"op":    "replace",
		"path":  "/spec/replicas",
		"value": count,
	}}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = ctl.k8sClient.AppsV1().StatefulSets(ctl.namespace).Patch(ctx, fmt.Sprintf("bitcoind%s-%s", ctl.networkInfix, ctl.podID), types.JSONPatchType, payloadBytes, v1.PatchOptions{})
	return err
}

// GetBitcoindStatus returns the status of a pod
func (ctl *BitcoindK8SCtl) GetBitcoindStatus(ctx context.Context) (BitcoindStatus, error) {

	var status BitcoindStatus

	r, err := ctl.k8sClient.AppsV1().StatefulSets(ctl.namespace).Get(ctx, fmt.Sprintf("bitcoind%s-%s", ctl.networkInfix, ctl.podID), v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return status, ErrBitcoindStopped
		} else {
			log.WithError(err).Error("fail to check statefulset for bitcoind")
			return status, ErrReadBitcoindStatus
		}
	}

	if r.Status.Replicas == 0 {
		return status, ErrBitcoindStopped
	}

	if r.Status.ReadyReplicas != 1 {
		return status, ErrBitcoindProcNotReady
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

// StartBitcoind updates the replicas of the statefulset of a autonomy pod to 1
func (ctl *BitcoindK8SCtl) StartBitcoind(ctx context.Context) error {
	return ctl.updateBitcoindReplica(ctx, 1)
}

// StopBitcoind updates the replicas of the statefulset of a autonomy pod to 0
func (ctl *BitcoindK8SCtl) StopBitcoind(ctx context.Context, force bool) error {
	return ctl.updateBitcoindReplica(ctx, 0)
}
