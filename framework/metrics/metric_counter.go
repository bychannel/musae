package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"strings"
)

// RegisterCounter
//
//	@Description: 注册Counter
//	@param name
//	@param labels
func RegisterCounter(name CounterType, labels ...string) {
	if baseconf.GetBaseConf() != nil && !baseconf.GetBaseConf().Metric {
		return
	}
	if labels == nil {
		labels = []string{}
	}
	tags := GetMetricMgr().GetServiceTags()
	if len(tags) > 0 {
		for _, tag := range tags {
			kvs := strings.Split(tag, "|")
			if len(kvs) == 2 {
				labels = append(labels, kvs[0])
			}
		}
	}
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      string(name),
		Namespace: NameSpace,
	}, labels)
	prometheus.MustRegister(counter)
	counterMap[name] = counter
}

// CounterInc
//
//	@Description: 增加 1
//	@param name
//	@param labels
func CounterInc(name CounterType, labels ...string) {
	if baseconf.GetBaseConf() != nil && !baseconf.GetBaseConf().Metric {
		return
	}
	if labels == nil {
		labels = []string{}
	}
	tags := GetMetricMgr().GetServiceTags()
	if len(tags) > 0 {
		for _, tag := range tags {
			kvs := strings.Split(tag, "|")
			if len(kvs) == 2 {
				labels = append(labels, kvs[0])
			}
		}
	}
	counter, ok := counterMap[name]
	if !ok {
		return
	}
	counter.WithLabelValues(labels...).Inc()
}

// CounterAdd
//
//	@Description: 增加一个指定的值
//	@param name
//	@param delta
//	@param labels
func CounterAdd(name CounterType, delta int64, labels ...string) {
	if baseconf.GetBaseConf() != nil && !baseconf.GetBaseConf().Metric {
		return
	}
	if labels == nil {
		labels = []string{}
	}
	tags := GetMetricMgr().GetServiceTags()
	if len(tags) > 0 {
		for _, tag := range tags {
			kvs := strings.Split(tag, "|")
			if len(kvs) == 2 {
				labels = append(labels, kvs[0])
			}
		}
	}
	counter, ok := counterMap[name]
	if !ok {
		return
	}
	counter.WithLabelValues(labels...).Add(float64(delta))
}
