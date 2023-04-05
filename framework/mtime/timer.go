package mtime

import (
	"context"
	"time"
)

type Timer struct {
	tw     *TimeWheel
	node   *TimeNode
	C      chan bool
	fn     func()
	stopFn func()

	cancel context.CancelFunc
	Ctx    context.Context
	async  bool
}

func NewTimer(delay time.Duration, async bool) *Timer {
	timer := &Timer{async: async}
	timer.SetTimer(delay)
	return timer
}

func (t *Timer) SetTimer(delay time.Duration) {
	queue := make(chan bool, 1)
	node := GTimeWheel.AfterFunc(delay,
		func() {
			notifyChannel(queue)
		},
		t.async,
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.tw = GTimeWheel
	t.node = node
	t.C = queue
	t.cancel = cancel
	t.Ctx = ctx
}
func (t *Timer) SetAfterFunc(delay time.Duration, callback func()) {
	queue := make(chan bool, 1)
	node := GTimeWheel.AfterFunc(delay,
		func() {
			callback()
			notifyChannel(queue)
		},
		t.async,
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.tw = GTimeWheel
	t.node = node
	t.C = queue
	t.fn = callback
	t.cancel = cancel
	t.Ctx = ctx

}

func (t *Timer) Reset(delay time.Duration) {
	t.node.Stop()
	var node *TimeNode
	if t.fn != nil {
		node = t.tw.AfterFunc(delay,
			func() {
				t.fn()
				notifyChannel(t.C)
			},
			t.async,
		)
	} else {
		node = t.tw.AfterFunc(delay,
			func() {
				notifyChannel(t.C)
			},
			t.async,
		)
	}

	t.node = node
}

func (t *Timer) Stop() {
	if t.stopFn != nil {
		t.stopFn()
	}
	t.node.Stop()
	t.cancel()
}

func (t *Timer) StopFunc(callback func()) {
	t.stopFn = callback
}
