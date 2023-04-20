package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/http"
	"gitlab.musadisca-games.com/wangxw/musae/framework/metrics"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var sugar *zap.SugaredLogger

// var _log *zap.Logger
var atom zap.AtomicLevel
var RedisCli *redis.Client

const (
	LOG_SIZE_CUT_TYPE = 0
	LOG_TIME_CUT_TYPE = 1
)

func init() {
	atom = zap.NewAtomicLevelAt(zap.DebugLevel)
}

func Caller(l zapcore.Level) string {
	funcName, file, line, ok := runtime.Caller(3)
	if ok {
		var level string
		switch l {
		case zap.DebugLevel:
			level = "DEBUG"
		case zap.InfoLevel:
			level = "INFO"
		case zap.WarnLevel:
			level = "WARN"
		case zap.ErrorLevel:
			level = "ERROR"
		case zap.FatalLevel:
			level = "FATAL"
		}
		return fmt.Sprintf("%s %s %s %s %s:%d%s ", time.Now().Format("2006-01-02 15:04:05.000 -0700 MST"), global.HostName, global.AppID, level, path.Base(file), line, path.Ext(runtime.FuncForPC(funcName).Name()))
	}
	return ""
}

func GetProcName() string {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("logger init failed")
		}
	}()

	dir, err := os.Executable()
	if err != nil {
		panic(any(err))
	}

	var realPath string
	realPath, err = filepath.EvalSymlinks(dir)
	if err != nil {
		panic(any(err))
	}

	filename := filepath.Base(realPath)
	baseName := strings.TrimSuffix(filename, path.Ext(filename))

	return strings.ToLower(baseName)
}

func ResetLogLevel(l zapcore.Level) {
	atom.SetLevel(l)
}

func Init(logPath, fileName string) error {

	if logPath == "" {
		return fmt.Errorf("log path error:%s", logPath)
	}

	if fileName == "" {
		fileName = GetProcName()
	}

	//_log, _ = zap.NewProduction(zap.AddCaller())
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		//TimeKey:        "T",
		//LevelKey: "L",
		//NameKey:        "N",
		//CallerKey:      "C",
		//FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	if utils.PathExists(logPath) == false {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			fmt.Printf("create logdir:%s err:%v \n", logPath, err)
			return err
		}
	}

	var (
		syncWriter zapcore.WriteSyncer
		//bws              *zapcore.BufferedWriteSyncer
		//bufSize          int
		//bufFlushInterval time.Duration
	)

	//if baseconf.GetBaseConf() != nil {
	//	bufSize = baseconf.GetBaseConf().LogBufSize * 1024
	//	bufFlushInterval = time.Duration(baseconf.GetBaseConf().LogFlushInterval) * time.Second
	//}
	if !global.IsCloud {
		//fileName = fmt.Sprintf("%s-%v", fileName, os.Getpid())
		//bufSize = 1 * 1024                 // 最小为1kb
		//bufFlushInterval = 1 * time.Second // 最小为1s
	}

	if baseconf.GetBaseConf() != nil {
		cutType := baseconf.GetBaseConf().LogCutType
		switch cutType {
		case LOG_SIZE_CUT_TYPE:
			syncWriter = zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(logPath, fileName+".log"),
				MaxAge:     baseconf.GetBaseConf().LogMaxAges,
				MaxBackups: baseconf.GetBaseConf().LogMaxBackups,
				MaxSize:    baseconf.GetBaseConf().LogMaxSize,
				Compress:   baseconf.GetBaseConf().LogCompress,
			})
		case LOG_TIME_CUT_TYPE:
			logWriter, err := rotatelogs.New(
				filepath.Join(logPath, fileName+"-%Y%m%d%H%M%S.log"),
				rotatelogs.WithMaxAge(time.Duration(baseconf.GetBaseConf().LogMaxAges)*time.Hour*24),
				rotatelogs.WithRotationTime(time.Minute*time.Duration(baseconf.GetBaseConf().LogRotationTime)),
				rotatelogs.WithRotationSize(int64(baseconf.GetBaseConf().LogMaxSize*1024*1024)),
			)
			if err != nil {
				return err
			}
			syncWriter = zapcore.AddSync(logWriter)
		}
		//bws = &zapcore.BufferedWriteSyncer{
		//	WS:            syncWriter,
		//	Size:          bufSize,
		//	FlushInterval: bufFlushInterval,
		//}
	} else {
		syncWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logPath, fileName+".log"),
			MaxAge:     30,
			MaxBackups: 100,
			MaxSize:    1024,
			Compress:   false,
		})
		//bws = &zapcore.BufferedWriteSyncer{
		//	WS:            syncWriter,
		//	Size:          bufSize,
		//	FlushInterval: bufFlushInterval,
		//}
	}

	core := zapcore.NewCore(encoder, syncWriter, zap.DebugLevel)
	sugar = zap.New(core, zap.AddCaller()).Sugar()

	return nil
}

