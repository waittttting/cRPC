package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/waittttting/cRPC-common/model"
)

const (
	// 业务消息
	MsgCodeTypeBus = 1
	// 框架消息
	MsgCodeTypeKit = 2
	// client 注册到 control center
	msgKitRegisterPing = "msgKitRegisterPing"
	// control center 回复 client 的注册消息
	msgKitRegisterPong = "msgKitRegisterPong"
	// 心跳
	msgKitHeartbeat = "msgKitHeartbeat"
	// 数据帧头长度
	msgHeaderLength uint16 = 400
)

type Message struct {
	Header  Header
	Payload []byte
	C       chan Message
}

type Header struct {
	// 用户id
	Uid string
	// token
	Token string
	// IP
	Ip string
	// 流量标记(用于灰度等)
	Tag uint8
	// 全链路压测
	MTest string
	// 会话ID
	SessionId string
	// 消息类型 (框架消息/业务消息)
	MsgCodeType uint8
	// 消息ID（用于对应消息）
	SeqId uint64
	// traceId 链路追踪
	TraceId string
	// 方法名
	ServerMethod string
	// 目标服务名
	ServerName string
	// 目标服务版本
	ServerVersion string
	// body 长度，在读数据的时候有用
	PayloadLen uint16
	// 主域ID（指代某个客户）
	DomainId uint8
	// AppID 某个客户的某个业务
	AppId uint8
	// 扩展字段
	Expend string
}

/**
 * @Description: 将 Message 转换成 byte 数组
 * @param msg
 * @return [msgHeaderLength]byte
 */
func transMsgToByte(msg *Message) [msgHeaderLength]byte {

	seqIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqIdBytes, msg.Header.SeqId)
	payloadLenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadLenBytes, msg.Header.PayloadLen)

	var headerBytes [msgHeaderLength]byte
	copy(headerBytes[:0+32], msg.Header.Uid)
	copy(headerBytes[32:32+32], msg.Header.Token)
	copy(headerBytes[64:64+15], msg.Header.Ip)
	copy(headerBytes[79:79+1], []uint8{msg.Header.Tag})
	copy(headerBytes[80:80+32], msg.Header.MTest)
	copy(headerBytes[112:112+32], msg.Header.SessionId)
	copy(headerBytes[144:144+1], []uint8{msg.Header.MsgCodeType})
	copy(headerBytes[145:145+64], seqIdBytes)
	copy(headerBytes[209:209+64], msg.Header.TraceId)
	copy(headerBytes[273:273+32], msg.Header.ServerMethod)
	copy(headerBytes[305:305+32], msg.Header.ServerName)
	copy(headerBytes[337:337+32], msg.Header.ServerVersion)
	copy(headerBytes[369:369+1], []uint8{msg.Header.DomainId})
	copy(headerBytes[370:370+1], []uint8{msg.Header.AppId})
	copy(headerBytes[371:371+2], payloadLenBytes)
	copy(headerBytes[373:373+27], msg.Header.Expend)
	return headerBytes
}

func transByteToHeader(headerBuff []byte) *Header {

	header := new(Header)
	header.MsgCodeType = headerBuff[144]
	header.SeqId = binary.BigEndian.Uint64(headerBuff[145 : 145+64])
	header.ServerMethod = string(bytes.Trim(headerBuff[273:273+32], "\x00"))
	header.ServerName = string(bytes.Trim(headerBuff[305:305+32], "\x00"))
	header.ServerVersion = string(bytes.Trim(headerBuff[337:337+32], "\x00"))
	header.PayloadLen = binary.BigEndian.Uint16(headerBuff[371 : 371+2])
	return header
}

func MsgRegisterPing(severName string, serverVersion string, serverPort int) (*Message, error) {

	port := model.PortConfig{
		Port: serverPort,
	}
	portJson, err := json.Marshal(port)
	if err != nil {
		return nil, err
	}
	header := Header{
		PayloadLen:    uint16(len(portJson)),
		MsgCodeType:   MsgCodeTypeKit,
		ServerMethod:  msgKitRegisterPing,
		ServerVersion: serverVersion,
		ServerName:    severName,
	}
	return &Message{
		Header:  header,
		Payload: portJson,
	}, nil
}

func MsgRegisterPong() *Message {

	header := Header{
		PayloadLen:   0,
		MsgCodeType:  MsgCodeTypeKit,
		ServerMethod: msgKitRegisterPong,
	}
	return &Message{
		Header:  header,
		Payload: nil,
	}
}

func MsgHeartbeat(serverName, serverVersion string) *Message {

	header := Header{
		ServerName: serverName,
		ServerVersion: serverVersion,
		PayloadLen:   0,
		MsgCodeType:  MsgCodeTypeKit,
		ServerMethod: msgKitHeartbeat,
	}
	return &Message{
		Header:  header,
		Payload: nil,
	}
}
