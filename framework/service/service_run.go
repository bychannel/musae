package service

import (
	"errors"
	dapr "github.com/dapr/go-sdk/client"
	"gitlab.musadisca-games.com/wangxw/musae/framework/base"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"net/http"
	"time"
)

func (s *Service) Run() error {

	// start dapr server
	go func() {
		defer func() {
			if err := recover(); err != any(nil) {
				logger.Fatal("[service] server run recover, err: ", err)
			}
		}()

		logger.Info("[service] init", s.InAddr)
		if err := s.svc.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("[service] dapr server start exit, err: ", err)
			}
		}
		logger.Info("[service] exit", s.InAddr)

		//退出进程
		s.ExitCh <- struct{}{}
	}()

	var err error
	times := 0
	time.Sleep(2 * time.Second)
	s.Daprc, err = dapr.NewClientWithPort(s.GRPCPort)
	for err != nil || s.Daprc == nil {
		if times++; times >= 3 {
			logger.Fatalf("[service] NewClient error: %v", err)
		}
		time.Sleep(1 * time.Second)
		s.Daprc, err = dapr.NewClientWithPort(s.GRPCPort)
		logger.Warnf("[service] NewClient error: %v", err)
	}
	logger.Info("[service] init dapr client success")

	// start tcp server
	if s.OutAddr != "" {
		go func() {
			defer func() {
				if err := recover(); err != any(nil) {
					logger.Fatal("tcp server run recover, err: ", err)
				}
			}()
			logger.Info("[service] init", s.OutAddr)
			if err := s.tcp.ListenAndServe("tcp", s.OutAddr); err != nil {
				logger.Fatal("tcp service exit, err: ", err)
			}
			logger.Info("[service] exit", s.OutAddr)
		}()
	}

	if s.WebAddr != "" {
		go func() {
			defer func() {
				if err := recover(); err != any(nil) {
					logger.Fatal("tcp server run recover, err: ", err)
				}
			}()
			logger.Info("[service] init", s.WebAddr)
			if err := s.http.Start(); err != nil {
				logger.Fatal("web server exit, err: %v", err)
			}
			logger.Info("[service] exit", s.WebAddr)
		}()
	}

	return nil
}

func (s *Service) Status() base.PState {
	return s.state
}
