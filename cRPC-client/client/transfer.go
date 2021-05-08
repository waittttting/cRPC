package client

import (
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/cerr"
	"github.com/waittttting/cRPC-common/snowFlake"
	"github.com/waittttting/cRPC-common/tcp"
	"math/rand"
	"time"
)

type transMessage struct {
	tcpMessage   tcp.Message
	responseChan chan *tcp.Message
}

func (rc *RpcClient) Send(message *tcp.Message) (*tcp.Message, error) {

	// msg 生成对应的 seq id 和 trace id
	// todo: node index
	traceId, err := snowFlake.GenSeqId(0)
	if err != nil {
		return nil, err
	}
	seqId, err := snowFlake.GenSeqId(0)
	if err != nil {
		return nil, err
	}

	message.Header.TraceId = traceId.String()
	message.Header.SeqId = uint64(seqId.Int64())

	rc.scsLock.Lock() // 服务掉线，服务上线都会操作 subConnMap
	conns, ok := rc.subConnMap[message.Header.ServerName]
	rc.scsLock.Unlock()
	if !ok {
		logrus.Error(cerr.ErrServiceNotFound.ErrMsg)
		return nil, cerr.ErrServiceNotFound
	}

	// 随机获取连接 random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(conns))
	conn := conns[index]

	if conn != nil {
		logrus.Error(cerr.ErrConnNotFound.ErrMsg)
		return nil, cerr.ErrConnNotFound
	}

	tMsg := transMessage{

	}
	timer := time.NewTimer(time.Second * 3)
	defer timer.Reset(0)
	var response *tcp.Message
	select {
	case  response = <- tMsg.responseChan:
		// todo: 接收消息
		timer.Stop()
		return response, nil
	case <- timer.C: // todo: time.C 这样写会创建一个协程吗，当有大量的消息时，这样写好还是用 tw 的方式好
		// 超时
		return nil, cerr.ErrTimeOut
	}
}


// 发送消息
func send(conn tcp.Connection, message *tcp.Message) {

	conn.Send(message)
}


