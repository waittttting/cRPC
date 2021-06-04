package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-client/conf"
	"github.com/waittttting/cRPC-common/model"
	"github.com/waittttting/cRPC-common/tcp"
	"net"
	"strconv"
	"sync"
	"time"
)

const controlCenterMsgChanLen = 100

type RpcClient struct {
	// 本地配置
	config *conf.LocalConf
	// 云端配置
	cloudConfig *model.CloudConfigInfo
	// 连接控制中心的长连接
	controlCenterConn *serverConn
	// 接收 socket 的队列
	receiveSocketChan chan *net.TCPConn
	// 订阅的服务的长连接
	subConnMap map[string][]*tcp.Connection
	// subConns 的锁
	scsLock sync.Mutex
}

func NewRpcClient(config *conf.LocalConf) *RpcClient {
	return &RpcClient{
		config: config,
		receiveSocketChan: make(chan *net.TCPConn, config.Server.ReceiveSocketChanLen),
		subConnMap: map[string][]*tcp.Connection{},
	}
}

func (rc *RpcClient) Start() {

	// 如果本服务提供 tcp 的服务，开启 tcp 端口
	if rc.config.Server.RpcOpen {
		// 开启长连接端口
		rc.startReceiveSocket()
		rc.handleSocket()
	}
	// 发送 http 请求获取 config file
	cloudConfig, err := rc.getServerConfig()
	if err != nil {
		logrus.Fatalf("getServerConfig err:[%v]", err)
	}
	rc.cloudConfig = cloudConfig
	// 连接控制中心
	configCenterHost := fmt.Sprintf("%s:%s", rc.cloudConfig.ControlCenterAddr.Host, rc.cloudConfig.ControlCenterAddr.TcpPort)
	conn, err := rc.createAndRegister(configCenterHost)
	if err != nil {
		logrus.Fatalf("create config center conn and register err:[%v]", err)
		return
	}
	logrus.Info("register success")
	sc := newServerConn("controlCenter", "0.0.1", controlCenterMsgChanLen, conn)
	rc.controlCenterConn = &sc
	// 处理 control center 的消息
	go rc.handleControlCenterMsg()
	// 开始读 消息
	go rc.controlCenterConn.StartLoop(true)

	// 创建和订阅的服务的长连接
	go rc.createConnsWithSubServiceInfos()
}


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

	var gid tcp.GID
	err := json.Unmarshal([]byte(nodeInfo), &gid)
	if err != nil {
		logrus.Errorf("unmarshal node gid err:[%v]", err)
		return
	}
	host := fmt.Sprintf("%s:%v", gid.IP, gid.Port)

	conn, err := rc.createAndRegister(host)
	if err != nil {
		logrus.Errorf("create sub server conn and register err:[%v], target serverName:[%s] ip:[%s]", err, gid.ServiceName, gid.IP)
		return
	}
	rc.scsLock.Lock()
	defer rc.scsLock.Unlock()
	if v, ok := rc.subConnMap[gid.ServiceName]; ok {
		// todo: 遍历
		v = append(v, conn)
	} else {
		subConns := make([]*tcp.Connection, 0)
		subConns = append(subConns, conn)
		rc.subConnMap[gid.ServiceName] = subConns
	}
}

/**
 * @Description: 创建和目标服务的长连接，并且注册
 * @receiver rc
 * @param host
 * @return *tcp.Connection
 * @return error
 */
func (rc *RpcClient) createAndRegister(host string) (*tcp.Connection, error) {

	conn, err := tcp.NewConnectionWithHost(host)
	if err != nil {
		logrus.Errorf("create sub service err:[%v]", err)
		return nil, err
	}
	// 创建注册消息
	p, err := tcp.MsgRegisterPing(rc.config.Client.ServerName, rc.config.Client.ServerVersion, rc.config.Server.TcpPort)
	if err != nil {
		logrus.Errorf("init MsgRegisterPing failed err:[%v]", err)
		return nil, err
	}
	// 发送自身的信息到订阅的服务
	err = conn.Send(p)
	if err != nil {
		logrus.Errorf("send server message error:[%v]", err)
		return nil, err
	}

	// 接收回复的消息
	_, err = conn.Receive(0 * time.Second)
	if err != nil {
		logrus.Errorf("register failed err:[%v]", err)
		return nil, err
	}
	return conn, nil
}


func (rc *RpcClient) startReceiveSocket() {

	tcp.AcceptSocket(strconv.Itoa(rc.config.Server.TcpPort), rc.receiveSocketChan, 3 * time.Second)
}

// 处理本服务接收到的 socket
func (rc *RpcClient) handleSocket()  {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("handleSocket goroutine error %v", err)
			}
		}()
		for socket := range rc.receiveSocketChan {
			logrus.Infof("socket addr = %v", socket.RemoteAddr())
			// todo: 处理新接收的（订阅本服务的） socket
		}
	}()
}


func (rc *RpcClient) handleControlCenterMsg() {

	for msg := range rc.controlCenterConn.receiveMsgChan {
		// todo: 处理控制中心发送的消息~
		print(msg)
	}
}

