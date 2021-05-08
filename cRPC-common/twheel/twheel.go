package twheel

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type TimeWheel struct {
	slots      []map[interface{}]interface{}
	posMap     map[interface{}]int
	curPos     int
	slotsCount int
	lock       sync.Mutex
	noticeChan chan interface{}
}

/**
 * @Description: 创建一个新的时间轮
 * @param cap 时间轮中循环数组的大小
 * @param noticeChanLen 元素超时回调队列
 * @return *TimeWheel
 */
func New(slotsCount int, noticeChan chan interface{}) (*TimeWheel, error) {

	if slotsCount <= 0 {
		return nil, errors.New("cap can not <= 0")
	}

	return &TimeWheel{
		// 存储数据的槽
		slots: make([]map[interface{}]interface{}, slotsCount),
		// 槽数量
		slotsCount: slotsCount,
		// tw 当前指针
		curPos: 0,
		// 到期回调 chan
		noticeChan: noticeChan,
		// 元素位置 map
		posMap: make(map[interface{}]int),
	}, nil
}

func (tw *TimeWheel) Start() {
	go func() {
		ticker := time.Tick(1 * time.Second)
		for {
			<-ticker
			tw.lock.Lock()
			tw.curPos = (tw.curPos + 1) % tw.slotsCount
			curBucket := tw.slots[tw.curPos]
			if len(curBucket) == 0 {
				tw.lock.Unlock()
				continue
			}
			for _, value := range curBucket {
				timer := time.NewTimer(500 * time.Millisecond)
				// 1. 通知
				select {
				case tw.noticeChan <- value:
					timer.Reset(0)
				case <-timer.C:
					logrus.Errorf("item send to noticeChan timeout")
				}
				// 2. 删除 index map 中的数据
				delete(tw.posMap, value)
			}
			// 3. 清空 curIndex 上的 map
			tw.slots[tw.curPos] = make(map[interface{}]interface{})
			tw.lock.Unlock()
		}
	}()
}

func (tw *TimeWheel) Add(item interface{}, timeout int) error {
	if timeout <= 0 {
		return errors.New("timeout can not <= 0")
	}
	tw.lock.Lock()
	defer tw.lock.Unlock()
	itemIndex := (tw.curPos + timeout) % tw.slotsCount
	tw.posMap[item] = itemIndex
	if tw.slots[itemIndex] == nil {
		tw.slots[itemIndex] = map[interface{}]interface{}{}
	}
	tw.slots[itemIndex][item] = item
	return nil
}

func (tw *TimeWheel) Delete(item interface{}) {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	oldPos := tw.posMap[item]
	// todo: if oldPos == nil
	delete(tw.slots[oldPos], item)
	delete(tw.posMap, item)
}

func (tw *TimeWheel) Refresh(item interface{}, timeout int) error {

	if timeout <= 0 {
		return errors.New("timeout can not <= 0")
	}

	tw.lock.Lock()
	defer tw.lock.Unlock()
	oldPos := 0
	ok := false
	if oldPos, ok = tw.posMap[item]; !ok {
		return errors.New(" item not find in time wheel slots ")
	}
	delete(tw.slots[oldPos], item)
	delete(tw.posMap, item)
	newPos := (tw.curPos + timeout) % tw.slotsCount
	tw.posMap[item] = newPos
	if tw.slots[newPos] == nil {
		tw.slots[newPos] = map[interface{}]interface{}{}
	}
	tw.slots[newPos][item] = item
	return nil
}
