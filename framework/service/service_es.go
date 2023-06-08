package service

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptsorttype"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"net/http"
	"time"
)

type RangeItem struct {
	Min float64
	Max float64
}

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

func (s *Service) ESPut(dbName, id string, data proto.Message) error {
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

func (s *Service) ESGet(dbName, id string) (error, []byte) {
	res, err := s.ES.Get(dbName, id).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es get err, db:%s, id:%s", dbName, id), nil
	}
	if !res.Found {
		return DB_ERROR_NOT_EXIST, nil
	}

	return nil, res.Source_
}

// ESMultiSearch
//
//	@Description: ES多条件查找，支持等值和范围条件
//	@param dbName 索引名
//	@param matchMap 等值条件 没有填nil
//	@param rangeMap 范围条件 没有填nil
//	@param hitSize 表示期望命中的数量，真实返回的结果可能会小于该值，[注意：取值0-10000，其他值非法]
//	@param random 是否需要结果随机，[注意：随机是对于整个doc而言，而非按字段]
//	@return error
//	@return *types.HitsMetadata
func (s *Service) ESMultiSearch(dbName string, matchMap map[string]string, rangeMap map[string]RangeItem, hitSize int, random bool) (error, *types.HitsMetadata) {
	var (
		matchOpt = make(map[string]types.MatchQuery)
		rangeOpt = make(map[string]types.RangeQuery)
		query    = make([]types.Query, 0)
		req      = &search.Request{}
	)

	if len(matchMap) == 0 && len(rangeMap) == 0 {
		return fmt.Errorf("es query param illegal"), nil
	}

	// 等值条件
	if len(matchMap) > 0 {
		for field, keyword := range matchMap {
			matchOpt[field] = types.MatchQuery{Query: keyword}
		}
		query = append(query, types.Query{Match: matchOpt})
	}
	// 范围条件
	if len(rangeMap) > 0 {
		for field, kvItem := range rangeMap {
			rangeOpt[field] = types.NumberRangeQuery{
				Gte: (*types.Float64)(&kvItem.Min),
				Lte: (*types.Float64)(&kvItem.Max),
			}
		}
		query = append(query, types.Query{Range: rangeOpt})
	}

	// 排序，随机配置
	if random {
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				Script_: &types.ScriptSort{
					Type:   &scriptsorttype.Number,
					Script: &types.InlineScript{Source: "Math.random()"},
				},
			},
		}
	}

	req.Query = &types.Query{Bool: &types.BoolQuery{Must: query}} // 等值条件
	req.Size = &hitSize                                           // 查询数量

	logger.Debugf("ESMultiSearch 请求: %+v", query)
	// 请求数据
	res, err := s.ES.Search().Index(dbName).Request(req).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es search err,dbName:%s, matchMap:%v", dbName, matchMap), nil
	}
	if res.TimedOut {
		return DB_ERROR_TIMEOUT, nil
	}
	return nil, &res.Hits
}
