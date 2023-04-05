package mtime

import (
	"context"
	"time"
)

type Ticker struct {
	tw   *TimeWheel
	node *TimeNode
	C    chan bool

	cancel context.CancelFunc
	Ctx    context.Context
	async  bool
}

func NewTicker(delay time.Duration, async bool) *Ticker {
	ticker := &Ticker{async: async}
	ticker.SetTicker(delay)
	return ticker
}

func (t *Ticker) SetTicker(delay time.Duration) {
	queue := make(chan bool, 1)
	node := GTimeWheel.ScheduleFunc(delay,
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

func (t *Ticker) Stop() {
	t.node.Stop()
	t.cancel()
}
