package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"strings"
	"time"
)

// RegisterBlockMetric
// RegisterReqMetric
//  @Description: 阻塞类信息采集： 数据库,rpc,http,包含时延,并发数,失败数的统计
//  @param name
//  @param labels
//
func RegisterBlockMetric(name BlockingType, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
		return
	}
	impl := &MetricBlock{Name: name}
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

	impl.Counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: string(name),
		Name:      TagFailed,
		Namespace: NameSpace,
	}, labels)
	impl.Gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: string(name),
		Name:      TagConcurrent,
		Namespace: NameSpace,
	}, labels)
	impl.Histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: string(name),
		Name:      TagDelay,
		Buckets:   defaultBucket,
		Namespace: NameSpace,
	}, labels)
	prometheus.MustRegister(impl.Histogram, impl.Gauge, impl.Counter)

	metricBlockMap[name] = impl
}

//
// calculateMetric
//  @Description:
//  @receiver m
//  @param latency
//  @param success
//  @param labels
//
func (m *MetricBlock) calculateMetric(latency int64, success bool, labels ...string) {
	if !baseconf.GetBaseConf().Metric {
		return
	}
	if labels == nil {
		labels = []string{}
	}
	//labels = append(labels, GetMetricMgr().GetServiceTag())
	// 时延
	m.Histogram.WithLabelValues(labels...).Observe(float64(latency))
	// 并发数
	m.Gauge.WithLabelValues(labels...).Dec()
	// 失败数量
	if !success {
		m.Counter.WithLabelValues(labels...).Inc()
	}
}

//
// StartBlock
//  @Description:
//  @param stub
//  @param name
//  @param labels
//  @return *BlockStub
//
func StartBlock(stub *BlockStub, name BlockingType, labels ...string) *BlockStub {
	if !baseconf.GetBaseConf().Metric {
		return nil
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
	metricBlock, ok := metricBlockMap[name]
	if !ok {
		return &BlockStub{MetricBlock: &MetricBlock{}}
	}
	stub.Labels = labels
	metricBlock.Gauge.WithLabelValues(stub.Labels...).Inc()
	stub.BeginTime = time.Now()
	stub.MetricBlock = metricBlock
	return stub
}

// End 结束 Blockuest 打点
func (s *BlockStub) End(success bool) {
	if !baseconf.GetBaseConf().Metric {
		return
	}
	s.Success = success
	deltaTime := time.Since(s.BeginTime).Milliseconds()
	s.MetricBlock.calculateMetric(deltaTime, s.Success, s.Labels...)
}
