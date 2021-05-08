package tcp

import "testing"

func TestMessageTrans(t *testing.T) {

	msg := Message{
		Header: Header{
			Uid:           "1",
			Token:         "duiuwddanld1221ndsaud1d0jczs91ld9aj",
			Ip:            "",
			Tag:           1,
			MTest:         "2",
			SessionId:     "1",
			MsgCodeType:   1,
			SeqId:         20897,
			TraceId:       "duiuwddanld1221ndsaud1d0jczs91ld9aj",
			ServerMethod:  "register",
			ServerName:    "user",
			ServerVersion: "0.1",
			PayloadLen:    355,
			DomainId:      1,
			AppId:         2,
			Expend:        "cccc",
		},
	}

	testMsgByte := transMsgToByte(&msg)
	header := transByteToHeader(testMsgByte[0:400])
	println(header)
}
