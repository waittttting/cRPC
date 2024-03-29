package server

import (
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/tcp"
	"time"
)

/**
 * @Description: 与其他服务的长连接的结构体
 */
type serverConn struct {
	serviceName    string
	serviceVersion string
	conn           *tcp.Connection
	exit           bool
	gid            *tcp.GID
	heartbeatTime  time.Time
	ccs            *ControlCenterServer
	redisKey       string
	redisValue     string
}

func (sc *serverConn) loop() {

	// todo: 上线记录
	sc.open()
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("serviceConn loop panic serviceName: %v, serviceVersion: %v, err: %v", sc.serviceName, sc.serviceVersion, err)
		}
	}()
	for !sc.exit {
		msg, err := sc.conn.Receive(0 * time.Second)
		if err != nil {
			logrus.Errorf("receive msg in loop occurred err : %v", err)
			sc.exit = true
		}
		// 刷新心跳
		sc.handleHeartbeat(msg)
	}
	logrus.Infof("complete loop gid: [%s]", sc.gid.String())
	// 关闭 sc
	sc.close()
}

func (sc *serverConn) open() {
	logrus.Infof("service conn onLine, git = [%v]", sc.gid)
	redisPubOnLine(sc.gid)
}

/**
 * @Description:
 * @receiver sc
 * @param cause
 */
func (sc *serverConn) letScOffLine(cause string) {
	logrus.Infof("service conn offLine, cause:[%v]", cause)
	sc.exit = true
	sc.conn.Close() // 当调用 sc.conn.Close() 时，sc.conn.Receive 回直接报错，此时 exit 已经是 true，结束 loop 循环
}

func (sc *serverConn) handleHeartbeat(msg *tcp.Message) {

	sc.heartbeatTime = time.Now()
	sc.ccs.refreshHeartbeat(sc)
	// todo: 刷新心跳逻辑是否写全
}

/**
 * @Description: sc.close 与 sc.offLine 语意的区别: close 是清理 ccs 内 sc 对应的信息，LetScOffLine 是外部 让 sc 关闭
 * @receiver sc
 */
func (sc *serverConn) close() {

	// 删除 ccs 内的 sc 信息
	sc.ccs.lock.Lock()
	defer sc.ccs.lock.Unlock()
	curS := sc.ccs.clientMap[sc.gid.ServiceName]
	if curS == nil {
		logrus.Error("not find sc from ccs.clientMap")
	} else {
		index := 0
		for i := 0; i < len(curS); i++ {
			if curS[i] == sc {
				index = i
			}
		}
		curS = append(curS[:index], curS[index+1:]...)
	}

	sc.ccs.deleteHeartbeat(sc)                                      // 清除 time wheel
	err := RedisOp(sc.redisKey, sc.redisValue, redisOpSRemServerIp) // 清除 redis 信息
	if err != nil {
		// 报警，人工介入
		logrus.Errorf("delete client ip from redis error:[%v]", err)
	}
	// 通知其他 ccs 上 <订阅了本服务的服务>，本服务下线
	redisPubOffline(sc.gid)
}
