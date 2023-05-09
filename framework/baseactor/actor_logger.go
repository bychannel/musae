package baseactor

import (
	"fmt"
	"github.com/dapr/go-sdk/actor"

	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
)

type ActorLogger struct {
	a actor.Server
	h string
}

func (b *ActorLogger) SetActor(actor actor.Server) {
	b.a = actor
}

func (b *ActorLogger) prefix() string {
	return b.a.ID() + " " + b.h
}

func (b *ActorLogger) Debug(args ...interface{}) {
	logger.DebugA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Debugf(template string, args ...interface{}) {
	logger.DebugAf(b.prefix()+" "+template, args...)
}

func (b *ActorLogger) Warn(args ...interface{}) {
	logger.WarnA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Warnf(template string, args ...interface{}) {
	logger.WarnAf(b.prefix()+" "+template, args...)
}

func (b *ActorLogger) WarnDelayf(delay int64, template string, args ...interface{}) {
	logger.WarnDelayAf(delay, b.prefix()+" "+template, args...)
}

func (b *ActorLogger) Error(args ...interface{}) {
	logger.ErrorA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Errorf(template string, args ...interface{}) {
	logger.ErrorAf(b.prefix()+" "+template, args...)
}

func (b *ActorLogger) Trace(args ...interface{}) {
	logger.ErrorA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Tracef(template string, args ...interface{}) {
	logger.ErrorAf(b.prefix()+" "+template, args...)
}

func (b *ActorLogger) Fatal(args ...interface{}) {
	logger.FatalA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Fatalf(template string, args ...interface{}) {
	logger.FatalAf(b.prefix()+" "+template, args...)
}

func (b *ActorLogger) Info(args ...interface{}) {
	logger.FatalA(b.prefix() + " " + fmt.Sprintln(args...))
}

func (b *ActorLogger) Infof(template string, args ...interface{}) {
	logger.InfoAf(b.prefix()+" "+template, args...)
}
