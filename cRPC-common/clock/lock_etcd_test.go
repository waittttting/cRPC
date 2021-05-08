package clock

import (
	"testing"
	"time"
)

var (
	testKey1 = "/testLock/key1"
	testKey2 = "/testLock/key2"
	etcdHosts = []string{"127.0.0.1:2379"}
)

func TestInitEtcdLock(t *testing.T) {

	// todo: etcd host 填写错误时，etcd 调用方法无超时，无返回
	lock1, err := EtcdLock(etcdHosts)
	if err != nil {
		t.Fatalf("create lock1 error: %v", err)
	}
	lock1.SetLease(5)
	ok, err := lock1.Acquire(testKey1)
	if err != nil || !ok {
		t.Fatalf("lock1 acquire error: %v", err)
	}
}

func TestEtcdLockAcquire(t *testing.T) {

	lock1, err := EtcdLock(etcdHosts)
	if err != nil {
		t.Fatalf("create lock1 error: %v", err)
	}
	lock1.SetLease(5)
	ok, err := lock1.Acquire(testKey1)
	if err != nil || !ok {
		t.Fatalf("lock1 acquire error: %v", err)
	}
	t.Logf("lock1 locked")
	go func() {
		lock2, err := EtcdLock(etcdHosts)
		if err != nil {
			t.Errorf("create lock2 error: %v", err)
		}
		i := 0
		for {
			ok, err := lock2.Acquire(testKey1)
			if err != nil {
				if err.Error() == errKeyExist.Error()  {
					t.Logf("lock2 try to get lock, loop count: %v", i)
				} else {
					t.Fatalf("lock2 acquire error: %v",  err)
				}
			}
			if ok {
				t.Logf("lock2 get locked, loop count: %v", i)
				ok, err = lock2.Release()
				if err != nil {
					t.Errorf("lock2 release err: %v", err)
				}
				if ok {
					t.Logf("lock2 release")
				}
				err = lock2.Close()
				if err != nil {
					t.Errorf("lock2 close err: %v", err)
				}
				break
			}
			time.Sleep(1 * time.Second)
			i ++
		}
	}()
	time.Sleep(10 * time.Second)
	ok, err = lock1.Release()
	if err != nil {
		t.Errorf("lock1 release err: %v", err)
	}
	if ok {
		t.Logf("lock1 release")
	}
	err = lock1.Close()
	if err != nil {
		t.Errorf("lock1 close err: %v", err)
	}
	time.Sleep(1 * time.Second)
}

func TestLockUseSameInstance(t *testing.T) {

	lock1, err := EtcdLock(etcdHosts)
	if err != nil {
		t.Fatalf("create lock1 error: %v", err)
	}
	lock1.SetLease(5)
	ok, err := lock1.Acquire(testKey1)
	if err != nil || !ok {
		t.Fatalf("lock1 acquire error: %v", err)
	}
	t.Logf("lock1 locked")

	ok, err = lock1.Acquire(testKey2)
	if err != nil && err.Error() == errSameInstance.Error() {
		t.Logf("lock2 acquire error: %v", err)
	}

	ok, err = lock1.Release()
	if err != nil {
		t.Fatalf("lock1 release error: %v", err)
	}
	if ok {
		t.Logf("lock1 released")
	}

	ok, err = lock1.Acquire(testKey2)
	if err != nil {
		t.Fatalf("lock2 acquire error: %v", err)
	}
	if ok {
		t.Logf("lock1 locked")
	}
	lock1.Release()
	lock1.Close()
}

func TestEtcdLockAcquireWithRetry(t *testing.T) {
	ok := testEtcdLockAcquireWithRetry(t,1, 3)
	if ok {
		t.Error("get lock error")
	} else {
		t.Logf("get lock fail")
	}
	ok = testEtcdLockAcquireWithRetry(t,2, 13)
	if !ok {
		t.Error("get lock error")
	} else {
		t.Logf("get lock sucess")
	}
}

func testEtcdLockAcquireWithRetry(t *testing.T, interval int, maxRetry int) bool {
	lock1, err := EtcdLock(etcdHosts)
	if err != nil {
		t.Fatalf("create lock1 error: %v", err)
	}
	lock1.SetLease(18)
	ok, err := lock1.Acquire(testKey1)
	if err != nil || !ok {
		t.Fatalf("lock1 acquire error: %v", err)
	}
	t.Logf("lock1 locked")
	go func() {
		lock2, err := EtcdLock(etcdHosts)
		if err != nil {
			t.Fatalf("create lock2 error: %v", err)
		}
		t.Log("lock2 try to get lock")
		ok, err = lock2.AcquireWithRetry(testKey1, interval, maxRetry)
		if err != nil {
			t.Fatalf("lock2 acquireWithRetry error: %v", err)
		}
		if ok {
			t.Log("lock2 get lock")
		} else {
			t.Log("lock2 not get lock")
		}
	}()
	time.Sleep(20 * time.Second)
	t.Log("lock1 try to release")
	lock1.Release()
	t.Log("lock1 released")
	lock1.Close()
	time.Sleep(5 * time.Second)
	return ok
}