func ToString(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func IsArgsTooLang(l zapcore.Level, str string) bool {
	if baseconf.GetBaseConf() == nil {
		return false
	}
	maxLen := baseconf.GetBaseConf().LogMaxLen
	strLen := len(str)
	if maxLen > 0 && strLen > maxLen {
		fragmentNum := strLen / maxLen
		fragmentLast := strLen % maxLen
		for i := 0; i < fragmentNum; i++ {
			fmt.Println(str[i*maxLen : (i+1)*maxLen])
			switch l {
			case zap.DebugLevel:
				sugar.Debug(str[i*maxLen : (i+1)*maxLen])
			case zap.InfoLevel:
				sugar.Info(str[i*maxLen : (i+1)*maxLen])
			case zap.WarnLevel:
				sugar.Warn(str[i*maxLen : (i+1)*maxLen])
			case zap.ErrorLevel:
				sugar.Error(str[i*maxLen : (i+1)*maxLen])
			case zap.FatalLevel:
				sugar.Fatal(str[i*maxLen : (i+1)*maxLen])
			}
		}
		if fragmentLast > 0 {
			fmt.Println(str[fragmentNum*maxLen : strLen])
			switch l {
			case zap.DebugLevel:
				sugar.Debug(str[fragmentNum*maxLen : strLen])
			case zap.InfoLevel:
				sugar.Info(str[fragmentNum*maxLen : strLen])
			case zap.WarnLevel:
				sugar.Warn(str[fragmentNum*maxLen : strLen])
			case zap.ErrorLevel:
				sugar.Error(str[fragmentNum*maxLen : strLen])
			case zap.FatalLevel:
				sugar.Fatal(str[fragmentNum*maxLen : strLen])
			}
		}
		return true
	}

	return false
}

func SaveToRedis(log ...string) {
	if RedisCli != nil && len(baseconf.GetBaseConf().RedisLogKey) > 0 {
		ctx, f := context.WithTimeout(context.Background(), global.DB_INVOKE_TIMEOUT*time.Second)
		defer f()
		RedisCli.RPush(ctx, baseconf.GetBaseConf().RedisLogKey, log)
	}
	for _, k := range log {
		fmt.Println(k)
	}
}

func Debug(args ...interface{}) {
	DebugA(args...)
}

func DebugA(args ...interface{}) {
	if atom.Enabled(zap.DebugLevel) {
		log := Caller(zap.DebugLevel) + fmt.Sprintln(args)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		if IsArgsTooLang(zap.DebugLevel, log) {
			return
		}
		sugar.Debug(log)
	}
}

func Debugf(template string, args ...interface{}) {
	DebugAf(template, args...)
}

func DebugAf(template string, args ...interface{}) {
	if atom.Enabled(zap.DebugLevel) {
		log := Caller(zap.DebugLevel) + fmt.Sprintf(template, args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		if IsArgsTooLang(zap.DebugLevel, log) {
			return
		}
		sugar.Debugf(log)
	}
}

func Warn(args ...interface{}) {
	WarnA(args...)
}

func WarnA(args ...interface{}) {
	if atom.Enabled(zap.WarnLevel) {
		log := Caller(zap.WarnLevel) + fmt.Sprintln(args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		metrics.GaugeInc(metrics.WarnCount)
		if IsArgsTooLang(zap.WarnLevel, log) {
			return
		}
		sugar.Warn(log)
	}
}

func Warnf(template string, args ...interface{}) {
	WarnAf(template, args...)
}

func WarnAf(template string, args ...interface{}) {
	if atom.Enabled(zap.WarnLevel) {
		log := Caller(zap.WarnLevel) + fmt.Sprintf(template, args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		metrics.GaugeInc(metrics.WarnCount)
		if IsArgsTooLang(zap.WarnLevel, log) {
			return
		}
		sugar.Warnf(log)
	}
}

func WarnDelayf(delay int64, template string, args ...interface{}) {
	WarnDelayAf(delay, template, args...)
}

func WarnDelayAf(delay int64, template string, args ...interface{}) {
	if atom.Enabled(zap.WarnLevel) {
		if delay < baseconf.GetBaseConf().DelayLogLimit {
			// 未达到阈值
			return
		}
		if args == nil {
			args = make([]interface{}, 0)
		}
		template += " cost-time:%dms"
		args = append(args, delay)
		log := Caller(zap.WarnLevel) + fmt.Sprintf(template, args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().IsDebug && len(baseconf.GetBaseConf().FeishuRobot) > 0 {
			// TODO 临时关闭
			// PushLog2Chat(baseconf.GetBaseConf().FeishuRobot, "DELAY", log)
		}
		metrics.GaugeInc(metrics.WarnCount)
		if IsArgsTooLang(zap.WarnLevel, log) {
			return
		}
		sugar.Warnf(log)
	}
}

func Error(args ...interface{}) {
	ErrorA(args...)
}

func ErrorA(args ...interface{}) {
	if atom.Enabled(zap.ErrorLevel) {
		log := Caller(zap.ErrorLevel) + fmt.Sprintln(args...)
		callStack := "===>>>CallStack\n" + string(debug.Stack())
		logStack := log + "\n" + callStack
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(logStack)
		}
		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().IsDebug && len(baseconf.GetBaseConf().FeishuRobot) > 0 {
			PushLog2Chat(baseconf.GetBaseConf().FeishuRobot, "ERROR", logStack)
		}
		metrics.GaugeInc(metrics.ErrorCount)
		if IsArgsTooLang(zap.ErrorLevel, log) {
			return
		}
		sugar.Error(logStack)
	}
}

func Errorf(template string, args ...interface{}) {
	ErrorAf(template, args...)
}

func ErrorAf(template string, args ...interface{}) {
	if atom.Enabled(zap.ErrorLevel) {
		log := Caller(zap.ErrorLevel) + fmt.Sprintf(template, args...)
		callStack := "===>>>CallStack\n" + string(debug.Stack())
		logStack := log + "\n" + callStack
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(logStack)
		}
		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().IsDebug && len(baseconf.GetBaseConf().FeishuRobot) > 0 {
			PushLog2Chat(baseconf.GetBaseConf().FeishuRobot, "ERROR", logStack)
		}
		metrics.GaugeInc(metrics.ErrorCount)
		if IsArgsTooLang(zap.ErrorLevel, log) {
			return
		}
		sugar.Error(logStack)
	}
}

func Trace(args ...interface{}) {
	ErrorA(args...)
}

func Tracef(template string, args ...interface{}) {
	ErrorAf(template, args...)
}

func Fatal(args ...interface{}) {
	FatalA(args...)
}

func FatalA(args ...interface{}) {
	if atom.Enabled(zap.FatalLevel) {
		log := Caller(zap.FatalLevel) + fmt.Sprintln(args...)
		callStack := "===>>>CallStack\n" + string(debug.Stack())
		logStack := log + "\n" + callStack
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(logStack)
		}
		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().IsDebug && len(baseconf.GetBaseConf().FeishuRobot) > 0 {
			PushLog2Chat(baseconf.GetBaseConf().FeishuRobot, "FATAL", logStack)
		}
		metrics.GaugeInc(metrics.FatalCount)
		if IsArgsTooLang(zap.FatalLevel, log) {
			return
		}
		sugar.Fatal(log)
	}
}

func Fatalf(template string, args ...interface{}) {
	FatalAf(template, args...)
}

func FatalAf(template string, args ...interface{}) {
	if atom.Enabled(zap.FatalLevel) {
		log := Caller(zap.FatalLevel) + fmt.Sprintf(template, args...)
		callStack := "===>>>CallStack\n" + string(debug.Stack())
		logStack := log + "\n" + callStack
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(logStack)
		}
		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().IsDebug && len(baseconf.GetBaseConf().FeishuRobot) > 0 {
			PushLog2Chat(baseconf.GetBaseConf().FeishuRobot, "FATAL", logStack)
		}
		metrics.GaugeInc(metrics.FatalCount)
		if IsArgsTooLang(zap.FatalLevel, log) {
			return
		}
		sugar.Fatalf(log)
	}
}

func Info(args ...interface{}) {
	InfoA(args...)
}

func InfoA(args ...interface{}) {
	if atom.Enabled(zap.InfoLevel) {
		log := Caller(zap.InfoLevel) + fmt.Sprintln(args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		if IsArgsTooLang(zap.InfoLevel, log) {
			return
		}
		sugar.Info(log)
	}
}

func Infof(template string, args ...interface{}) {
	InfoAf(template, args...)
}

func InfoAf(template string, args ...interface{}) {
	if atom.Enabled(zap.InfoLevel) {
		log := Caller(zap.InfoLevel) + fmt.Sprintf(template, args...)
		if baseconf.GetBaseConf() != nil && !global.IsCloud {
			SaveToRedis(log)
		}
		if IsArgsTooLang(zap.InfoLevel, log) {
			return
		}
		sugar.Infof(log)

	}
}

func PushLog2Chat(url, title, text string) {
	var result interface{}
	text = fmt.Sprintf("server: [%s]\ngate: [%s]\nlog: %s\n", global.HostName, global.GateWay, text)
	msg := &FeishuMsg{MsgType: "post", Content: FeishuContent{Post: FeishuZh_cn{Zh_cn: FeishuTitle{Title: title, Content: [][]FeishuTitleContent{[]FeishuTitleContent{FeishuTitleContent{Tag: "text", Text: text}}}}}}}
	err := http.Post(url, msg, &result, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func Flush() {
	sugar.Sync()
}
