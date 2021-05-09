package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-control-center/conf"
	"github.com/waittttting/cRPC-control-center/cweb"
	"github.com/waittttting/cRPC-control-center/server"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()

	var config conf.CCSConf
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		logrus.Fatal(err)
	}

	// 初始化 Redis
	server.RedisInit(config)
	ch, err := server.RedisSubServerOnLine()
	if err != nil {
		logrus.Fatalf("sub redis err:[%v]", err)
	}
	// todo: 获取所有服务的订阅关系
	// todo: 订阅所有服务的 key

	ccs := server.NewControlCenterServer(&config, ch)
	ccs.Start()

	ws := cweb.NewWebServer()
	ws.Start(config)
}