package dlog

import (
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/global"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	LogChanSizeDefault = 1024
)

type Dlog struct {
	Id      uint64
	AppId   string
	Type    int32
	Content string
}

type DlogThread struct {
	log   *zap.SugaredLogger
	logCh chan *Dlog
}

type Dlogger struct {
	logger    *zap.SugaredLogger
	cutType   int32
	logThread map[int]*DlogThread
	threads   int
}

type Conf struct {
	DataLogPath   string `json:"datalogpath"`   // 统计日志输出路径
	DataLogThread int    `json:"datalogthread"` // 运营日志写线程数
}

var dlogger *Dlogger

func GetDlogger() *Dlogger {
	if dlogger == nil {
		dlogger = &Dlogger{cutType: baseconf.GetBaseConf().LogCutType}
	}
	return dlogger
}

func Init(logPath, fileName string, threads int) error {
	return GetDlogger().Init(logPath, fileName, threads)
}

func Flush() {
	GetDlogger().Flush()
}

func Write(content string) {
	GetDlogger().Write(content)
}

func (p *Dlogger) Init(logPath, fileName string, threads int) error {

	if logPath != "" {

		if utils.PathExists(logPath) == false {
			os.MkdirAll(logPath, 0755)
		}

		if fileName == "" {
			fileName = utils.GetProcName()
		}
		if threads < 0 {
			threads = 0
		}
		if global.IsCloud {
			fileName = "dlog-" + fileName + "-" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + "-" + strconv.Itoa(os.Getpid())
		} else {
			fileName = "dlog-" + fileName
		}
		atomLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
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

		var syncWriter zapcore.WriteSyncer

		switch p.cutType {
		case logger.LOG_SIZE_CUT_TYPE:
			syncWriter = zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(logPath, fileName+".log"),
				MaxAge:     baseconf.GetBaseConf().LogMaxAges,
				MaxBackups: baseconf.GetBaseConf().LogMaxBackups,
				MaxSize:    baseconf.GetBaseConf().LogMaxSize,
				Compress:   baseconf.GetBaseConf().LogCompress,
			})
		case logger.LOG_TIME_CUT_TYPE:
			logWriter, err := rotatelogs.New(
				filepath.Join(logPath, fileName+".log"),
				rotatelogs.WithMaxAge(time.Duration(baseconf.GetBaseConf().LogMaxAges)*time.Hour*24),
				rotatelogs.WithRotationTime(time.Minute*time.Duration(baseconf.GetBaseConf().LogRotationTime)),
				rotatelogs.WithRotationSize(int64(baseconf.GetBaseConf().LogMaxSize*1024*1024)),
			)
			if err != nil {
				return err
			}
			syncWriter = zapcore.AddSync(logWriter)
		}
		core := zapcore.NewCore(encoder, syncWriter, atomLevel)
		p.logger = zap.New(core, zap.AddCaller()).Sugar()
	} else {
		return errors.New(fmt.Sprint("dlog init failed,log file is nil!", logPath, fileName, threads))
	}

	p.threads = threads
	p.logThread = make(map[int]*DlogThread)
	for i := 0; i < p.threads; i++ {
		p.logThread[i] = &DlogThread{}
		p.logThread[i].Init(p.logger)
	}
	logger.Info("dlog init succeed,", logPath, fileName, p.threads)
	return nil
}

func (p *Dlogger) Write(content string) {
	if p.threads > 0 {
		p.logThread[int(time.Now().Unix()%int64(p.threads))].Log(&Dlog{Content: content})
	} else {
		//p.logger.Info(appId, zap.Uint64("logId", logId), zap.String("appId", appId), zap.String("token", ""), zap.Int32("type", ntype), zap.String("content", content))
		p.logger.Info(content)
	}
}

func (p *Dlogger) Flush() {
	if p.logger == nil {
		return
	}
	err := p.logger.Sync()
	if err != nil {
		logger.Errorf("error Flush:%v", err.Error())
	}
}

func (p *DlogThread) Init(log *zap.SugaredLogger) {
	p.log = log
	p.logCh = make(chan *Dlog, LogChanSizeDefault)
	go p.runThreadAsync()
}

func (p *DlogThread) Log(log *Dlog) {
	p.logCh <- log
}

func (p *DlogThread) runThreadAsync() {
	for log := range p.logCh {
		//p.log.Info(log.AppId, zap.Uint64("logId", log.Id), zap.String("appId", log.AppId), zap.String("token", ""), zap.Int32("type", log.Type), zap.String("content", log.Content))
		p.log.Info(log.Content)
	}
}
