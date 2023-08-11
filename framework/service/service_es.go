package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptsorttype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
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

func (s *Service) ESPutNoId(dbName string, data proto.Message) error {
	res, err := s.ES.Index(dbName).
		Request(data).
		Refresh(refresh.True).
		Do(context.Background())
	if err != nil || (res.Result != result.Created && res.Result != result.Updated) {
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

func (s *Service) ESDel(dbName, id string) error {
	res, err := s.ES.Delete(dbName, id).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es delete got err, db:%s, id:%s", dbName, id)
	}
	if res.Result != result.Deleted {
		return fmt.Errorf("es delete failed. result:%s", res.Result)
	}

	return nil
}

// ESCheckIndex
//
//	@Description: 检查指定的索引是否存在
//	@receiver s
//	@param dbName 索引名
//	@return bool 存在则返回true，否则返回false
func (s *Service) ESCheckIndex(dbName string) bool {
	ok, _ := s.ES.Indices.Exists(dbName).IsSuccess(context.Background())
	return ok
}

func (s *Service) ESDelIndex(dbName string) error {
	if !s.ESCheckIndex(dbName) {
		return nil
	}
	_, err := s.ES.Indices.Delete(dbName).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es delete index  got err, db:%s", dbName)
	}
	return nil
}

// ESMultiSearch
//
//	@Description: ES多条件查找，支持等值和范围条件
//	@param dbName 索引名
//	@param matchMap 等值条件 没有填nil 默认模糊查询,需要精准查找需要在字段后加keyword
//	@param rangeMap 范围条件 没有填nil
//	@param filterMap 排除条件 没有填nil
//	@param hitSize 表示期望命中的数量，真实返回的结果可能会小于该值，[注意：取值0-10000，其他值非法]
//	@param random 是否需要结果随机，[注意：随机是对于整个doc而言，而非按字段]
//	@return error
//	@return *types.HitsMetadata
func (s *Service) ESMultiSearch(dbName string, matchMap map[string]string, rangeMap map[string]RangeItem, filterMap map[string][]string, hitSize int, random bool) (error, *types.HitsMetadata) {
	var (
		query       = make([]types.Query, 0)
		filterQuery = make([]types.Query, 0)
		req         = &search.Request{}
	)

	//if len(matchMap) == 0 && len(rangeMap) == 0 {
	//	return fmt.Errorf("es query param illegal"), nil
	//}

	// 等值条件
	if len(matchMap) > 0 {
		for field, keyword := range matchMap {
			query = append(query, types.Query{Match: map[string]types.MatchQuery{field: types.MatchQuery{Query: keyword}}})
		}
	}
	// 范围条件
	if len(rangeMap) > 0 {
		for field, kvItem := range rangeMap {
			tempItem := kvItem // 临时遍历，规避foreach取地址的bug
			query = append(query, types.Query{Range: map[string]types.RangeQuery{field: types.NumberRangeQuery{
				Gte: (*types.Float64)(&tempItem.Min),
				Lte: (*types.Float64)(&tempItem.Max),
			}}})
		}
	}
	// 排除条件
	if len(filterMap) > 0 {
		for field, keywords := range filterMap {
			filterQuery = append(filterQuery, types.Query{Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{field: keywords},
			}})
		}
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

	t := &types.BoolQuery{}
	if len(query) > 0 {
		t.Must = query
	}
	if len(filterQuery) > 0 {
		t.MustNot = filterQuery
	}
	req.Query = &types.Query{Bool: t} // 条件填充
	req.Size = &hitSize               // 查询数量

	reqStr, _ := json.Marshal(req)
	logger.Debugf("ESMultiSearch 请求: %s", string(reqStr))
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

// ESMultiSearchPage 分页+排序+范围查找
func (s *Service) ESMultiSearchPage(dbName string, rangeMap map[string]RangeItem, hitSize int, sortType *sortorder.SortOrder, from int) (error, *types.HitsMetadata) {
	var (
		query = make([]types.Query, 0)
		req   = &search.Request{}
	)

	// 范围条件
	if len(rangeMap) > 0 {
		for field, kvItem := range rangeMap {
			tempItem := kvItem // 临时遍历，规避foreach取地址的bug
			query = append(query, types.Query{Range: map[string]types.RangeQuery{field: types.NumberRangeQuery{
				Gte: (*types.Float64)(&tempItem.Min),
				Lte: (*types.Float64)(&tempItem.Max),
			}}})
		}
	}

	// 排序
	if sortType != nil {
		req.Sort = []types.SortCombinations{
			types.SortOptions{SortOptions: map[string]types.FieldSort{
				"timeStamp": {Order: sortType},
			}},
		}
	}

	// 分页
	req.From = &from
	req.Size = &hitSize // 查询数量

	reqStr, _ := json.Marshal(req)
	logger.Debugf("ESMultiSearch 请求: %s", string(reqStr))
	// 请求数据
	res, err := s.ES.Search().Index(dbName).Request(req).Do(context.Background())
	if err != nil {
		return errors.Wrapf(err, "es search err,dbName:%s, matchMap:%v", dbName), nil
	}
	if res.TimedOut {
		return DB_ERROR_TIMEOUT, nil
	}
	return nil, &res.Hits
}
