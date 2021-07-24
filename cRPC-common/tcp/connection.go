package tcp

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/assert"
	"io"
	"net"
	"time"
)

type Connection struct {
	IP     string
	socket net.Conn
}

func NewConnection(socket net.Conn) (*Connection, error) {
	host, _, err := net.SplitHostPort(socket.RemoteAddr().String())
	if err != nil {
		return nil, err
	}

	return &Connection{
		socket: socket,
		IP:     "[" + host + "]", // todo: IP 转换
	}, nil
}

func NewConnectionWithHost(host string) (*Connection, error) {

	socket, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return nil, err
	}
	return NewConnection(socket)
}

func (conn *Connection) Close() {
	conn.socket.Close()
}

func (conn *Connection) Send(msg *Message) error {

	msgBuf := bytes.NewBuffer(make([]byte, 0, msgHeaderLength + msg.Header.PayloadLen))

	// header buf
	headerBytes := transMsgToByte(msg)
	err := binary.Write(msgBuf, binary.BigEndian, &headerBytes)
	if err != nil {
		logrus.Errorf("write packet header to buf error [%v]", err)
		return err
	}
	// payload buf
	if msg.Header.PayloadLen != 0 {
		err = binary.Write(msgBuf, binary.BigEndian, msg.Payload)
		if err != nil {
			logrus.Errorf("write packet payload to buf error [%v]", err)
			return err
		}
	}

	err = socketSend(conn.socket, int64(len(msgBuf.Bytes())), msgBuf.Bytes(), 0)
	if err != nil {
		return err
	}
	return nil
}

func socketSend(conn net.Conn, packetLen int64, buffer []byte, timeout time.Duration) error {

	assert.CheckParam(packetLen <= int64(len(buffer)) && timeout >= 0)
	hasWrittenLen := int64(0)
	for hasWrittenLen < packetLen {
		n, err := conn.Write(buffer[hasWrittenLen:packetLen])
		if err != nil {
			logrus.Errorf("conn write buffer error [%v]", err)
			return err
		}
		if n == 0 {
			// 一般不会出现 n == 0 的情况，如果出现，则可认定为不正常状态，打印告警日志，且返回 EOF 即可；
			//  Reference: http://stackoverflow.com/questions/2176443/is-a-return-value-of-0-from-write2-in-c-an-error
			logrus.Warningf("write socket[%s] zero bytes!", conn.RemoteAddr())
			return io.EOF
		} else {
			hasWrittenLen = hasWrittenLen + int64(n)
		}
	}
	return nil
}

func (conn *Connection) Receive(timeout time.Duration) (*Message, error) {

	msg := new(Message)
	// 读取消息头
	headerBuff := make([]byte, msgHeaderLength)
	err := socketReceive(conn.socket, int64(msgHeaderLength), headerBuff, timeout)
	if err != nil {
		logrus.Warningf("read header err, socket : [%s], error [%v]", conn.socket.RemoteAddr(), err)
		return nil, err
	}
	header := transByteToHeader(headerBuff)
	// 根据消息头中 payload 长度读取 payload
	payloadBuff := make([]byte, header.PayloadLen)
	err = socketReceive(conn.socket, int64(header.PayloadLen), payloadBuff, timeout)
	if err != nil {
		logrus.Warningf("read payload err, socket : [%s], error [%v]", conn.socket.RemoteAddr(), err)
		return nil, err
	}
	msg.Header = *header
	msg.Payload = payloadBuff
	return msg, nil
}

func socketReceive(conn net.Conn, packetLen int64, Buffer []byte, timeout time.Duration) (err error) {

	var receiveLen int64
	if timeout <= 0 {
		conn.SetReadDeadline(time.Time{})
	} else {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}
	for receiveLen < packetLen {
		tempLen, err := conn.Read(Buffer[receiveLen:packetLen])
		if err == io.EOF {
			logrus.Warningf("peer socket[%s] exit", conn.RemoteAddr())
			return err
		} else if err != nil {
			logrus.Warningf("read socket[%s] connection error[%v]", conn.RemoteAddr(), err)
			return err
		} else {
			receiveLen += int64(tempLen)
		}
	}
	return err
}
