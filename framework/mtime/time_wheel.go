package mtime

import (
	"context"
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/antlabs/stl/list"
)

const (
	nearShift = 8

	nearSize = 1 << nearShift

	levelShift = 6

	levelSize = 1 << levelShift

	nearMask = nearSize - 1

	levelMask = levelSize - 1
)

type TimeWheel struct {
	// 单调递增累加值, 走过一个时间片就+1
	jiffies uint64

	// 256个槽位
	t1 [nearSize]*Time

	// 4个64槽位, 代表不同的刻度
	t2Tot5 [4][levelSize]*Time

	// 时间只精确到10ms
	// curTimePoint 为1就是10ms 为2就是20ms
	curTimePoint time.Duration

	// 上下文
	ctx context.Context

	// 取消函数
	cancel context.CancelFunc
}

func NewTimeWheel() *TimeWheel {
	ctx, cancel := context.WithCancel(context.Background())
	t := &TimeWheel{ctx: ctx, cancel: cancel}
	t.init()
	return t
}

func (t *TimeWheel) init() {

	for i := 0; i < nearSize; i++ {
		t.t1[i] = newTimeHead(1, uint64(i))

	}

	for i := 0; i < 4; i++ {
		for j := 0; j < levelSize; j++ {
			t.t2Tot5[i][j] = newTimeHead(uint64(i+2), uint64(j))
		}
	}

	t.curTimePoint = get10Ms()
}

func maxVal() uint64 {
	return (1 << (nearShift + 4*levelShift)) - 1
}

func levelMax(index int) uint64 {
	return 1 << (nearShift + index*levelShift)
}

func (t *TimeWheel) index(n int) uint64 {
	return (t.jiffies >> (nearShift + levelShift*n)) & levelMask
}

func (t *TimeWheel) add(node *TimeNode, jiffies uint64) *TimeNode {
	var head *Time
	expire := node.expire
	idx := expire - jiffies

	level, index := uint64(1), uint64(0)

	if idx < nearSize {
		index = uint64(expire) & nearMask
		head = t.t1[index]

	} else {
		max := maxVal()
		for i := 0; i <= 3; i++ {
			if idx > max {
				idx = max
				expire = idx + jiffies
			}
			if idx < levelMax(i+1) {
				index = uint64(expire >> (nearShift + i*levelShift) & levelMask)
				head = t.t2Tot5[i][index]
				level = uint64(i) + 2
				break
			}
		}
	}

	if head == nil {
		logger.Error("TimeWheel head error")
	}

	head.lockPushBack(node, level, index)
	return node
}

func (t *TimeWheel) AfterFunc(expire time.Duration, callback func(), async bool) *TimeNode {

	jiffies := atomic.LoadUint64(&t.jiffies)

	expire = expire/(time.Millisecond*10) + time.Duration(jiffies)

	node := &TimeNode{
		expire:   uint64(expire),
		callback: callback,
		async:    async,
	}

	return t.add(node, jiffies)
}

func getExpire(expire time.Duration, jiffies uint64) time.Duration {
	return expire/(time.Millisecond*10) + time.Duration(jiffies)
}

func (t *TimeWheel) ScheduleFunc(userExpire time.Duration, callback func(), async bool) *TimeNode {

	jiffies := atomic.LoadUint64(&t.jiffies)

	expire := getExpire(userExpire, jiffies)

	node := &TimeNode{
		userExpire: userExpire,
		expire:     uint64(expire),
		callback:   callback,
		isSchedule: true,
		async:      async,
	}

	return t.add(node, jiffies)
}

func (t *TimeWheel) Stop() {
	t.cancel()
}

// 移动链表
func (t *TimeWheel) cascade(levelIndex int, index int) {

	tmp := newTimeHead(0, 0)

	l := t.t2Tot5[levelIndex][index]
	l.Lock()
	if l.Len() == 0 {
		l.Unlock()
		return
	}

	l.ReplaceInit(&tmp.Head)

	// 每次链表的元素被移动走，都修改version
	atomic.AddUint64(&l.version, 1)
	l.Unlock()

	offset := unsafe.Offsetof(tmp.Head)
	tmp.ForEachSafe(func(pos *list.Head) {
		node := (*TimeNode)(pos.Entry(offset))
		t.add(node, atomic.LoadUint64(&t.jiffies))
	})

}

// moveAndExec函数功能
//1. 先移动到near链表里面
//2. near链表节点为空时，从上一层里面移动一些节点到下一层
//3. 再执行
func (t *TimeWheel) moveAndExec() {

	// 这里时间溢出
	if uint32(t.jiffies) == 0 {
		// TODO
		// return
	}

	//如果本层的盘子没有定时器，这时候从上层的盘子移动一些过来
	index := t.jiffies & nearMask
	if index == 0 {
		for i := 0; i <= 3; i++ {
			index2 := t.index(i)
			t.cascade(i, int(index2))
			if index2 != 0 {
				break
			}
		}
	}

	atomic.AddUint64(&t.jiffies, 1)

	t.t1[index].Lock()
	if t.t1[index].Len() == 0 {
		t.t1[index].Unlock()
		return
	}

	head := newTimeHead(0, 0)
	t1 := t.t1[index]
	t1.ReplaceInit(&head.Head)
	atomic.AddUint64(&t1.version, 1)
	t.t1[index].Unlock()

	// 执行,链表中的定时器
	offset := unsafe.Offsetof(head.Head)

	head.ForEachSafe(func(pos *list.Head) {
		val := (*TimeNode)(pos.Entry(offset))
		head.Del(pos)

		if atomic.LoadUint32(&val.stop) == haveStop {
			return
		}

		if val.async {
			go val.callback()
		} else {
			val.callback()
		}

		if val.isSchedule {
			jiffies := t.jiffies
			// 这里的jiffies必须要减去1
			// 当前的callback被调用，已经包含一个时间片,如果不把这个时间片减去，
			// 每次多一个时间片，就变成累加器, 最后周期定时器慢慢会变得不准
			val.expire = uint64(getExpire(val.userExpire, jiffies-1))
			t.add(val, jiffies)
		}
	})

}

// get10Ms函数通过参数传递，为了方便测试
func (t *TimeWheel) run(get10Ms func() time.Duration) {
	// 先判断是否需要更新
	// 内核里面实现使用了全局jiffies和本地的jiffies比较,应用层没有jiffies，直接使用时间比较
	// 这也是skynet里面的做法

	ms10 := get10Ms()

	if ms10 < t.curTimePoint {

		fmt.Printf("github.com/antlabs/timer:Time has been called back?from(%d)(%d)\n",
			ms10, t.curTimePoint)

		t.curTimePoint = ms10
		return
	}

	diff := ms10 - t.curTimePoint
	t.curTimePoint = ms10

	for i := 0; i < int(diff); i++ {
		t.moveAndExec()
	}

}

func (t *TimeWheel) Run() {

	// 10ms精度
	tk := time.NewTicker(time.Millisecond * 10)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			t.run(get10Ms)
		case <-t.ctx.Done():
			return
		}
	}
}
