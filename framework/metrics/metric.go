package metrics

import (
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// service metric on prometheus

type GaugeType string
type CounterType string
type HistogramType string
type BlockingType string

type GaugeMap map[GaugeType]*prometheus.GaugeVec
type CounterMap map[CounterType]*prometheus.CounterVec
type HistogramMap map[HistogramType]*prometheus.HistogramVec
type MetricBlockMap map[BlockingType]*MetricBlock
type MetricResetMap map[string]bool
type RegisterMetricFunc func()

const (
	NameSpace  = "aniwar"
	MetricName = "metric"
	MetricTags = "tags"
)

const (
	LOG_SIZE_CUT_TYPE = 0
	LOG_TIME_CUT_TYPE = 1
)

var (
	GaugeFuncCallInterval = 10 // gauge func 函数调用时间间隔，单位秒
	OutputInterval        = 10 // tick logger 输出的时间间隔，单位秒
)

// defaultBucket 延统计的 bucket 范围,单位: ms
var defaultBucket = []float64{5, 10, 25, 50, 75, 100, 150, 250, 500, 750, 1000, 1500, 2000, 3500, 5000}

var counterMap CounterMap
var gaugeMap GaugeMap
var histogramMap HistogramMap
var metricBlockMap MetricBlockMap
var metricResetMap MetricResetMap
var customMetricRegisterFunc RegisterMetricFunc // 业务层metric注册函数
var metricMgr MetricMgr

type MetricBlock struct {
	Name      BlockingType
	Histogram *prometheus.HistogramVec
	Gauge     *prometheus.GaugeVec
	Counter   *prometheus.CounterVec
}

// BlockStub
//  BlockStub
//  @Description: 监控时的打点类，记录监控起点，终点
//
type BlockStub struct {
	BeginTime   time.Time    // 打点起始时间
	Success     bool         // 请求是否成功
	Labels      []string     // 打点的标签内容
	MetricBlock *MetricBlock // 所操作的 metric 组件
}

//=========================== MetricMgr 管理器 =============================//

type MetricMgr struct {
	serviceTags       []string //当前进程服务
	labelMap          map[string]string
	monotonicTimeFunc func() int64
	loggerDone        chan struct{}
	log               *zap.SugaredLogger
	logPath           string
	logFile           string
}

func (m *MetricMgr) Init(logPath, fileName string, service []string, monotonicTimeFunc func() int64) error {
	m.serviceTags = service
	m.labelMap = map[string]string{}
	m.loggerDone = make(chan struct{}, 1)
	m.monotonicTimeFunc = monotonicTimeFunc
	if err := m.InitLog(logPath, fileName); err != nil {
		return err
	}
	var labels []*io_prometheus_client.LabelPair
	if len(m.serviceTags) > 0 {
		for _, tag := range m.serviceTags {
			kvs := strings.Split(tag, "|")
			if len(kvs) == 2 {
				labels = append(labels, &io_prometheus_client.LabelPair{Name: &kvs[0], Value: &kvs[1]})
			}
		}
	}
	for _, label := range labels {
		m.labelMap[label.GetName()] = label.GetValue()
	}
	return nil
}

func (m *MetricMgr) InitLog(logPath, fileName string) error {

	if logPath == "" {
		return errors.New(fmt.Sprint("dlog init failed,log file is nil!", logPath, fileName))
	}

	if utils.PathExists(logPath) == false {
		os.MkdirAll(logPath, 0755)
	}

	baseConf := baseconf.GetBaseConf()
	if baseConf == nil {
		return fmt.Errorf("baseConf is nil")
	}

	defer func() {
		//svcPath := filepath.Join(logPath, "/", service)
		logFile := fmt.Sprintf("metric-%s-%v.prom", global.AppID, global.SID)
		GetMetricMgr().Start(logPath, logFile, global.IsCloud)
	}()

	if global.IsCloud {
		return nil
	}

	if fileName == "" {
		fileName = utils.GetProcName()
	}

	//fileName = "metric-" + fileName + "-" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10) +
	//	"-" + strconv.Itoa(os.Getpid())

	atomLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	var syncWriter zapcore.WriteSyncer
	cutType := baseconf.GetBaseConf().LogCutType
	switch cutType {
	case LOG_SIZE_CUT_TYPE:
		syncWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logPath, "metric-"+fileName+".log"),
			MaxAge:     baseconf.GetBaseConf().LogMaxAges,
			MaxBackups: baseconf.GetBaseConf().LogMaxBackups,
			MaxSize:    baseconf.GetBaseConf().LogMaxSize,
			Compress:   baseconf.GetBaseConf().LogCompress,
		})
	case LOG_TIME_CUT_TYPE:
		logWriter, err := rotatelogs.New(
			filepath.Join(logPath, "metric-"+fileName+"-%Y%m%d%H%M%S.log"),
			rotatelogs.WithMaxAge(time.Duration(baseconf.GetBaseConf().LogMaxAges)*time.Hour*24),
			rotatelogs.WithRotationTime(time.Minute*time.Duration(baseconf.GetBaseConf().LogRotationTime)),
			rotatelogs.WithRotationSize(int64(baseconf.GetBaseConf().LogMaxSize*1024*1024)),
		)
		if err != nil {
			return err
		}
		syncWriter = zapcore.AddSync(logWriter)
	}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		//TimeKey:        "T",
		//LevelKey:       "L",
		//NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(encoder, syncWriter, atomLevel)
	var options []zap.Option
	m.log = zap.New(core, options...).Sugar()

	return nil
}

