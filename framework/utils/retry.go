package utils

import (
	"github.com/pkg/errors"
	"time"
)

type retryFunc1[T any] func() (T, error)

// RetryDoSync 同步重试执行
//@param retryCount 重试次数
//@param intervalMs 重试间隔时间(毫秒)
//@param fn 重试执行的函数
func RetryDoSync[T any](retryCount int, fn retryFunc1[T]) (T, error) {
	return RetryDoSyncInterval(retryCount, 100, fn)
}

// RetryDoSyncInterval 同步重试执行
//@param retryCount 重试次数
//@param intervalMs 重试间隔时间(毫秒)
//@param fn 重试执行的函数
func RetryDoSyncInterval[T any](retryCount int, intervalMs int, fn retryFunc1[T]) (T, error) {
	var nilT T

	if retryCount <= 0 {
		return nilT, errors.New("retry count must be positive")
	}

	var res T
	var err error
	for i := 0; i < retryCount; i++ {
		res, err = fn()
		if err == nil {
			return res, nil
		}

		// 间隔时间
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
	return nilT, errors.New("retry got nil return value")
}
