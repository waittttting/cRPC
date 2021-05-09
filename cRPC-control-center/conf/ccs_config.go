package conf

type CCSConf struct {
	Server    Server
	Redis     Redis
	TimeWheel TimeWheel
}

type Server struct {

	HttpPort int
	TcpPort int
	// 接收 socket 队列的长度
	ReceiveSocketChanLen int
	// 接收 serviceConn 队列的长度
	ReceiveServiceConnChanLen int
}

type Redis struct {
	Host  string
	Pwd   string
	Index int
}

type TimeWheel struct {
	Cap           int
	NoticeChanLen int
}
