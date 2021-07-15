// SPDX-License-Identifier: ISC
// Copyright (c) 2019-2021 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeConfig *rest.Config

func main() {
	runtime := flag.String("runtime", "k8s", "bitcoind runtime environment (k8s|local)")

	var configFile string
	flag.StringVar(&configFile, "c", "./config.yaml", "[optional] path of configuration file")
	flag.StringVar(&configFile, "config", "./config.yaml", "[optional] path of configuration file")
	flag.Parse()

	LoadConfig(configFile)

	{
		var err error
		if viper.GetBool("k8s.use_local_context") {
			var kubeConfigFile *string
			if home := homedir.HomeDir(); home != "" {
				kubeConfigFile = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeConfigFile = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}

			kubeConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfigFile)
			if err != nil {
				panic(err)
			}
		} else {
			kubeConfig, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}
	}

	var bitcoindCtl BitcoindCtl
	if *runtime == "local" {
		bitcoindCtl = NewBitcoindLocalCtl(viper.GetString("owner_did"),
			viper.GetString("bitcoind_network"),
			viper.GetString("bitcoind_endpoint"))
	} else {
		bitcoindCtl = NewBitcoindK8SCtl(kubeConfig,
			viper.GetString("k8s.namespace"),
			viper.GetString("owner_did"),
			viper.GetString("bitcoind_network"),
			viper.GetString("bitcoind_endpoint"))
	}

	route := gin.Default()

	route.POST("/bitcoind/start", func(c *gin.Context) {
		if err := bitcoindCtl.StartBitcoind(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": 1})
	})

	route.POST("/bitcoind/stop", func(c *gin.Context) {
		if err := bitcoindCtl.StopBitcoind(c, false); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": 1})
	})

	route.GET("/bitcoind/status", func(c *gin.Context) {
		status, err := bitcoindCtl.GetBitcoindStatus(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, status)
	})

	route.Run(viper.GetString("port"))
}
