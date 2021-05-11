package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/cerr"
	http2 "github.com/waittttting/cRPC-common/http"
	"github.com/waittttting/cRPC-common/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type WebHandler struct {
	db *gorm.DB
	redisCli *redis.Client
}

func NewWebHandler(db *gorm.DB, client *redis.Client) *WebHandler {

	return &WebHandler{
		db: db,
		redisCli: client,
	}
}

/**
 * @Description: 获取服务配置
 * @receiver wh
 * @param c
 */

func (wh *WebHandler) GetConfig(c *gin.Context) {

	serverName := c.PostForm("server_name")
	serverVersion := c.PostForm("server_version")

	sc := model.ServerConfig{
		ServerName: serverName,
		ServerVersion: serverVersion,
	}

	// 获取 control center 的 url
	ccAddr, err := wh.getControlCenterUrl()
	if err != nil {
		logrus.Errorf("get ControlCenterUrl err:[%v]", err)
		c.JSON(http.StatusOK, http2.NewResponseWithErr(cerr.ErrDB))
		return
	}

	// 在数据库中查询配置信息
	result := wh.db.Table("server_config").Where(&sc).Find(&sc)
	if result.Error != nil {
		logrus.Errorf("query server config err:[%v]", result.Error)
		c.JSON(http.StatusOK, http2.NewResponseWithErr(cerr.ErrDB))
		return
	}
	// 查询订阅的服务的IP列表
	ssi, err := getServersIPList(ccAddr, sc.SubServers)
	if err != nil {
		logrus.Errorf("query sub infos err:[%v]", err)
		c.JSON(http.StatusOK, http2.NewResponseWithErr(cerr.ErrInternal))
		return
	}
	sci := &model.CloudConfigInfo{
		ServerConfig: sc,
		SubServersInfos: ssi,
		ControlCenterAddr: *ccAddr,
	}
	c.JSON(http.StatusOK, http2.NewResponseWithData(sci))
}

type GlobalConfig struct {
	Key string
	Value string
}

/**
 * @Description: 获取 control center 的 url
 * @receiver wh
 * @return string 地址
 */
func (wh *WebHandler) getControlCenterUrl() (*model.ControlCenterAddr, error) {

	var cca model.ControlCenterAddr

	// 查缓存
	ccas, err := wh.redisCli.Get(configCenterRedisKeyControlCenterUrl).Result()
	if err != nil && err.Error() != cerr.ErrRedisNil.Error() {
		return nil, err
	}
	if ccas == "" {
		gc := GlobalConfig{
			Key: "control_center_url",
		}
		result := wh.db.Table("global_config").Where(&gc).Find(&gc)
		if result.Error != nil {
			return nil, result.Error
		}
		ccas = gc.Value
	}

	// 查库
	err = json.Unmarshal([]byte(ccas), &cca)
	if err != nil {
		return nil, err
	}
	// 写入缓存
	ok, err := wh.redisCli.SetNX(configCenterRedisKeyControlCenterUrl, ccas, 30 * time.Second).Result()
	if !ok || err != nil {
		logrus.Errorf("set control_center_url to redis err:[%v]", err)
	}
	return &cca, nil
}

/**
 * @Description: 获取订阅的服务的信息
 * @param cca
 * @param serverNameList
 * @return []*model.SubServerInfos
 * @return error
 */
func getServersIPList(cca *model.ControlCenterAddr, serverNameList string) ([]*model.SubServerInfos, error) {

	params := map[string]string{
		"serverNameList" : serverNameList,
	}
	address := "http://" + "cx-control-center-svc" + ":" + cca.HttpPort + "/get/serverIpLists"
	result, err := http2.Post(address, params)
	if err != nil {
		return nil, err
	}
	// rets 是一个map，key = 订阅的服务，value = 每个服务的节点列表
	rets := result.(map[string]interface{})
	subServerInfos := make([]*model.SubServerInfos, 0)
	for k, v := range rets {
		truthV := v.([]interface{})
		nodes := make([]string, 0)
		for _, ts := range truthV {
			nodes = append(nodes, ts.(string))
		}
		ssi := &model.SubServerInfos{
			ServerName: k,
			Infos: nodes,
		}
		subServerInfos = append(subServerInfos, ssi)
	}
	if err != nil {
		return nil, err
	}
	// TODO: 是否需要加缓存，缓存与实际数据不一致问题~~
	return subServerInfos, nil
}
