package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/tcp"
)

func (rc *RpcClient) createConnsWithSubServiceInfos() {
	// 遍历订阅服务的信息
	for _, subServerInfo := range rc.cloudConfig.SubServersInfos {
		// 遍历订阅的服务的可用节点
		for _, nodeInfo := range subServerInfo.Infos {
			go rc.createConnWithNodeInfo(nodeInfo)
		}
	}
}

func (rc *RpcClient) createConnWithNodeInfo(nodeInfo string) {

	logrus.Infof("create conn with NodeInfo : [%v]", nodeInfo)

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("createConnWithNodeInfo err : [%v]", err)
		}
	}()

	var gid tcp.GID
	err := json.Unmarshal([]byte(nodeInfo), &gid)
	if err != nil {
		logrus.Errorf("unmarshal node gid err:[%v]", err)
		return
	}

	host := fmt.Sprintf("%s:%v", gid.IP, gid.Port)
	conn, err := rc.createConnAndRegister(host)
	if err != nil {
		logrus.Errorf("create sub server conn and register err:[%v], target serverName:[%s] ip:[%s]", err, gid.ServiceName, gid.IP)
		return
	}

	receiveMsgChan := make(chan *tcp.Message, rc.localConfig.Server.ReceiveTcpMsgChanLen)
	newSc := newServerConn(gid, conn, receiveMsgChan)

	sp := &serverPackage{
		sc: &newSc,
		receiveMsgChan: receiveMsgChan,
	}
	rc.scsLock.Lock()
	defer rc.scsLock.Unlock()
	if servers, ok := rc.subConnMap[gid.ServiceName]; ok {
		if _, ok := servers[gid.String()]; !ok {
			servers[gid.String()] = sp
		}
	} else {
		subServerConns := make(map[string]*serverPackage)
		subServerConns[gid.String()] = sp
		rc.subConnMap[gid.ServiceName] = subServerConns
	}
	newSc.StartLoop(true)
	go handleMsg(receiveMsgChan)
}
