package client

import (
	"github.com/waittttting/cRPC-common/tcp"
	"time"
)

type clientConn struct {
	conn *tcp.Connection
}

func (cc *clientConn) heartbeat() {
	ticker := time.Tick(3 * time.Second)
	for {
		<- ticker
		// todo: 连接关闭时，停止发送心跳
		cc.conn.Send(tcp.MsgHeartbeat())
	}
}
