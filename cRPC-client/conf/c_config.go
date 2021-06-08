package conf

type LocalConf struct {
	// 本服务作为 Client 的一些配置
	Client Client
	// 本服务作为 Server 的一些配置
	Server Server
}

type Client struct {
	ServerName 			string
	ServerVersion 		string
	ConfigCenterHost 	string
}

type Server struct {
	TcpPort int
	HttpPort int
	// 是否需要开启本地 RPC 接口
	RpcOpen bool
	// 接收 socket 队列的长度
	ReceiveSocketChanLen int
	ReceiveTcpMsgChanLen int
}