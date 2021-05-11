package main

import (

	"cx-mp-gateway/crisk"
	"cx-mp-gateway/cweb"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-client/client"
	"github.com/waittttting/cRPC-client/conf"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()
	var config conf.LocalConf
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		logrus.Fatal(err)
	}
	// cRPC client
	rpcClient := client.NewRpcClient(&config)
	rpcClient.Start()
	// 风控启动
	crisk.GRS.Start()
	// web
	ws := cweb.NewWebServer(rpcClient)
	ws.Start(config)
}
