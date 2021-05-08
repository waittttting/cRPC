package tcp

import (
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func AcceptSocket(port string, socketChan chan *net.TCPConn, timeoutMs time.Duration) {
	// 获取监听地址
	addr, err := net.ResolveTCPAddr("tcp4", ":" + port)
	if err != nil {
		logrus.Fatalf("resolve TCP addr error [%v]", err)
	}
	// 获取监听器
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logrus.Fatalf("listen TCP error [%v]", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("access tcp_server occur panic: %v", r)
			}
		}()
		for {
			// 接收长连接
			socket, err := ln.AcceptTCP()
			if err != nil {
				logrus.Errorf("accept TCP error [%v]", err)
			}
			timer := time.NewTimer(timeoutMs * time.Millisecond)
			select {
			case socketChan <- socket:
				timer.Reset(0)
				logrus.Infof("socket accepted [%v]", socket.RemoteAddr())
			case <-timer.C:
				timer.Reset(0)
				socket.Close()
				logrus.Errorf("receive socket timeout")
			}
		}
	}()
}