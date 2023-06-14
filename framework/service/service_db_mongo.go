package service

import (
	"context"
	"encoding/json"
	dapr "github.com/dapr/go-sdk/client"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
	"gitlab.musadisca-games.com/wangxw/musae/framework/state"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
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
//// SaveMongoGame save to mongo game
//func (s *Service) SaveMongoGame(key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
//	return s.saveMongo(MongoDbType_MongoGame, key, table, meta, so...)
//}
//
//// GetMongoGame load from mongo game
//func (s *Service) GetMongoGame(key string, meta map[string]string) (*state.KvTable, error) {
//	return s.getMongo(MongoDbType_MongoGame, key, meta)
//}
//
//// SaveMongoMail save to mongo mail
//func (s *Service) SaveMongoMail(key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
//	return s.saveMongo(MongoDbType_MongoMail, key, table, meta, so...)
//}
//
//// GetMongoMail load from mongo mail
//func (s *Service) GetMongoMail(key string, meta map[string]string) (*state.KvTable, error) {
//	return s.getMongo(MongoDbType_MongoMail, key, meta)
//}
//
//// SaveMongoAccount save to mongo account
//func (s *Service) SaveMongoAccount(key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
//	return s.saveMongo(MongoDbType_MongoAccount, key, table, meta, so...)
//}
//
//// GetMongoAccount load from mongo account
//func (s *Service) GetMongoAccount(key string, meta map[string]string) (*state.KvTable, error) {
//	return s.getMongo(MongoDbType_MongoAccount, key, meta)
//}
//
//// SaveMongoGmt save to mongo gmt
//func (s *Service) SaveMongoGmt(key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
//	return s.saveMongo(MongoDbType_MongoGmt, key, table, meta, so...)
//}
//
//// GetMongoGmt load from mongo gmt
//func (s *Service) GetMongoGmt(key string, meta map[string]string) (*state.KvTable, error) {
//	return s.getMongo(MongoDbType_MongoGmt, key, meta)
//}
//
func (s *Service) SaveMongo(db MongoDbType, key string, table *state.KvTable, meta map[string]string, so ...dapr.StateOption) error {
	data, err := json.Marshal(table)
	if err != nil {
		logger.Error("saveMongo Marshal err: ", err, db, key, meta, table.Str())
		return DB_ERROR_MARSHAL
	}

	dataLen := len(table.Data)
	ctx := context.Background()
	now := time.Now()
	//err = s.Daprc.SaveState(ctx, string(db), key, data, meta, so...)
	_, err = utils.RetryDoSyncInterval(
		baseconf.GetBaseConf().MusaeDbSetRetryCount,
		baseconf.GetBaseConf().MusaeDbRetryInterval,
		func() (any, error) {
			return nil, s.Daprc.SaveState(ctx, string(db), key, data, meta, so...)
		})
	if err != nil {
		logger.Error("saveMongo err: ", err, db, key, dataLen, table.Str())
		metrics.GaugeInc(metrics.MongoWErr)
		return DB_ERROR_TIMEOUT
	}
	delay := time.Since(now).Milliseconds()
	metrics.HistogramPut(metrics.MongoWDelayHist, delay, metrics.Mongo)
	metrics.GaugeInc(metrics.MongoWCount)
	metrics.GaugeAdd(metrics.MongoWSize, int64(dataLen))
	logger.WarnDelayf(delay, "SaveMongo db:%v key:%v Delay:%v kvTable:%v", db, key, delay, table.Str())
	return nil
}

func (s *Service) GetMongo(db MongoDbType, key string, meta map[string]string) (*state.KvTable, error) {
	ctx := context.Background()
	now := time.Now()
	//item, err := s.Daprc.GetState(ctx, string(db), key, meta)
	item, err := utils.RetryDoSyncInterval(
		baseconf.GetBaseConf().MusaeDbGetRetryCount,
		baseconf.GetBaseConf().MusaeDbRetryInterval,
		func() (*dapr.StateItem, error) {
			return s.Daprc.GetState(ctx, string(db), key, meta)
		})
	if err != nil {
		logger.Error("getMongo GetState err: ", err)
		metrics.GaugeInc(metrics.MongoRErr)
		return nil, DB_ERROR_TIMEOUT
	}
	delay := time.Since(now).Milliseconds()
	metrics.HistogramPut(metrics.MongoRDelayHist, delay, metrics.Mongo)
	logger.Debugf("getMongo db:%v key:%v len:%v", db, key, len(item.Value))

	metrics.GaugeInc(metrics.MongoRCount)
	// 初始状态
	if len(item.Value) == 0 {
		return nil, DB_ERROR_NOT_EXIST
	}

	table := &state.KvTable{}
	err = json.Unmarshal(item.Value, table)
	if err != nil {
		logger.Errorf("getMongo Unmarshal KvTable err: %v %v %+v", err, key, item)
		return nil, DB_ERROR_UNMARSHAL
	}
	logger.WarnDelayf(delay, "getMongo db:%v key:%v Delay:%v kvTable:%s", db, key, delay, table.Str())

	return table, nil
}

//
//// UpsertMongoTableTransaction update or insert to mongo by transaction
//func (s *Service) UpsertMongoTableTransaction(db MongoDbType, meta map[string]string, kvTableMap map[string]*state.KvTable) error {
//	opts := make([]*dapr.StateOperation, 0)
//
//	for dbKey, kvTable := range kvTableMap {
//		data, err := json.Marshal(kvTable)
//		if err != nil {
//			logger.Errorf("UpsertMongoTableTransaction Marshal err: ", kvTable, err)
//			return err
//		}
//
//		opt := &dapr.StateOperation{
//			Type: dapr.StateOperationTypeUpsert,
//			Item: &dapr.SetStateItem{
//				Key:   dbKey,
//				Value: data,
//				Options: &dapr.StateOptions{
//					Concurrency: dapr.StateConcurrencyLastWrite, // 最终一致性
//					Consistency: dapr.StateConsistencyStrong,
//				},
//			},
//		}
//		opts = append(opts, opt)
//	}
//
//	return s.SaveMongoTransaction(db, meta, opts)
//}
//
//// SaveMongoTransaction save to mongo by transaction
//func (s *Service) SaveMongoTransaction(db MongoDbType, meta map[string]string, opts []*dapr.StateOperation) error {
//	optCount := len(opts)
//
//	ctx, _ := context.WithTimeout(context.Background(), global.DB_INVOKE_TIMEOUT*time.Second)
//	now := time.Now()
//	logger.Debugf("SaveMongoTransaction === db:%v, meta:%v, opts:%v", db, meta, opts)
//	err := s.Daprc.ExecuteStateTransaction(ctx, string(db), meta, opts)
//	if err != nil {
//		logger.Error("SaveMongoTransaction err: ", err, optCount)
//		metrics.GaugeInc(metrics.MongoWErr)
//		return DB_ERROR_TIMEOUT
//	}
//	metrics.HistogramPut(metrics.MongoWDelayHist, time.Since(now).Milliseconds(), metrics.Mongo)
//	metrics.GaugeInc(metrics.MongoWCount)
//	metrics.GaugeAdd(metrics.MongoWSize, int64(optCount))
//	logger.Debugf("SaveMongoTransaction db:[%v], opts: %v, meta: %v", db, opts, meta)
//	return nil
//}
