package cweb

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/cerr"
	http2 "github.com/waittttting/cRPC-common/http"
	"github.com/waittttting/cRPC-control-center/server"
	"net/http"
)

type WebHandler struct {}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

type ServerNameList struct {
	SubServerName []string `json:"sub_server_name"`
}

func (wh *WebHandler) getServersIpList(c *gin.Context) {

	snls := c.PostForm("serverNameList")
	var snl ServerNameList
	err := json.Unmarshal([]byte(snls), &snl)
	if err != nil {
		c.JSON(http.StatusOK, cerr.ErrInternal)
	}
	rets := make(map[string][]string, 0)
	// 根据服务名，在 redis 里获取该服务的所有节点列表
	for _, serverName := range snl.SubServerName {
		ret, err := server.RedisCli.SMembers(serverName).Result()
		if err != nil {
			c.JSON(http.StatusOK, cerr.ErrInternal)
			continue
		}
		rets[serverName] = ret
	}
	logrus.Infof("------ response of getServersIpList :[%v]", rets)
	c.JSON(http.StatusOK, http2.NewResponseWithData(rets))
}


