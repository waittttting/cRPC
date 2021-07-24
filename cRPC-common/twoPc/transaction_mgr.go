package twoPc

import "errors"

type ResourceManage interface {
	Prepare() bool
	Commit() bool
}

type TransactionManager struct {
	rm []*ResourceManage
}

func newTransactionMgr(rms ...*ResourceManage) *TransactionManager {
	rm := TransactionManager{
		rm: rms,
	}
	return &rm
}

func (tm *TransactionManager) Do() error {
	if len(tm.rm) <= 0 {
		return errors.New("rm len <= 0")
	}

	return nil
}

