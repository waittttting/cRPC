package cweb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/waittttting/cRPC-control-center/conf"
)

type WebServer struct {}

func NewWebServer() *WebServer {
	return &WebServer{

	}
}

func (ws *WebServer) Start(config conf.CCSConf)  {

	// todo: 开启特定个数的 goroutine 处理 web 消息

	wh := NewWebHandler()
	r := gin.Default()

	get := r.Group("/get")
	get.POST("/serverIpLists", wh.getServersIpList)
	r.Run(fmt.Sprintf(":%v", config.Server.HttpPort))
}
