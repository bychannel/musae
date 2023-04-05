package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"strings"
)

//
// RegisterGauge
//  @Description: 普通的 gauge 注册，需要业务主动调用更新接口
//  @param name
//  @param labelNames
//
func RegisterGauge(name GaugeType, reset bool, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      string(name),
		Namespace: NameSpace,
	}, labels)

	prometheus.MustRegister(gauge)
	gaugeMap[name] = gauge

	if reset {
		metricResetMap[string(name)] = reset
	}
}

//
// GaugeSet
//  @Description: 将指标设置为指定的值
//  @param name
//  @param value
//  @param labels
//
func GaugeSet(name GaugeType, value int64, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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
	gauge, ok := gaugeMap[name]
	if !ok {
		return
	}
	gauge.WithLabelValues(labels...).Set(float64(value))
}

//
// GaugeInc
//  @Description: 指标增加 1
//  @param name
//  @param labels
//
func GaugeInc(name GaugeType, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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
	gauge, ok := gaugeMap[name]
	if !ok {
		return
	}
	gauge.WithLabelValues(labels...).Inc()
}

// GaugeAdd
//  @Description: 指标增加指定的值
//  @param name
//  @param value
//  @param labels
//
func GaugeAdd(name GaugeType, value int64, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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

	gauge, ok := gaugeMap[name]
	if !ok {
		return
	}
	gauge.WithLabelValues(labels...).Add(float64(value))
}

//
// GaugeDec
//  @Description: 指标减小 1
//  @param name
//  @param labels
//
func GaugeDec(name GaugeType, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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
	gauge, ok := gaugeMap[name]
	if !ok {
		return
	}
	gauge.WithLabelValues(labels...).Dec()
}

//
// GaugeSub
//  @Description: 指标减小指定的值
//  @param name
//  @param value
//  @param labels
//
func GaugeSub(name GaugeType, value int64, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
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
	gauge, ok := gaugeMap[name]
	if !ok {
		return
	}
	gauge.WithLabelValues(labels...).Sub(float64(value))
}
