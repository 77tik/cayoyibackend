package util

import (
	"my_backend/weedfilesys/glog"
	"strings"
	"time"
)

var RetryWaitTime = 6 * time.Second

// 重试任务
func RetryUntil(name string, job func() error, onErrFn func(err error) (shouldContinue bool)) {
	waitTime := time.Second
	for {
		err := job()
		if err == nil {
			waitTime = time.Second
			break
		}

		if onErrFn(err) {
			if strings.Contains(err.Error(), "transport") {
				glog.V(0).Infof("retry %s: err: %v", name, err)
			}
			time.Sleep(waitTime)
			if waitTime < RetryWaitTime {
				waitTime += waitTime / 2
			}
			continue
		} else {
			break
		}
	}
}
