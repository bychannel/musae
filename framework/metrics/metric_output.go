package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"go.uber.org/zap"
	"os"
	"sort"
	"strings"
	"time"
)

//
// CalculateHistogram
//  @Description: 直方图
//  @param histogram
//  @return uint64
//  @return uint64
//

// 默认记录方式 默认为 Node_Exporeter 上报格式
var DEFAULT_RECORD_KIND = Metric_Type_Node_Exporter

const (
	Metric_Type_Remote_Write  int = iota // 程序初始上报格式
	Metric_Type_Node_Exporter            // Node_Exporter agent 上报格式
)

// 只生成一次 #HELP所使用的去重map
var NameRecord = make(map[string]struct{})

func CalculateHistogram(histogram *io_prometheus_client.Histogram) (uint64, uint64) {

	if histogram == nil {
		return 0, 0
	}

	type Bucket struct {
		Cumulative uint64
		UpperBound uint64
	}

	total := uint64(0)
	buckets := make([]Bucket, len(histogram.GetBucket()))
	for idx, metricBucket := range histogram.GetBucket() {
		if idx >= len(buckets) {
			break
		}
		buckets[idx].Cumulative = metricBucket.GetCumulativeCount()
		buckets[idx].UpperBound = uint64(metricBucket.GetUpperBound())
		if total < metricBucket.GetCumulativeCount() {
			total = metricBucket.GetCumulativeCount()
		}
	}

	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].UpperBound < buckets[j].UpperBound
	})

	rankMedium := total >> 1
	rank95 := uint64(float64(total) * 0.95)

	calculateSpecifiedVal := func(startBucketIdx int, targetRank uint64) (uint64, int) {
		var lastBucketIdx = -1
		if startBucketIdx > 0 {
			lastBucketIdx = startBucketIdx - 1
		}

		var curBucketIdx int
		var targetVal uint64
		for i := startBucketIdx; i < len(buckets); i++ {
			if buckets[i].Cumulative < targetRank {
				lastBucketIdx = i
				continue
			}

			lastBucketBound := uint64(0)
			bucketCount := buckets[i].Cumulative
			targetCountInBucket := targetRank
			bucketBoundSpan := buckets[i].UpperBound
			if lastBucketIdx != -1 {
				bucketCount = buckets[i].Cumulative - buckets[lastBucketIdx].Cumulative
				targetCountInBucket = targetRank - buckets[lastBucketIdx].Cumulative
				bucketBoundSpan = buckets[i].UpperBound - buckets[lastBucketIdx].UpperBound
				lastBucketBound = buckets[lastBucketIdx].UpperBound
			}

			if bucketCount == 0 {
				break
			}

			targetVal = uint64((float64(targetCountInBucket)/float64(bucketCount))*float64(bucketBoundSpan) + float64(lastBucketBound))
			curBucketIdx = i
			break
		}
		return targetVal, curBucketIdx
	}
	mediumVal, curBucketIdx := calculateSpecifiedVal(0, rankMedium)
	per95Val, _ := calculateSpecifiedVal(curBucketIdx, rank95)

	return mediumVal, per95Val
}

type MetricPoint struct {
	Metric  string            `json:"metric"` // 指标名称
	TagsMap map[string]string `json:"tags"`   // 数据标签
	Time    int64             `json:"time"`   // unix时间戳
	Value   float64           `json:"value"`
}

func writeLog(log *zap.SugaredLogger, metric *MetricPoint) {
	b, err := json.Marshal(metric)
	if err == nil && len(b) > 0 {
		log.Info(string(b))
	}
}

// OutputMetricLog
//
//	@Description: 采集信息写入log,异步推送到prometheus
//	@param log
func OutputMetricLog(log *zap.SugaredLogger) {
	if baseconf.GetBaseConf() != nil && !baseconf.GetBaseConf().Metric {
		return
	}
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return
	}
	now := time.Now().Unix()
	labelMap := GetMetricMgr().GetLabelMap()
	for _, mf := range mfs {
		for _, metric := range mf.GetMetric() {
			name := mf.GetName()
			switch mf.GetType() {
			case io_prometheus_client.MetricType_GAUGE:
				writeLog(log, &MetricPoint{Metric: name, TagsMap: labelMap, Time: now, Value: metric.GetGauge().GetValue()})
				namespaceTag := fmt.Sprintf("%s_", NameSpace)
				if strings.Contains(name, namespaceTag) {
					metricName := strings.Replace(name, namespaceTag, "", 1)
					_, ok := metricResetMap[metricName]
					if ok {
						GaugeSet(GaugeType(metricName), 0)
					}
				}
			case io_prometheus_client.MetricType_COUNTER:
				writeLog(log, &MetricPoint{Metric: name, TagsMap: labelMap, Time: now, Value: metric.GetCounter().GetValue()})
			case io_prometheus_client.MetricType_HISTOGRAM:
				avgVal := uint64(metric.GetHistogram().GetSampleSum()) / metric.GetHistogram().GetSampleCount()
				mediumVal, percent95Val := CalculateHistogram(metric.GetHistogram())
				writeLog(log, &MetricPoint{Metric: fmt.Sprintf("%s_%s", name, "count"), TagsMap: labelMap, Time: now,
					Value: float64(metric.GetHistogram().GetSampleCount())})
				writeLog(log, &MetricPoint{Metric: fmt.Sprintf("%s_%s", name, "average"), TagsMap: labelMap, Time: now,
					Value: float64(avgVal)})
				writeLog(log, &MetricPoint{Metric: fmt.Sprintf("%s_%s", name, "medium"), TagsMap: labelMap, Time: now,
					Value: float64(mediumVal)})
				writeLog(log, &MetricPoint{Metric: fmt.Sprintf("%s_%s", name, "percent95"), TagsMap: labelMap, Time: now,
					Value: float64(percent95Val)})
			}
		}
	}
}

