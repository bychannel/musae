package service

import (
	"context"
	"encoding/json"
	dapr "github.com/dapr/go-sdk/client"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
	"gitlab.musadisca-games.com/wangxw/musae/framework/state"
	"time"
)

//
//import (
//	"context"
//	"encoding/json"
//	dapr "github.com/dapr/go-sdk/client"
//	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
//	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
//	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
//	"gitlab.musadisca-games.com/wangxw/musae/framework/state"
//	"time"
//)
//
//func (s *Service) SaveRedis(key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
//	return s.saveRedis(RedisGlobal, key, table, meta, so...)
//}
//
//func (s *Service) GetRedis(key string, meta map[string]string) (*state.KvTable, error) {
//	return s.getRedis(RedisGlobal, key, meta)
//}
//
func (s *Service) SaveRedis(db RedisDbType, key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
	data, err := json.Marshal(table)
	if err != nil {
		logger.Errorf("saveRedis Marshal err: ", table, err)
		return DB_ERROR_MARSHAL
	}

	dataLen := len(data)
	ctx := context.Background()
	now := time.Now()
	err = s.Daprc.SaveState(ctx, string(db), key, data, meta, so...)
	if err != nil {
		logger.Error("saveRedis err: ", err, db, key, dataLen, meta, table.Str())
		metrics.GaugeInc(metrics.RedisWErr)
		return DB_ERROR_TIMEOUT
	}
	metrics.HistogramPut(metrics.RedisWDelayHist, time.Since(now).Milliseconds(), metrics.Redis)
	metrics.GaugeInc(metrics.RedisWCount)
	metrics.GaugeAdd(metrics.RedisWSize, int64(dataLen))
	logger.Debugf("saveRedis db:[%v], key:[%v], kvTable: %v, meta: %v", db, key, table.Str(), meta)
	return nil
}

func (s *Service) GetRedis(db RedisDbType, key string, meta map[string]string) (*state.KvTable, error) {
	ctx := context.Background()
	now := time.Now()
	item, err := s.Daprc.GetState(ctx, string(db), key, meta)
	if err != nil {
		logger.Error("getRedis GetState err: ", err)
		metrics.GaugeInc(metrics.RedisRErr)
		return nil, DB_ERROR_TIMEOUT
	}
	logger.Debugf("getRedis db:[%v], key:[%v], len:[%v]", db, key, len(item.Value))
	metrics.HistogramPut(metrics.RedisRDelayHist, time.Since(now).Milliseconds(), metrics.Redis)
	metrics.GaugeInc(metrics.RedisRCount)
	// 初始状态
	if len(item.Value) == 0 {
		return nil, DB_ERROR_NOT_EXIST
	}

	table := &state.KvTable{}
	err = json.Unmarshal(item.Value, table)
	if err != nil {
		logger.Errorf("getRedis Unmarshal KvTable err: %v, %v, %+v", err, key, item)
		return nil, DB_ERROR_UNMARSHAL
	}
	logger.Debugf("getRedis db:[%v], key:[%v], kvTable:%s", db, key, table.Str())
	return table, nil
}

//// UpsertRedisTableTransaction update or insert to redis by transaction
//func (s *Service) UpsertRedisTableTransaction(db RedisDbType, meta map[string]string, kvTableMap map[string]*state.KvTable) error {
//	opts := make([]*dapr.StateOperation, 0)
//
//	for dbKey, kvTable := range kvTableMap {
//		data, err := json.Marshal(kvTable)
//		if err != nil {
//			logger.Errorf("UpsertRedisTableTransaction Marshal err: ", kvTable, err)
//			return err
//		}
//
//		opt := &dapr.StateOperation{
//			Type: dapr.StateOperationTypeUpsert,
//			Item: &dapr.SetStateItem{
//				Key:   dbKey,
//				Value: data,
//				Options: &dapr.StateOptions{
//					Concurrency: dapr.StateConcurrencyFirstWrite,
//					Consistency: dapr.StateConsistencyStrong,
//				},
//			},
//		}
//		opts = append(opts, opt)
//	}
//
//	return s.SaveRedisTransaction(db, meta, opts)
//}
//
//// SaveRedisTransaction save to redis by transaction
//func (s *Service) SaveRedisTransaction(db RedisDbType, meta map[string]string, opts []*dapr.StateOperation) error {
//	optCount := len(opts)
//
//	if optCount <= 0 {
//		return nil
//	}
//
//	ctx, _ := context.WithTimeout(context.Background(), global.DB_INVOKE_TIMEOUT*time.Second)
//	now := time.Now()
//	err := s.Daprc.ExecuteStateTransaction(ctx, string(db), meta, opts)
//	if err != nil {
//		logger.Error("SaveRedisTransaction err: ", err, optCount)
//		metrics.GaugeInc(metrics.RedisWErr)
//		return DB_ERROR_TIMEOUT
//	}
//	metrics.HistogramPut(metrics.RedisWDelayHist, time.Since(now).Milliseconds(), metrics.Redis)
//	metrics.GaugeInc(metrics.RedisWCount)
//	metrics.GaugeAdd(metrics.RedisWSize, int64(optCount))
//	logger.Debugf("SaveRedisTransaction db:[%v], opts: %v, meta: %v", db, opts, meta)
//	return nil
//}
