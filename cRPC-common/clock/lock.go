package clock

import "time"

type lockType int

const (
	lockTypeEtcdCP  = 1
	lockTypeRedisAP = 2
	defaultLease = 10
	lockDefaultValue = "etcd_default_value"
)

type Lock interface {
	Lock(key string, value string, lease int64) (bool, error)
	UnLock() (bool, error)
	Close() error
}


type CXLock struct {
	lock Lock
	lt lockType
	locked bool
	lease int64
	singleChan chan int
}

/**
 * @Description: 设置续约时间 默认为 10
 * @receiver cl
 * @param lease 续约时间
 */
func (cl *CXLock) SetLease(lease int64) {
	cl.lease = lease
}


/**
 * @Description: 基于 etcd 的 CP 分布式锁
 * @param host etcd 集群地址
 * @return *CXLock
 * @return error
 */
func EtcdLock(host []string) (*CXLock, error) {
	if len(host) <= 0 {
		return nil, errHostLen
	}
	cl := &CXLock{
		locked: false,
		singleChan: make(chan int, 1),
	}
	cl.lt = lockTypeEtcdCP
	lock, err := newEtcdLock(host)
	if err != nil {
		return nil, err
	}
	cl.lock = lock
	return cl, nil
}

/**
 * @Description: 基于 redis 的 AP 分布式锁
 * @param host redis 地址
 * @param pwd 密码
 * @param dbIndex 索引
 * @return *CXLock
 * @return error
 */
func RedisLock(host string, pwd string, dbIndex int) (*CXLock, error) {
	if len(host) <= 0 {
		return nil, errHostLen
	}
	cl := &CXLock{
		locked: false,
		singleChan: make(chan int, 1),
	}
	cl.lt = lockTypeRedisAP
	lock, err := newRedisLock(host, pwd, dbIndex, cl.singleChan)
	if err != nil {
		return nil, err
	}
	cl.lock = lock
	return cl, nil
}

/**
 * @Description: 获取锁不阻塞
 * @receiver lc
 * @param ttl 锁过期时间
 * @return bool 是否获取锁
 * @return error
 */
func (cl *CXLock) Acquire(key string) (bool, error) {
	if key == "" {
		return false, errKeyLen
	}
	if cl.locked {
		return false, errSameInstance
	}
	if cl.lock == nil {
		return false, errLockInit
	}
	if cl.lease <= 0 {
		cl.lease = defaultLease
	}
	ok, err := cl.lock.Lock(key, lockDefaultValue, cl.lease)
	if ok {
		cl.locked = true
		return true, nil
	}
	return false, err
}

/**
 * @Description: 获取锁，直到重试超时
 * @receiver lc LockClient
 * @param ttl 锁过期时间
 * @param interval 获取锁时间间隔
 * @param maxRetry 获取锁重试次数
 * @return bool 是否获取锁
 * @return error
 */
func (cl *CXLock) AcquireWithRetry(key string, interval int, maxRetry int) (bool, error) {

	if cl.locked {
		return false, errSameInstance
	}
	if cl.lock == nil {
		return false, errLockInit
	}
	if interval < 0  {
		interval = 0
	}
	if maxRetry < 0 {
		maxRetry = 0
	}
	if cl.lease <= 0 {
		cl.lease = defaultLease
	}
	for i := 0; i < maxRetry; i ++ {
		ok, err := cl.lock.Lock(key, lockDefaultValue, cl.lease)
		if err != nil && err.Error() != errKeyExist.Error() {
			return false, err
		}
		if ok {
			cl.locked = true
			return true, nil
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return false, nil
}

/**
 * @Description: 释放锁
 * @receiver lc
 */
func (cl *CXLock) Release() (bool, error) {
	if !cl.locked {
		return false, errReleaseUnLockedKey
	}
	ok, err := cl.lock.UnLock()
	if ok {
		cl.locked = false
		return ok, nil
	}
	return false, err
}

/**
 * @Description: 释放与锁相关的资源
 * @receiver cl
 * @return error
 */
func (cl *CXLock) Close() error {
	return cl.lock.Close()
}