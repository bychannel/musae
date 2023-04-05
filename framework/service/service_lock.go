package service

import (
	"context"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"time"
)

const (
	UnlockResponse_SUCCESS                int32 = 0
	UnlockResponse_LOCK_DOES_NOT_EXIST    int32 = 1
	UnlockResponse_LOCK_BELONGS_TO_OTHERS int32 = 2
	UnlockResponse_INTERNAL_ERROR         int32 = 3
)

var TryLockFail = fmt.Errorf("TryLockFail")
var UnLockFail = fmt.Errorf("UnLockFail")

func (s *Service) TryLock(owner, resId string, expirySec int32) (bool, error) {
	logger.Debug("TryLock Begin:", RedisLock, owner, resId, expirySec)
	var retryTimes int
	err := s.Lock(owner, resId, expirySec)
	if err == TryLockFail && retryTimes < 10 {
		time.Sleep(200 * time.Millisecond)
		logger.Warnf("Retry TryLock idx:%d, err:%v", retryTimes, err)
		err = s.Lock(owner, resId, expirySec)
		retryTimes++
	}
	if err != nil {
		return false, errors.Wrap(err, "TryLock Lock")
	}
	logger.Debug("TryLock End:", RedisLock, owner, resId, expirySec, true)
	return true, nil
}

func (s *Service) Lock(owner, resId string, expirySec int32) error {
	ctx := context.Background()
	r, err := s.Daprc.TryLockAlpha1(ctx, string(RedisLock), &dapr.LockRequest{
		LockOwner:       owner,
		ResourceID:      resId,
		ExpiryInSeconds: expirySec,
	})
	if err != nil {
		return errors.Wrapf(err, "TryLock fail,owner:%s, resId:%s, expiry:%d", owner, resId, expirySec)
	}
	if r == nil {
		return errors.New(fmt.Sprintf("TryLock fail,owner:%s, resId:%s, expiry:%d", owner, resId, expirySec))
	}
	if !r.Success {
		return TryLockFail
	}
	return nil
}

func (s *Service) UnLock(owner, resId string) (bool, error) {
	logger.Debug("UnLock Begin:", RedisLock, owner, resId)
	ctx := context.Background()
	r, err := s.Daprc.UnlockAlpha1(ctx, string(RedisLock), &dapr.UnlockRequest{
		LockOwner:  owner,
		ResourceID: resId,
	})
	if err != nil || r == nil {
		logger.Warnf("UnLock error, err:%v, res:%+v", err, r)
		return false, err
	}

	if r.StatusCode != UnlockResponse_SUCCESS {
		logger.Warnf("UnLock fail, res:%+v", r)
		return false, UnLockFail
	}
	logger.Debug("UnLock End:", RedisLock, owner, resId, r.StatusCode, r.Status)
	return true, nil
}
