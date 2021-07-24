package server

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/tcp"
	"github.com/waittttting/cRPC-control-center/conf"
)


var RedisCli *redis.Client

type redisOp int

const (
	redisOpSAddServerIp = 1
	redisOpSRemServerIp = 2
)

func RedisInit(conf conf.CCSConf) {

	RedisCli = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Pwd,
		DB:       conf.Redis.Index,
	})
}

func RedisOp(key, value string, op redisOp) error {

	// redis 单线程（线程安全）
	var err error
	switch op {
	case redisOpSAddServerIp:
		logrus.Infof("register info: key:[%v], ip:[%v]", key, value)
		_, err = RedisCli.SAdd(key, value).Result()
	case redisOpSRemServerIp:
		logrus.Infof("delete info: key:[%v], ip:[%v]", key, value)
		_, err = RedisCli.SRem(key, value).Result()
	}
	if err != nil {
		return err
	}
	logrus.Infof("redis client set, key:[%v], value:[%v]", key, RedisCli.SMembers(key))
	return nil
}

const (
	serverOnLinAndOffLine = "serverOnLinAndOffLine"
	serviceOnLine = "1"
	serviceOffLine = "2"
)

func RedisSubServerOnLine() (<-chan *redis.Message, error) {

	pubSub := RedisCli.Subscribe(serverOnLinAndOffLine)
	_, err := pubSub.Receive()
	if err != nil {
		return nil, err
	}
	ch := pubSub.Channel()
	return ch, nil
}

type messageOfRedisMQ struct {
	gid tcp.GID
	onLineState string
}

func redisPubOnLine(gid *tcp.GID) {
	mor := &messageOfRedisMQ{
		gid: *gid,
		onLineState: serviceOnLine,
	}
	RedisCli.Publish(serverOnLinAndOffLine, mor)
}

func redisPubOffline(gid *tcp.GID) {
	mor := &messageOfRedisMQ{
		gid: *gid,
		onLineState: serviceOffLine,
	}
	RedisCli.Publish(serverOnLinAndOffLine, mor)
}