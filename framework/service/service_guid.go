package service

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/guid"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"strconv"
	"time"
)

func (s *Service) DBNext(name string, delta uint64) (uint64, error) {
	ctx := context.Background()
	/*ok, err := s.TryLock(s.AppId, name, LOCK_TTL_SEC)
	defer s.UnLock(s.AppId, name)
	if !ok || err != nil {
		logger.Errorf("DBNext TryLock err: %v, %v, %+v", err, name, s.AppId)
		return 0, err
	}*/
	var f context.CancelFunc
	ctx, f = context.WithTimeout(ctx, global.DB_INVOKE_TIMEOUT*time.Second)
	defer f()
	item, err := s.Daprc.GetStateWithConsistency(ctx, string(MongoDbType_MongoAccount), name, nil, dapr.StateConsistencyStrong)
	if err != nil {
		logger.Errorf("DBNext GetStateWithConsistency err: %v, %v, %+v", err, name, item)
		return 0, err
	}
	curMax, err := strconv.ParseUint(string(item.Value), 10, 64)

	newMax := curMax + delta
	upsers := []*dapr.StateOperation{{
		Type: dapr.StateOperationTypeUpsert,
		Item: &dapr.SetStateItem{
			Key: item.Key,
			Etag: &dapr.ETag{
				Value: item.Etag,
			},
			Value: []byte(strconv.FormatUint(newMax, 10)),
			Options: &dapr.StateOptions{
				Concurrency: dapr.StateConcurrencyFirstWrite,
				Consistency: dapr.StateConsistencyStrong,
			},
		},
	}}

	var f2 context.CancelFunc
	ctx, f2 = context.WithTimeout(ctx, global.DB_INVOKE_TIMEOUT*time.Second)
	defer f2()
	err = s.Daprc.ExecuteStateTransaction(ctx, string(MongoDbType_MongoAccount), nil, upsers)
	if err != nil {
		logger.Errorf("DBNext  SaveState err: %v, %v, %+v", err, name, item)
		return 0, err
	}
	logger.Debugf("GUIDNext, curMax:%v, delta:%v, newMax:%v", curMax, delta, newMax)
	return newMax, nil
}

func (s *Service) GenGUID(typ guid.GUID_TYPE) uint64 {
	id, err := s.GuidPool.NextByPool(typ)
	if id == 0 || err != nil {
		id, err = s.GuidPool.NextByPool(typ)
		if id == 0 || err != nil {
			logger.Errorf("GenGUID retry err:%+v", err)
			return 0
		}
	}
	return id
}
