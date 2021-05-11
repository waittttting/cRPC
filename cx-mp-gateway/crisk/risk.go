package crisk

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

var GRS RiskServer

type RiskServer struct {}

// 风控 check 维度 ip uid 等，数据量小占用内存小的时候的时候可以采用本地内存缓存名单
// 数据量大的时候，采用布隆过滤器的方案

func init() {
	GRS = RiskServer{}
}

func (rs *RiskServer) Start() {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("updateRsData occurred err:[%v]", err)
			}
		}()
		ticker := time.Tick(1 * time.Second)
		for  {
			rs.updateRsData()
			<- ticker
		}
	}()
}

func (rs *RiskServer) updateRsData() {
	// todo: 更新 Redis 中的黑名单数据到本地缓存
}

func (rs *RiskServer) CheckRisk(c *gin.Context) bool {
	// todo: check 逻辑
	return true
}


