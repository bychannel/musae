package service

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"net/http"
	"time"
)

func (s *Service) initES() error {
	cfg := elasticsearch.Config{
		Username: baseconf.GetBaseConf().ESConf.UserName,
		Password: baseconf.GetBaseConf().ESConf.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(baseconf.GetBaseConf().ESConf.Timeout) * time.Second,
		},
	}
	if global.IsCloud {
		cfg.Addresses = baseconf.GetBaseConf().ESConf.Addr
	} else {
		cfg.Addresses = []string{baseconf.GetBaseConf().ESConf.AddrDev}
	}

	var err error
	s.ES, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Index(dbName, id string, data proto.Message) error {
	res, err := s.ES.Index(dbName).
		Id(id).
		Request(data).
		Refresh(refresh.True).
		Do(context.Background())
	if err != nil || (res.Result != result.Created && res.Result != result.Updated) {
		return errors.Wrap(err, res.Result.Name)
	}
	return nil
}

func (s *Service) Get(dbName, id string) (error, []byte) {
	res, err := s.ES.Get(dbName, id).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es get err, db:%s, id:%s", dbName, id), nil
	}
	if !res.Found {
		return DB_ERROR_NOT_EXIST, nil
	}

	return nil, res.Source_
}

func (s *Service) Search(dbName, filed, keyword string) (error, *types.HitsMetadata) {
	res, err := s.ES.Search().
		Index(dbName).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					filed: {Query: keyword},
				},
			},
		}).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es search err,dbName:%s, filed:%s, keyword:%s", dbName, filed, keyword), nil
	}
	if res.TimedOut {
		return DB_ERROR_TIMEOUT, nil
	}
	return nil, &res.Hits
}
