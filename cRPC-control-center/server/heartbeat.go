package server

import "github.com/sirupsen/logrus"

func (ccs *ControlCenterServer) heartbeatTimeWheelStart() {
	ccs.timeWheel.Start()
	logrus.Info("server's timeWheel started ")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("heartbeatTimeoutChan err %v", err)
			}
		}()
		for inter := range ccs.heartbeatTimeoutChan {
			sc := inter.(*serviceConn)
			logrus.Infof("heartbeat timeout %s", sc.gid.String())
			// todo: 接收消息超时的逻辑是否写全
			sc.LetScOffLine("heartbeat timeout")
			ccs.timeWheel.Delete(sc)
		}
	}()
}

func (ccs *ControlCenterServer) addToHeartbeatTimeWheel(client *serviceConn) {
	ccs.timeWheel.Add(client, 5)
}

func (ccs *ControlCenterServer) refreshHeartbeat(client *serviceConn) {
	logrus.Infof("received heartbeat:[%v]", client.gid.String())
	ccs.timeWheel.Refresh(client, 5)
}

func (ccs *ControlCenterServer) deleteHeartbeat(client *serviceConn) {
	ccs.timeWheel.Delete(client)
}
