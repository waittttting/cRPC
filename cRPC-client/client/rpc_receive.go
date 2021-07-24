package client

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/model"
	"github.com/waittttting/cRPC-common/tcp"
	"runtime/debug"
)

// 处理本服务接收到的 socket
func (rc *RpcClient) handleReceivedSocket()  {

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("handleSocket goroutine error %v", err)
		}
	}()
	for socket := range rc.receiveSocketChan {

		logrus.Infof("socket addr = %v", socket.RemoteAddr())
		conn, err := tcp.NewConnection(socket)
		if err != nil {
			logrus.Errorf("NewConnection error address:[%v], err:[%v]", socket.RemoteAddr(), err)
			socket.Close() // 调用 socket.Close(), Client 调用 read 时会直接报错 connection reset by peer
			return
		}
		msg, err := conn.Receive(0)
		if err != nil {
			logrus.Errorf("receive msg error address:[%v], err:[%v]", socket.RemoteAddr(), err)
			socket.Close()
			return
		}

		var portConf model.PortConfig
		err = json.Unmarshal(msg.Payload, &portConf)
		if err != nil {
			logrus.Errorf("unmarshal server port err:[%v]", socket.RemoteAddr())
			socket.Close()
			return
		}

		// 回复注册成功消息
		if err = conn.Send(tcp.MsgRegisterPong()); err != nil {
			logrus.Errorf("send MsgRegisterPong error: [%v]", err)
			if err != nil {
				// 报警，人工介入
				logrus.Errorf("delete client ip from redis error: [%v]", err)
			}
			socket.Close()
			return
		}

		// 创建 serverConn
		gid := tcp.NewGid(msg.Header.ServerName, msg.Header.ServerVersion, conn.IP, portConf.Port)
		receiveMsgChan := make(chan *tcp.Message, rc.localConfig.Server.ReceiveTcpMsgChanLen)
		newSc := newServerConn(*gid, conn, receiveMsgChan, rc)
		sp := &serverPackage{
			sc: &newSc,
			receiveMsgChan: receiveMsgChan,
		}

		rc.scsLock.Lock()
		if servers, ok := rc.subConnMap[gid.ServiceName]; ok {
			if _, ok := servers[gid.String()]; !ok {
				servers[gid.String()] = sp
			} else {
				return
			}
		} else {
			subServerConns := make(map[string]*serverPackage)
			subServerConns[gid.String()] = sp
			rc.subConnMap[gid.ServiceName] = subServerConns
		}
		rc.scsLock.Unlock()
		newSc.StartLoop(false)

		go handleMsg(receiveMsgChan)
	}
}

/**
 * @Description:
 * @param receiveMsgChan
 */
func handleMsg(receiveMsgChan chan *tcp.Message) {
	logrus.Info("begin to handle msg")
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			logrus.Errorf("handleMsg goroutine error %v", err)
		}
	}()
	for msg := range receiveMsgChan {
		logrus.Infof("received msg from : [%v]", msg)
	}
}