package server

import (
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-config-center/conf"
	"net"
)

type ConfigCenterServer struct {
	config *conf.CCSConf
	ReceiveSocketChan chan *net.TCPConn
}

func NewConfigCenterServer(config *conf.CCSConf) *ConfigCenterServer {
	return &ConfigCenterServer{
		config: config,
		ReceiveSocketChan: make(chan *net.TCPConn, config.Server.ReceiveSocketChanLen),
	}
}

func (ccs *ConfigCenterServer) Start() {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("handle received socket occur panic: %v", r)
			}
		}()
	}()
}