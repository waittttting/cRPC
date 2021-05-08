package clock

import (
	"github.com/go-redis/redis"
	"time"
)

const redisUnLockSingle = 0

type redisLock struct {
	client *redis.Client
	singleChan chan int
	key string
}

func newRedisLock(host string, pwd string, dbIndex int, singleChan chan int) (*redisLock, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       dbIndex,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &redisLock{
		client: client,
		singleChan: singleChan,
	}, nil
}

func (rl *redisLock) Lock(key string, value string, lease int64) (bool, error) {

	ret, err := rl.client.SetNX(key, value, time.Duration(lease) * time.Second).Result()
	if err != nil {
		return false, err
	}
	if !ret {
		return false, errKeyExist
	}
	if !ret {
		return false, nil
	}
	rl.key = key
	// 续租协程
	go func(singleChan chan int) {
		for {
			select {
			case single := <-singleChan:
				if single == redisUnLockSingle {
					break
				}
			default:
				time.Sleep(time.Duration(lease / 3 * 2) * time.Second)
				_, _ = rl.client.Set(key, value, time.Duration(lease) * time.Second).Result()
				// todo: 续租错误通过 chan 返回给应用层~~
			}
		}
	}(rl.singleChan)
	return true, nil
}

func (rl *redisLock) UnLock() (bool, error) {
	// 停止续租
	rl.singleChan <- redisUnLockSingle
	_, err := rl.client.Del(rl.key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rl *redisLock) Close() error {
	return rl.client.Close()
}