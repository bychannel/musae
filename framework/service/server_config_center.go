package service

import (
	"context"
	"errors"
	"github.com/dapr/go-sdk/client"
	"github.com/go-redis/redis/v8"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/threading"
	"google.golang.org/grpc/metadata"
	"strconv"
	"strings"
)

func (s *Service) InitConfigCenter() {

	s.SubscribeConfigCenter(baseconf.GetBaseConf().CfgKeys, s.OnCfgCenterCB)

	for _, key := range baseconf.GetBaseConf().CfgKeys {
		val, err := s.GetFromConfigCenter(key)
		if err == nil && val != "" {
			s.CfgKeys.Store(key, val)
		} else {
			if errors.Is(err, DB_ERROR_NOT_EXIST) || strings.Contains(err.Error(), redis.Nil.Error()) {
				logger.Errorf("GetFromConfigCenter [%s] not exist, err: %+v", key, err)
			} else {
				logger.Errorf("GetFromConfigCenter [%s] err:%s", key, err)

			}
		}
	}
	s.CfgKeys.Range(func(key, value any) bool {
		logger.Infof("CfgKeys [%s]:[%s]", key, value)
		return true
	})
}

func (s *Service) SubscribeConfigCenter(keys []string, handler client.ConfigurationHandleFunction) error {
	ctx := context.Background()
	md := metadata.Pairs("dapr-app-id", s.AppId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	threading.GoSafe(func() {
		logger.Infof("subscribe config center keys: %v", keys)
		if err := s.Daprc.SubscribeConfigurationItems(ctx, string(ConfigCenter), keys, func(id string, m map[string]*client.ConfigurationItem) {
			threading.RunSafe(func() {
				for k, v := range m {
					logger.Infof("config center updated id = %s, key = %s, value = %s", id, k, v.Value)
					s.CfgKeys.Store(k, v.Value)
				}
				handler(id, m)
			})
		}); err != nil {
			logger.Errorf("subscribe config center failed key=%+v, err: %v", keys, err)
		}
	})
	return nil
}

func (s *Service) GetFromConfigCenter(key string) (string, error) {
	ctx := context.Background()
	items, err := s.Daprc.GetConfigurationItem(ctx, string(ConfigCenter), key)
	if err != nil {
		return "", err
	}
	if items == nil {
		return "", DB_ERROR_NOT_EXIST
	}
	logger.Debugf("get from config center,%s:%s", key, (*items).Value)
	return (*items).Value, nil
}

func (s *Service) SaveToConfigCenter(key, value string) error {
	//ctx := context.Background()
	//err := s.Daprc.SaveState(ctx, RedisGlobal, key, []byte(value), nil)
	//if err != nil {
	//	return err
	//}

	ctx := context.Background()
	md := metadata.Pairs("saveConfig", key)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// set config value
	s.Redis.Set(ctx, key, value, -1)
	logger.Debugf("save to config center,%s:%s", key, value)
	return nil
}

func (s *Service) GetConfigKeyForInt(key string) (int32, error) {
	if value, ok := s.CfgKeys.Load(key); ok {
		if v, ok1 := value.(int32); ok1 {
			return v, nil
		}
		return 0, DB_ERROR_MARSHAL
	}
	return 0, DB_ERROR_NOT_EXIST
}

func (s *Service) GetConfigKeyForStr(key string) (string, error) {
	if value, ok := s.CfgKeys.Load(key); ok {
		if v, ok1 := value.(string); ok1 {
			return v, nil
		}
		return "", DB_ERROR_MARSHAL
	}
	return "", DB_ERROR_NOT_EXIST
}

func (s *Service) GetConfigKeyForBool(key string) (bool, error) {
	if value, ok := s.CfgKeys.Load(key); ok {
		if v, ok1 := value.(string); ok1 {
			b, err := strconv.ParseBool(v)
			if err == nil {
				return b, nil
			}
		}

		return false, DB_ERROR_MARSHAL
	}

	return false, DB_ERROR_NOT_EXIST
}
