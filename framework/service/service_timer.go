package service

import (
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/mtime"
	"sync/atomic"
	"time"
)

func (s *Service) AddTimer(cycle bool, dur time.Duration, cb base.TimerEventCB) (uint64, bool) {

	timerId := atomic.AddUint64(&s.TimerId, 1)
	var node *mtime.TimeNode
	if cycle {
		node = mtime.AddCron(dur, func() {
			select {
			case s.TimeCh <- cb:
			default:
				logger.Warnf("AddTimer timeCh full, %s, %v", s.AppId, dur)
			}

		}, false)
	} else {
		node = mtime.Add(dur, func() {
			defer s.TimerMap.Delete(timerId)
			select {
			case s.TimeCh <- cb:
			default:
				logger.Warnf("AddTimer timeCh full, %s, %v", s.AppId, dur)
			}

		}, false)
	}

	if node == nil {
		logger.Warnf("AddTimer timeCh node nil, %s, %v", s.AppId, dur)
		return 0, false
	}

	s.TimerMap.Store(timerId, node)

	logger.Debugf("AddTimer, %v %v %v", s.AppId, timerId, dur)
	return timerId, true
}

func (s *Service) RemoveTimer(timerId uint64) bool {
	node, ok := s.TimerMap.LoadAndDelete(timerId)
	if !ok {
		logger.Warnf("RemoveTimer TimerMap nothing, %s, %v", s.AppId, timerId)
		return false
	}

	s.TimerMap.Delete(node)
	return true
}