func OutputMetricLogForExporter(logFile string) {
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return
	}

	//now := time.Now().Unix()

	for _, mf := range mfs {
		for _, metric := range mf.GetMetric() {
			name := mf.GetName()
			/*namespaceTag := fmt.Sprintf("%s_", NameSpace)
			if !strings.Contains(name, namespaceTag) {
				continue
			}
			labels := metric.GetLabel()
			if len(labels) == 0 {
				//labelName := "server"
				tags := GetMetricMgr().GetServiceTags()
				if len(tags) > 0 {
					for _, tag := range tags {
						kvs := strings.Split(tag, "|")
						if len(kvs) == 2 {
							labels = append(labels, &io_prometheus_client.LabelPair{Name: &kvs[0], Value: &kvs[1]})
						}
					}
				}
				//labels = []*io_prometheus_client.LabelPair{{Name: &labelName, Value: &labelValue}}
			}
			labelMap := map[string]string{}
			for _, label := range labels {
				labelMap[label.GetName()] = label.GetValue()
			}
			*/
			var strTag string
			labelMap := GetMetricMgr().GetLabelMap()
			for key, value := range labelMap {
				strTag = fmt.Sprintf("%s%s=\"%s\",", strTag, key, value)
			}

			//switch mf.GetType() {
			//case io_prometheus_client.MetricType_GAUGE:
			//	//switch DEFAULT_RECORD_KIND {
			//	//case Metric_Type_Remote_Write:
			//	//case Metric_Type_Node_Exporter:
			//	str := fmt.Sprintf("%s{%s} %v\n", name, strTag, metric.GetGauge().GetValue())
			//	_, err := file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	//log.Info(str)
			//	//}
			//
			//case io_prometheus_client.MetricType_COUNTER:
			//	//switch DEFAULT_RECORD_KIND {
			//	//case Metric_Type_Remote_Write:
			//	//
			//	//case Metric_Type_Node_Exporter:
			//	str := fmt.Sprintf("%s{%s} %v\n", name, strTag, metric.GetGauge().GetValue())
			//	_, err := file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	//log.Info(str)
			//	//}
			//
			//case io_prometheus_client.MetricType_HISTOGRAM:
			//	avgVal := uint64(metric.GetHistogram().GetSampleSum()) / metric.GetHistogram().GetSampleCount()
			//	mediumVal, percent95Val := CalculateHistogram(metric.GetHistogram())
			//
			//	//switch DEFAULT_RECORD_KIND {
			//	//case Metric_Type_Remote_Write:
			//	//
			//	//case Metric_Type_Node_Exporter:
			//
			//	str := fmt.Sprintf("%s_%s{%s} %v\n", name, "count", strTag, metric.GetHistogram().GetSampleCount())
			//	_, err := file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	str = fmt.Sprintf("%s_%s{%s} %v\n", name, "average", strTag, avgVal)
			//	_, err = file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	str = fmt.Sprintf("%s_%s{%s} %v\n", name, "medium", strTag, mediumVal)
			//	_, err = file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	str = fmt.Sprintf("%s_%s{%s} %v\n", name, "percent95", strTag, percent95Val)
			//	_, err = file.WriteString(str)
			//	if err != nil {
			//		return
			//	}
			//	//}
			//}
			switch mf.GetType() {
			case io_prometheus_client.MetricType_GAUGE:
				//switch DEFAULT_RECORD_KIND {
				//case Metric_Type_Remote_Write:
				//case Metric_Type_Node_Exporter:
				str := fmt.Sprintf("# HELP %s record %s\n# Type %s gauge\n%s{%s} %v\n", name, name, name, name, strTag, metric.GetGauge().GetValue())
				_, err := file.WriteString(str)
				if err != nil {
					return
				}
				//log.Info(str)
				//}

			case io_prometheus_client.MetricType_COUNTER:
				//switch DEFAULT_RECORD_KIND {
				//case Metric_Type_Remote_Write:
				//
				//case Metric_Type_Node_Exporter:
				str := fmt.Sprintf("# HELP %s record %s\n# Type %s counter\n%s{%s} %v\n", name, name, name, name, strTag, metric.GetGauge().GetValue())
				_, err := file.WriteString(str)
				if err != nil {
					return
				}
				//log.Info(str)
				//}

			case io_prometheus_client.MetricType_HISTOGRAM:
				avgVal := uint64(metric.GetHistogram().GetSampleSum()) / metric.GetHistogram().GetSampleCount()
				mediumVal, percent95Val := CalculateHistogram(metric.GetHistogram())

				//switch DEFAULT_RECORD_KIND {
				//case Metric_Type_Remote_Write:
				//
				//case Metric_Type_Node_Exporter:

				str := fmt.Sprintf("# HELP %s record %s\n# Type %s histogram\n%s_%s{%s} %v\n", name, name, name, name, "count", strTag, metric.GetHistogram().GetSampleCount())
				_, err = file.WriteString(str)
				if err != nil {
					return
				}
				str = fmt.Sprintf("%s_%s{%s} %v\n", name, "average", strTag, avgVal)
				_, err = file.WriteString(str)
				if err != nil {
					return
				}
				str = fmt.Sprintf("%s_%s{%s} %v\n", name, "medium", strTag, mediumVal)
				_, err = file.WriteString(str)
				if err != nil {
					return
				}
				str = fmt.Sprintf("%s_%s{%s} %v\n", name, "percent95", strTag, percent95Val)
				_, err = file.WriteString(str)
				if err != nil {
					return
				}
				//}
			}
		}
	}
}
