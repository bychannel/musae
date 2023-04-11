package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"strings"
)

// RegisterHistogram
//
//	@Description: 直方图分布统计
//	@param name
//	@param bucket
//	@param labels
func RegisterHistogram(name HistogramType, bucket []float64, labels ...string) {
	if baseconf.GetBaseConf() != nil && !baseconf.GetBaseConf().Metric {
		return
	}
	if bucket == nil {
		bucket = defaultBucket
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
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      string(name),
		Buckets:   bucket,
		Namespace: NameSpace,
	}, labels)

	prometheus.MustRegister(histogram)
	histogramMap[name] = histogram
	//metricResetMap[string(name)] = false
}

// HistogramPut
//
//	@Description: 添加数据到直方图统计
//	@param name
//	@param value
//	@param labels
func HistogramPut(name HistogramType, value int64, labels ...string) {
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
	histogram, ok := histogramMap[name]
	if !ok {
		return
	}
	histogram.WithLabelValues(labels...).Observe(float64(value))
}
