package conf

import "time"

type CCSConf struct {
	Server Server
	MySQL MySQL
	Redis Redis
}

type Server struct {
	TcpPort int
	HttpPort int
	ReceiveSocketChanLen int
	ReceiveSocketTimeoutMs time.Duration
}

type MySQL struct {
	Host string
	User string
	Password string
}

type Redis struct {
	Host  string
	Pwd   string
	Index int
}