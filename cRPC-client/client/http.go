package client

import (
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/http"
	"github.com/waittttting/cRPC-common/model"
)

/**
 * @Description: 获取保存在云端（配置中心）的配置
 * @receiver rc
 * @return *model.cloudConfig
 */
func (rc *RpcClient) getServerConfig() (*model.CloudConfigInfo, error) {

	params := map[string]string{
		"server_name" : rc.localConfig.Client.ServerName,
		"server_version" : rc.localConfig.Client.ServerVersion}
	result, err := http.Post(rc.localConfig.Client.ConfigCenterHost + "/get/config", params)
	if err != nil {
		return nil, err
	}
	ret := result.(map[string]interface{})
	var cci model.CloudConfigInfo
	serverConfigMap := ret["server_config"].(map[string]interface{})
	cci.ServerConfig.ServerName = serverConfigMap["server_name"].(string)
	cci.ServerConfig.ServerVersion = serverConfigMap["server_version"].(string)
	cci.ServerConfig.SubServers = serverConfigMap["sub_servers"].(string)

	temp := ret["sub_servers_infos"].([]interface{})
	cci.SubServersInfos = make([]*model.SubServerInfos, 0)
	for _, v := range temp {
		truthV := v.(map[string]interface{})
		truthKey := truthV["server_name"].(string)
		truthValue := truthV["infos"].([]interface{})
		infos := make([]string, 0)
		for _, ts := range truthValue {
			infos = append(infos, ts.(string))
		}
		ssi := model.SubServerInfos{
			ServerName: truthKey,
			Infos: infos,
		}
		cci.SubServersInfos = append(cci.SubServersInfos, &ssi)
	}

	controlCenterAddrMap := ret["control_center_addr"].(map[string]interface{})
	cci.ControlCenterAddr.Host = controlCenterAddrMap["host"].(string)
	cci.ControlCenterAddr.TcpPort = controlCenterAddrMap["tcp_port"].(string)
	cci.ControlCenterAddr.HttpPort = controlCenterAddrMap["http_port"].(string)

	if err != nil {
		logrus.Fatalf("load server config err [%v]", err)
		return nil, err
	}
	return &cci, nil
}
