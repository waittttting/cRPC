package clock


type LockError struct {
	errString string
}

func (le *LockError) Error() string {
	return le.errString
}

func newLockError(err string) error {
	return &LockError{
		errString: err,
	}
}

var (
	errKeyLen     				= newLockError("key len is 0")
	errKeyExist   				= newLockError("key already exist")
	errSameInstance				= newLockError("try to get lock use same instance")
 	errLockInit   				= newLockError("lock init err, lock can not be nil")
 	errHostLen   				= newLockError("host len can not be 0")
 	errRedisNil					= newLockError("redis: nil")
	errRedisDel					= newLockError("redis del int ret not be 1")
	errReleaseUnLockedKey		= newLockError("the lock of the key not locked")
)