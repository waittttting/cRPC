package model

/**
 * @Description: 从云端（配置中心）获取的配置信息
 */
type CloudConfigInfo struct {
	// 本服务在配置中心的配置
	ServerConfig ServerConfig `json:"server_config"`
	// 订阅的服务的信息
	SubServersInfos []*SubServerInfos `json:"sub_servers_infos"`
	// 控制中心 URL
	ControlCenterAddr ControlCenterAddr `json:"control_center_addr"`
}

type ServerConfig struct {
	// 服务名
	ServerName string `json:"server_name"`
	// 服务版本号
	ServerVersion string `json:"server_version"`
	// 订阅的服务的名字的列表
	SubServers string `json:"sub_servers"`
}

type ControlCenterAddr struct {
	Host     string `json:"host"`
	TcpPort  string `json:"tcp_port"`
	HttpPort string `json:"http_port"`
}

type SubServerInfos struct {
	ServerName string   `json:"server_name"`
	Infos      []string `json:"infos"` // gid 的 json 数组
}

type PortConfig struct {
	Port int `json:"port"`
}
