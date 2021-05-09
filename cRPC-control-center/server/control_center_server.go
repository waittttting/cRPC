package server

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/model"
	"github.com/waittttting/cRPC-common/tcp"
	"github.com/waittttting/cRPC-common/twheel"
	"github.com/waittttting/cRPC-control-center/conf"
	"net"
	"strconv"
	"sync"
	"time"
)

type ControlCenterServer struct {

	config *conf.CCSConf
	// 接收 socket 的队列
	receiveSocketChan chan *net.TCPConn
	// 存储注册到本节点的服务的 名字和 连接的关系
	clientMap map[string][]*serviceConn
	lock sync.RWMutex
	receiveServiceConnChan chan *serviceConn
	timeWheel *twheel.TimeWheel
	heartbeatTimeoutChan chan interface{}
	redisMsgChan <-chan *redis.Message
}

func NewControlCenterServer(config *conf.CCSConf, ch <-chan *redis.Message) *ControlCenterServer {


	heartbeatChan := make(chan interface{}, config.TimeWheel.NoticeChanLen)
	tw, err := twheel.New(config.TimeWheel.Cap, heartbeatChan)
	if err != nil {
		logrus.Fatalf("init time wheel err %v", err)
	}
	return &ControlCenterServer{
		config: config,
		receiveSocketChan: make(chan *net.TCPConn, config.Server.ReceiveSocketChanLen),
		receiveServiceConnChan: make(chan *serviceConn, config.Server.ReceiveServiceConnChanLen),
		clientMap: make(map[string][]*serviceConn),
		timeWheel: tw,
		heartbeatTimeoutChan: heartbeatChan,
		redisMsgChan: ch,
	}
}

// 控制中心集群中 各个节点如何进行通信 ？ Gossip 协议？
// https://github.com/hashicorp/memberlist Gossip 的一个实现
// 各个服务节点 通过 dns 获取 控制中心的 IP
// 目前使用 k8s service 做长连接的转发

func (ccs *ControlCenterServer) Start() {

	// ccs 集群通信
	ccs.clusterCommunication()
	// 启动心跳时间轮
	ccs.heartbeatTimeWheelStart()
	// 处理接收到的 serviceConn
	ccs.handleServiceConn()
	// 处理接收到的 socket
	ccs.handleSocket()
	// 启动端口监听
	ccs.acceptSocket()
	// 处理监听的 redis 队列
	ccs.handleRedisMsg()
	logrus.Info("control server started")
}

// 集群通信
func (ccs *ControlCenterServer) clusterCommunication() {
	// todo: 集群内部依赖 Gossip 协议 进行通信
}

func (ccs *ControlCenterServer) acceptSocket() {

	tcp.AcceptSocket(strconv.Itoa(ccs.config.Server.TcpPort), ccs.receiveSocketChan, 3 * time.Second)
}

func (ccs *ControlCenterServer)handleSocket() {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("handleSocket goroutine error %v", err)
			}
		}()
		for socket := range ccs.receiveSocketChan {
			logrus.Infof("socket addr = %v", socket.RemoteAddr())
			ccs.checkRegisterMsg(socket)
		}
	}()
}

func (ccs *ControlCenterServer)checkRegisterMsg(socket *net.TCPConn) {

	// 在注册的过程中，如果出现了 err 直接 调用 socket.Close(),
	// 注册到 redis 失败 结束本次注册，关闭 socket，等待 client 重连
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

	newSc := &serviceConn{
		serviceName: msg.Header.ServerName,
		serviceVersion: msg.Header.ServerVersion,
		conn: conn,
		gid: tcp.NewGid(msg.Header.ServerName, msg.Header.ServerVersion, conn.IP, portConf.Port),
		ccs: ccs, // todo: 是否可以循环引用
	}

	// 注册 Service 的 IP 注册到 redis, key = serviceName + serviceVersion
	// Client 在注册与 Server 的 tcp 连接的时候，会根据 serviceName + serviceVersion 找到对应的IP 列表

	jgid, err := json.Marshal(newSc.gid)
	if err != nil {
		logrus.Errorf("marshal gid err:[%v]", socket.RemoteAddr())
		socket.Close()
		return
	}

	newSc.redisKey = newSc.gid.ServiceName
	newSc.redisValue = string(jgid)

	err = RedisOp(newSc.redisKey, newSc.redisValue, redisOpSAddServerIp)
	if err != nil {
		logrus.Errorf("register to redis error: [%v]", err)
		socket.Close()
		return
	}

	// 回复注册成功消息
	if err = conn.Send(tcp.MsgRegisterPong()); err != nil {
		logrus.Errorf("send MsgRegisterPong error: [%v]", err)
		err = RedisOp(newSc.redisKey, newSc.redisValue, redisOpSRemServerIp)
		if err != nil {
			// 报警，人工介入
			logrus.Errorf("delete client ip from redis error: [%v]", err)
		}
		socket.Close()
		return
	}

	// 虽然此处的 map 操作都是顺序的，但是可能会有 server 与 client 的 tcp 断开的情况，在这种情况下，也会写操作这个 map
	ccs.lock.Lock()
	defer ccs.lock.Unlock()
	if v, ok := ccs.clientMap[newSc.serviceName]; !ok {
		scs := []*serviceConn{newSc}
		ccs.clientMap[newSc.gid.ServiceName] = scs
	} else {
		v = append(v, newSc)
	}

	// 添加到心跳时间轮
	ccs.addToHeartbeatTimeWheel(newSc)

	// 放入到 接收 serviceConn 的队列中
	timer := time.NewTimer(1 * time.Second)
	select {
	case ccs.receiveServiceConnChan <- newSc:
	case <- timer.C:
		logrus.Errorf("send to receiveServiceConnChan timeout %s", newSc.gid.String())
	}
}

func (ccs *ControlCenterServer) handleServiceConn() {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("handleServiceConn goroutine error %v", err)
			}
		}()
		for sc := range ccs.receiveServiceConnChan {
			// 新开一个协程开始等待 sc 的消息
			go sc.loop()
		}
	}()
}


// todo: 没有使用 Gossip 导致的引入了 redis 的订阅, 需要保证消息的可靠性，增加了架构的复杂性，使用 Gossip 则，则会通过 Gossip 协议主键通知到其他的服务
func (ccs *ControlCenterServer) handleRedisMsg()  {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("handleServiceConn goroutine error %v", err)
			}
		}()
		for msg := range ccs.redisMsgChan {
			// todo: 收到消息后看那些服务订阅了本条消息的服务，然后将这条消息通过 tcp 长连接推送到这些服务上
			println(msg.Channel, msg.Payload)
		}
	}()
}


