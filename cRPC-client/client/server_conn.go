package client

import (
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/tcp"
	"time"
)

type serverConn struct {

	gid tcp.GID
	conn *tcp.Connection
	// 是否断开 true 断开，false 未断开
	exit bool
	receiveMsgChan chan *tcp.Message
}

func newServerConn(
	gid tcp.GID,
	connection *tcp.Connection,
	msgChan chan *tcp.Message) serverConn {
	
	sc := serverConn{
		gid: gid,
		conn: connection,
		receiveMsgChan: msgChan,
	}
	return sc
}


func (sc *serverConn) StartLoop(heartbeat bool) {
	if heartbeat {
		go sc.heartbeat()
	}
	go sc.loop()
}

func (sc *serverConn) heartbeat() {

	ticker := time.Tick(3 * time.Second)
	for !sc.exit {
		<- ticker
		err := sc.conn.Send(tcp.MsgHeartbeat())
		if err != nil {
			logrus.Errorf("heartbeat err : [%v]", err)
		}
		logrus.Infof("...send heartbeat... [%v]", sc.gid.ServiceName)
	}
}

func (sc *serverConn) loop() {

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("serverConn loop panic [%v], [%v]", sc.gid.String(), err)
		}
	}()

	for !sc.exit {
		msg, err := sc.conn.Receive(0 * time.Second)
		if err != nil {
			logrus.Errorf("receive msg in loop occurred err : %v", err)
			sc.exit = true
		}
		logrus.Infof("server_comm loop header: %s, %v", msg.Header.ServerName, msg.Header.ServerName)
		timer := time.NewTimer( 2 * time.Second)
		select {
		case sc.receiveMsgChan <- msg:
		case <- timer.C:
			logrus.Errorf("send msg to ControlCenterMsgChan time out %v", err)
		}
	}
	logrus.Infof("complete loop gid: %s", sc.gid.String())
}
