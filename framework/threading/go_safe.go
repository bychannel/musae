package threading

import (
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
)

// GoSafeWithParam runs the given fn using another goroutine with param, recovers if fn panics.
func GoSafeWithParam(fn func(param interface{}), param interface{}) {
	go RunSafeWithParam(fn, param)
}

// RunSafeWithParam runs the given fn with param, recovers if fn panics.
func RunSafeWithParam(fn func(param interface{}), param interface{}) {
	defer func() {
		if err := recover(); err != any(nil) {
			logger.Error("RunSafe recover, err: ", err)
		}
	}()

	fn(param)
}

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer func() {
		if err := recover(); err != any(nil) {
			logger.Error("RunSafe recover, err: ", err)
		}
	}()

	fn()
}
