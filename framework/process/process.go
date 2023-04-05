package process

import (
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"os"
	"os/signal"
	"syscall"
)

type Process struct {
	srv base.IServer
}

func NewProcess(srv base.IServer) base.IProcess {
	process := &Process{srv: srv}
	return process
}

func Start(srv base.IServer, opts ...base.FProcessOption) {
	proc := NewProcess(srv)
	proc.Start(opts...)
}

func (p *Process) Start(opts ...base.FProcessOption) {
	defer func() {
		if err := recover(); err != any(nil) {
			logger.Fatal("[process] Start failed, err: ", err)
		}
	}()

	p.srv.SetState(base.PState_Starting)
	p.initSignal()
	if err := p.srv.Init(); err != nil {
		logger.Fatal("[server] init cfg failed, err:", err)
	}

	for _, o := range opts {
		if err := o(p.srv); err != nil {
			logger.Fatal("[server] init opts failed, err:", err)
		}
	}

	if err := p.srv.Start(); err != nil {
		logger.Fatal("[server] Start failed, err:", err)
	}

	p.Main()
}

func (p *Process) initSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGHUP,  //终端结束
		syscall.SIGINT,  //用户发送, 字符(Ctrl+C)触发
		syscall.SIGQUIT, //用户发送, QUIT字符(Ctrl+/)触发
		syscall.SIGILL,  //非法指令(程序错误、试图执行数据段、栈溢出等)
		syscall.SIGTRAP, //断点调试错误
		syscall.SIGABRT, //调用abort函数触发
		syscall.SIGBUS,  //硬件错误
		syscall.SIGFPE,  //浮点错误
		//syscall.SIGKILL,   //无条件结束程序(不能被捕获、阻塞或忽略) kill -9
		syscall.SIGSEGV,     //段错误
		syscall.SIGPIPE,     //消息管道损坏
		syscall.SIGALRM,     //时钟定时信号
		syscall.SIGTERM,     //结束程序(可以被捕获、阻塞或忽略)
		syscall.Signal(0xa), // syscall.SIGUSR1
		syscall.Signal(0xc), // syscall.SIGUSR2
	)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[process] signal catch error,err:", err)
			}
		}()
		for sig := range ch {
			switch sig {
			case syscall.Signal(0xa):
				logger.Info("[process] signal SIGUSR1,reload config")
				p.Reload()
			case syscall.Signal(0xc):
				logger.Info("[process] signal SIGUSR2", sig)
			case syscall.SIGINT:
				logger.Infof("[process] signal :%v", sig)
				p.Exit()
			case syscall.SIGQUIT:
				logger.Infof("[process] signal :%v", sig)
				p.Exit()
			case syscall.SIGTERM:
				logger.Infof("[process] signal close,%v", sig)
				p.Exit()
			case syscall.SIGHUP:
				logger.Infof("[process] signal:%v", sig)
			default:
				logger.Warnf("[process] signal:%v", sig)
			}
		}
	}()
}

func (p *Process) Reload() {
	defer func() {
		if err := recover(); err != any(nil) {
			logger.Fatal("[process] reload failed, err: ", err)
		}
	}()
	logger.Info("[process] reload")
	p.srv.SetState(base.PState_Loading)
	p.srv.Reload()
	p.srv.SetState(base.PState_Running)
}

func (p *Process) Exit() {
	logger.Infof("[process] graceful exit")
	p.srv.SetState(base.PState_Exiting)
	p.srv.Exit()
	p.srv.GracefulStop()
}

func (p *Process) Status() base.PState {
	logger.Info("[process] state:", p.srv.State())
	return p.srv.State()
}

func (p *Process) Main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("[process] main error:", err)
		}
	}()
	p.srv.SetState(base.PState_Running)
	logger.Infof("[process] run success, %s", p.srv.GetAppID())
	p.srv.Main()
}
