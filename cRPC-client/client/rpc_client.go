package client

import (
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

const controlCenterServiceName = "controlCenterService"

const controlCenterServiceVersion = "v0.0.1"

type RpcClient struct {
	// 本地配置
	localConfig *conf.LocalConf
	// 云端配置
	cloudConfig *model.CloudConfigInfo
	// 连接控制中心的长连接
	controlCenterConn *serverConn
	// 接收 控制中心消息的队列
	receiveControlCenterMsgChan chan *tcp.Message
	// 接收 socket 的队列
	receiveSocketChan chan *net.TCPConn
	// 订阅的服务的长连接
	subConnMap map[string]map[string]*serverPackage
	// subConnMap 的锁
	scsLock sync.Mutex
}

type serverPackage struct {
	sc *serverConn
	receiveMsgChan chan *tcp.Message
}

func NewRpcClient(config *conf.LocalConf) *RpcClient {
	return &RpcClient{
		localConfig: config,
		receiveSocketChan: make(chan *net.TCPConn, config.Server.ReceiveSocketChanLen),
		receiveControlCenterMsgChan: make(chan *tcp.Message, controlCenterMsgChanLen),
		subConnMap: make(map[string]map[string]*serverPackage),
	}
}

func (rc *RpcClient) Start() {

	// 如果本服务提供 tcp 的服务，开启 tcp 端口
	if rc.localConfig.Server.RpcOpen {
		// 开启长连接端口
		rc.startReceiveSocket()
		go rc.handleReceivedSocket()
	}

	// 发送 http 请求获取 config file
	cloudConfig, err := rc.getServerConfig()
	if err != nil {
		logrus.Fatalf("getServerConfig err:[%v]", err)
	}
	rc.cloudConfig = cloudConfig
	// 1. 连接控制中心
	configCenterHost := fmt.Sprintf("%s:%s", rc.cloudConfig.ControlCenterAddr.Host, rc.cloudConfig.ControlCenterAddr.TcpPort)
	conn, err := rc.createConnAndRegister(configCenterHost)
	if err != nil {
		logrus.Fatalf("create config center conn and register err:[%v]", err)
		return
	}
	logrus.Infof("register success, cloud config = [%v]", rc.cloudConfig)

	ccsPort, err := strconv.Atoi(rc.cloudConfig.ControlCenterAddr.TcpPort)
	if err != nil {
		logrus.Fatalf("ControlCenterAddr tcp port a to i err:[%v]", err)
		return
	}
	gid := tcp.NewGid(
		controlCenterServiceName,
		controlCenterServiceVersion,
		rc.cloudConfig.ControlCenterAddr.Host,
		ccsPort,
		)

	sc := newServerConn(*gid, conn, rc.receiveControlCenterMsgChan, rc)

	rc.controlCenterConn = &sc
	// 1.1 处理 control center 的消息
	go rc.handleControlCenterMsg()
	// 1.2 开始 ccs 的 loop
	go rc.startControlCenterConnLoop()
	// 2. 创建和订阅的服务的长连接
	go rc.createConnsWithSubServiceInfos()
}

func (rc *RpcClient) handleControlCenterMsg() {

	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorf("handleControlCenterMsg panic, err : [%v]", err)
		}
	}()

	for msg := range rc.receiveControlCenterMsgChan {
		// todo: 处理控制中心发送的消息~
		print(msg)
	}
}

func (rc *RpcClient) startControlCenterConnLoop() {

	rc.controlCenterConn.StartLoop(true)
}

/**
 * @Description: 创建和目标服务的长连接，并且注册
 * @receiver rc
 * @param host
 * @return *tcp.Connection
 * @return error
 */
func (rc *RpcClient) createConnAndRegister(host string) (*tcp.Connection, error) {

	conn, err := tcp.NewConnectionWithHost(host)
	if err != nil {
		logrus.Errorf("create sub service err:[%v]", err)
		return nil, err
	}
	// 创建注册消息
	p, err := tcp.MsgRegisterPing(
		rc.localConfig.Client.ServerName,
		rc.localConfig.Client.ServerVersion,
		rc.localConfig.Server.TcpPort)

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

func (rc *RpcClient) scLoopEnd(gid tcp.GID)  {

	fmt.Println(" ------ scLoopEnd 1------ ")

	if gid.ServiceName == controlCenterServiceName {
		// todo
	} else {
		rc.scsLock.Lock()
		defer rc.scsLock.Unlock()
		fmt.Println(" ------ scLoopEnd 2 ------ ")
		if servers, ok := rc.subConnMap[gid.ServiceName]; ok {
			if sp, ok := servers[gid.String()]; ok {
				close(sp.receiveMsgChan)
				delete(servers, gid.String())
				logrus.Infof(" server_conn close [%v]", gid.String())
			} else {
				logrus.Errorf("sp not find when server_conn loop over, server name = [%v]", gid.ServiceName)
			}
		} else {
			logrus.Errorf("service map not find when server_conn loop over, server name = [%v]", gid.ServiceName)
		}
	}
}

// 开启 tcp 端口
func (rc *RpcClient) startReceiveSocket() {

	tcp.AcceptSocket(strconv.Itoa(rc.localConfig.Server.TcpPort), rc.receiveSocketChan, 3 * time.Second)
}


