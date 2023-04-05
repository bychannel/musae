package mtime

import (
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"sync/atomic"
	"time"
)

var (
	GTimeWheel   *TimeWheel
	startTime    time.Time // 进程启动时间
	monotonicSec int64     // 单调时间, 进程启动时间时长, 精度秒
	timeOffset   int64     // 系统时间偏移量, 用于debug模式动态调整系统时间
)

func init() {
	timeOffset = 0
	now := time.Now()
	startTime = now.Add(-5 * time.Second)
	atomic.StoreInt64(&monotonicSec, int64(time.Since(startTime)/time.Second))
	GTimeWheel = NewTimeWheel()

	go func() {
		t := time.NewTicker(time.Second * 1)
		defer t.Stop()
		for range t.C {
			deltaTime := int64(time.Since(startTime) / time.Second)
			atomic.StoreInt64(&monotonicSec, deltaTime)
		}
	}()

	go GTimeWheel.Run()
}

//
// Now
//  @Description: 系统时间,包含偏移量
//  @return time.Time
//
func Now() time.Time {
	offset := time.Duration(atomic.LoadInt64(&timeOffset))
	return time.Now().Add(offset * time.Second)
}

//
// RealNow 不包含时间偏移的当前时间
//  @Description: 系统时间, 不包含偏移量
//  @return time.Time
//
func RealNow() time.Time {
	return time.Now()
}

//
// SetTimeOffset
//  @Description: 设置时间偏移量
//  @param offset
//
func SetTimeOffset(offset int64) {
	curOffset := atomic.LoadInt64(&timeOffset)
	if curOffset > offset {
		logger.Warn("[mtime] offset can't back", offset, curOffset)
		return
	}

	atomic.StoreInt64(&timeOffset, offset)
}

//
// Add
//  @Description: 添加定时器
//  @param delay
//  @param callback
//  @param async 异步时注意data race
//  @return *TimeNode
//
func Add(delay time.Duration, callback func(), async bool) *TimeNode {
	return GTimeWheel.AfterFunc(delay, callback, async)
}

//
// AddCron
//  @Description: 添加周期循环定时器
//  @param delay
//  @param callback
//  @param async 异步时注意data race
//  @return *TimeNode
//
func AddCron(delay time.Duration, callback func(), async bool) *TimeNode {
	return GTimeWheel.ScheduleFunc(delay, callback, async)
}

//
// Remove
//  @Description: 移除定时器
//  @param task
//
func Remove(task *TimeNode) {
	task.Stop()
}

func AfterFunc(delay time.Duration, callback func()) *Timer {
	timer := &Timer{}
	timer.SetAfterFunc(delay, callback)
	return timer
}

func After(delay time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	GTimeWheel.AfterFunc(delay,
		func() {
			queue <- time.Now()
		},
		false,
	)
	return queue
}

func Sleep(delay time.Duration) {
	queue := make(chan bool, 1)
	GTimeWheel.AfterFunc(delay,
		func() {
			queue <- true
		},
		false,
	)
	<-queue
}

func Stop() {
	GTimeWheel.Stop()
}

func notifyChannel(q chan bool) {
	select {
	case q <- true:
	default:
		logger.Debug("[mtime] timer_wheel lost")
	}
}