func (m *MetricMgr) Start(svcPath, logFile string, isk8s bool) {
	m.logPath = svcPath
	m.logFile = logFile
	if utils.PathExists(m.logPath) == false {
		os.MkdirAll(m.logPath, 0755)
	}
	go m.tickLogger(isk8s)
}

//
// tickLogger
//  @Description: 定期输出打点日志内容
//  @receiver m
//
func (m *MetricMgr) tickLogger(isk8s bool) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("tickLogger recover:", err)
		}
	}()

	ticker := time.NewTicker(time.Duration(OutputInterval) * time.Second) //TickLoggerInterval s 输出一次
	if !isk8s && m.log == nil {
		fmt.Println("tickLogger : init metric error m.log is nil")
		return
	}
	for {
		select {
		case <-ticker.C:
			if baseconf.GetBaseConf().Metric {
				func() {
					defer func() {
						if err := recover(); err != any(nil) {
							fmt.Println("tickLogger recover, err: ", err)
						}
					}()
					if isk8s {
						logFile := filepath.Join(m.logPath, m.logFile)
						tmpPromFile := logFile + ".tmp"
						OutputMetricLogForExporter(tmpPromFile)
						os.Rename(tmpPromFile, logFile)
					} else {
						OutputMetricLog(m.log)
					}
				}()
			}
		case <-m.loggerDone:
			fmt.Println("tickLogger stop logger")
			return
		}
	}
}

func (m *MetricMgr) Stop() {
	//logger.Info("[Metrics] ready to stop")
	m.loggerDone <- struct{}{}
}

func (m *MetricMgr) GetServiceTags() []string {
	return m.serviceTags
}

func (m *MetricMgr) GetLabelMap() map[string]string {
	return m.labelMap
}

func (m *MetricMgr) GetMonotonicTime() int64 {
	return m.monotonicTimeFunc()
}

func GetMetricMgr() *MetricMgr {
	return &metricMgr
}

func Final() {
	GetMetricMgr().Stop()
}

func InitMetric(logPath, fileName string, service []string, registerFunc RegisterMetricFunc, monotonicTimeFunc func() int64) error {
	gaugeMap = make(GaugeMap)
	counterMap = make(CounterMap)
	histogramMap = make(HistogramMap)
	metricBlockMap = make(MetricBlockMap)
	metricResetMap = make(MetricResetMap)
	customMetricRegisterFunc = registerFunc
	err := GetMetricMgr().Init(logPath, fileName, service, monotonicTimeFunc)
	if err != nil {
		return err
	}
	RegisterMetrics()
	if customMetricRegisterFunc != nil {
		customMetricRegisterFunc()
	}

	//logger.Infof("BaseConf: %v", baseconf.GetBaseConf())
	//logger.Info("[Metrics] init metric success")
	return nil
}
